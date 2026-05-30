package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/autodev-sh/autodev/installer"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	var fix bool
	cmd := &cobra.Command{
		Use:     "doctor",
		Short:   "Check the health of your development environment",
		Long:    `Verify that all common development tools are installed and working correctly. Reports versions, warns about missing tools, and suggests fixes.`,
		Example: `  autodev doctor
  autodev doctor --fix`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(fix)
		},
	}
	cmd.Flags().BoolVar(&fix, "fix", false, "automatically attempt to install missing tools")
	return cmd
}

type check struct {
	name string
	cmd  string
	args []string
	hint string
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

func runDoctor(fix bool) error {
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

	var mu sync.Mutex
	results := make(map[int]checkResult)
	var wg sync.WaitGroup
	wg.Add(len(checks))

	for i, c := range checks {
		go func(idx int, ch check) {
			defer wg.Done()

			// 1. Quick check if the binary exists in PATH
			_, err := exec.LookPath(ch.cmd)
			if err != nil {
				mu.Lock()
				results[idx] = checkResult{version: "", err: err}
				mu.Unlock()
				return
			}

			// 2. If it exists, run version check command with a timeout
			ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
			defer cancel()

			cmd := exec.CommandContext(ctx, ch.cmd, ch.args...)
			out, err := cmd.CombinedOutput()

			mu.Lock()
			results[idx] = checkResult{
				version: strings.TrimSpace(strings.Split(string(out), "\n")[0]),
				err:     err,
			}
			mu.Unlock()
		}(i, c)
	}

	// Channel to signal completion
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(800 * time.Millisecond):
		// Timeout to keep doctor command extremely fast, cutting off hung check processes
	}

	checkToRuntime := map[string]string{
		"Node.js":    "nodejs",
		"pnpm":       "pnpm",
		"Bun":        "bun",
		"Go":         "go",
		"Python 3":   "python",
		"Rust":       "rust",
		"Docker":     "docker",
		"kubectl":    "kubectl",
		"Terraform":  "terraform",
		"Flutter":    "flutter",
		"Java":       "java",
		"PHP":        "php",
		"Ruby":       "ruby",
	}

	var missingRuntimes []string

	for i, c := range checks {
		mu.Lock()
		res, exists := results[i]
		mu.Unlock()

		if !exists || res.err != nil || strings.TrimSpace(res.version) == "" {
			fmt.Printf("  %-10s %-20s %s\n",
				warnStyle.Render("[MISSING]"),
				c.name,
				dimStyle.Render(fmt.Sprintf("not found — %s", c.hint)),
			)
			missing++
			if rt, ok := checkToRuntime[c.name]; ok {
				missingRuntimes = append(missingRuntimes, rt)
			}
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
		if fix && len(missingRuntimes) > 0 {
			fmt.Println()
			fmt.Println(titleStyle.Render("🔧 Attempting auto-remediation (--fix)..."))
			inst := installer.New(false)
			for _, rtName := range missingRuntimes {
				rt, _ := installer.GetRuntime(rtName)
				fmt.Printf("\n  Installing %s...\n", rt.Name)
				if err := inst.Install(rtName); err != nil {
					fmt.Printf("  %s Failed to install %s: %v\n", warnStyle.Render("✗"), rt.Name, err)
				} else {
					fmt.Printf("  %s %s installed successfully\n", okStyle.Render("✓"), rt.Name)
				}
			}
		} else {
			fmt.Println()
			fmt.Println(dimStyle.Render("  Run 'autodev doctor --fix' to automatically install missing tools."))
		}
	}
	fmt.Println()

	PrintGitHubCTA()
	return nil
}
