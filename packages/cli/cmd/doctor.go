package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the health of your development environment",
		Long:  `Verify that all common development tools are installed and working correctly. Reports versions, warns about missing tools, and suggests fixes.`,
		Example: `  autodev doctor`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor()
		},
	}
}

type check struct {
	name    string
	cmd     string
	args    []string
	hint    string
}

var checks = []check{
	{name: "Git", cmd: "git", args: []string{"--version"}, hint: "https://git-scm.com/downloads"},
	{name: "Node.js", cmd: "node", args: []string{"--version"}, hint: "autodev install nodejs"},
	{name: "npm", cmd: "npm", args: []string{"--version"}, hint: "Comes with Node.js"},
	{name: "pnpm", cmd: "pnpm", args: []string{"--version"}, hint: "npm install -g pnpm"},
	{name: "yarn", cmd: "yarn", args: []string{"--version"}, hint: "npm install -g yarn"},
	{name: "Bun", cmd: "bun", args: []string{"--version"}, hint: "https://bun.sh"},
	{name: "Go", cmd: "go", args: []string{"version"}, hint: "autodev install go"},
	{name: "Python 3", cmd: "python3", args: []string{"--version"}, hint: "autodev install python"},
	{name: "pip", cmd: "pip3", args: []string{"--version"}, hint: "Comes with Python 3"},
	{name: "Rust", cmd: "rustc", args: []string{"--version"}, hint: "autodev install rust"},
	{name: "Docker", cmd: "docker", args: []string{"--version"}, hint: "autodev install docker"},
	{name: "docker compose", cmd: "docker", args: []string{"compose", "version"}, hint: "Upgrade Docker Desktop"},
	{name: "kubectl", cmd: "kubectl", args: []string{"version", "--client", "--short"}, hint: "autodev install kubectl"},
	{name: "Terraform", cmd: "terraform", args: []string{"version"}, hint: "autodev install terraform"},
	{name: "Flutter", cmd: "flutter", args: []string{"--version"}, hint: "autodev install flutter"},
	{name: "Java", cmd: "java", args: []string{"-version"}, hint: "autodev install java"},
	{name: "PHP", cmd: "php", args: []string{"--version"}, hint: "autodev install php"},
	{name: "Ruby", cmd: "ruby", args: []string{"--version"}, hint: "autodev install ruby"},
}

func runDoctor() error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	boldStyle := lipgloss.NewStyle().Bold(true)

	fmt.Println()
	fmt.Println(titleStyle.Render("AUTODEV DOCTOR - ENVIRONMENT HEALTH CHECK"))
	fmt.Println()

	// System info
	info, err := osinfo.Detect()
	if err == nil {
		fmt.Println(titleStyle.Render("SYSTEM SPECIFICATIONS"))
		fmt.Printf("  %-20s %s\n", boldStyle.Render("OS"), info.Version)
		fmt.Printf("  %-20s %s\n", boldStyle.Render("Architecture"), info.Arch)
		fmt.Printf("  %-20s %d cores\n", boldStyle.Render("CPU"), info.CPUCores)
		fmt.Printf("  %-20s %s\n", boldStyle.Render("RAM"), osinfo.FormatRAM(info.RAMBytes))
		fmt.Printf("  %-20s %s\n", boldStyle.Render("Package Manager"), info.PackageManager)
		fmt.Println()
	}

	// Tool checks
	fmt.Println(titleStyle.Render("MANAGED TOOLS CHECK"))
	installed := 0
	missing := 0

	type checkResult struct {
		version string
		err     error
	}

	results := make([]checkResult, len(checks))
	var wg sync.WaitGroup
	wg.Add(len(checks))

	for i, c := range checks {
		go func(idx int, ch check) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
			defer cancel()

			cmd := exec.CommandContext(ctx, ch.cmd, ch.args...)
			out, err := cmd.CombinedOutput()
			results[idx] = checkResult{
				version: strings.TrimSpace(strings.Split(string(out), "\n")[0]),
				err:     err,
			}
		}(i, c)
	}

	wg.Wait()

	for i, c := range checks {
		res := results[i]
		if res.err != nil || strings.TrimSpace(res.version) == "" {
			fmt.Printf("  %-10s %-20s %s\n",
				warnStyle.Render("[MISSING]"),
				c.name,
				dimStyle.Render(fmt.Sprintf("not found — %s", c.hint)),
			)
			missing++
		} else {
			version := res.version
			// Trim verbose output
			if len(version) > 50 {
				version = version[:50]
			}
			fmt.Printf("  %-10s %-20s %s\n",
				okStyle.Render("[OK]"),
				c.name,
				dimStyle.Render(version),
			)
			installed++
		}
	}

	fmt.Println()
	fmt.Printf("  %s and %s\n",
		okStyle.Render(fmt.Sprintf("%d tools installed", installed)),
		warnStyle.Render(fmt.Sprintf("%d missing", missing)),
	)

	if missing > 0 {
		fmt.Println()
		fmt.Println(dimStyle.Render("  Run 'autodev setup' to install missing tools."))
	}
	fmt.Println()

	return nil
}
