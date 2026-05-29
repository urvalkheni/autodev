package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newScanCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a repository for languages, frameworks, and dependencies",
		Long:  `Scan the specified directory (or current directory) to detect all technologies, frameworks, package managers, databases, and infrastructure requirements.`,
		Example: `  autodev scan
  autodev scan ./my-project
  autodev scan --json ./my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			if path == "" {
				path = "."
			}

			return runScan(path, jsonOut)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to scan")
	return cmd
}

func runScan(path string, asJSON bool) error {
	start := time.Now()

	if !asJSON {
		printScanHeader(path)
	}

	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	elapsed := time.Since(start)

	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	printScanResult(result, elapsed)
	return nil
}

func printScanHeader(path string) {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD700")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFD700")).
		Padding(0, 2)

	fmt.Println(headerStyle.Render(fmt.Sprintf("🔍 Scanning: %s", path)))
	fmt.Println()
}

func printScanResult(result *scanner.ScanResult, elapsed time.Duration) {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).PaddingLeft(2)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#FFD700")).
		Foreground(lipgloss.Color("#000000")).
		Padding(0, 1).
		Bold(true)

	if len(result.Languages) > 0 {
		fmt.Println(titleStyle.Render("◆ Languages"))
		for _, lang := range result.Languages {
			fmt.Println(itemStyle.Render("  • " + lang))
		}
		fmt.Println()
	}

	if len(result.Frameworks) > 0 {
		fmt.Println(titleStyle.Render("◆ Frameworks"))
		for _, fw := range result.Frameworks {
			fmt.Println(itemStyle.Render("  • " + fw))
		}
		fmt.Println()
	}

	if len(result.PackageManagers) > 0 {
		fmt.Println(titleStyle.Render("◆ Package Managers"))
		for _, pm := range result.PackageManagers {
			fmt.Println(itemStyle.Render("  • " + pm))
		}
		fmt.Println()
	}

	if len(result.Databases) > 0 {
		fmt.Println(titleStyle.Render("◆ Databases & Services"))
		for _, db := range result.Databases {
			fmt.Println(itemStyle.Render("  • " + db))
		}
		fmt.Println()
	}

	if len(result.Infra) > 0 {
		fmt.Println(titleStyle.Render("◆ Infrastructure"))
		for _, inf := range result.Infra {
			fmt.Println(itemStyle.Render("  • " + inf))
		}
		fmt.Println()
	}

	if result.HasDocker {
		fmt.Println(badgeStyle.Render("🐳 Docker") + "  " + badgeStyle.Render("Container-ready"))
	}
	if result.HasK8s {
		fmt.Println(badgeStyle.Render("☸  Kubernetes") + "  Found k8s manifests")
	}

	if len(result.RecommendedSetup) > 0 {
		fmt.Println()
		fmt.Println(titleStyle.Render("◆ Setup Plan"))
		for i, step := range result.RecommendedSetup {
			fmt.Println(itemStyle.Render(fmt.Sprintf("  %d. %s", i+1, step)))
		}
	}

	fmt.Println()
	fmt.Println(dimStyle.Render(fmt.Sprintf("  Scanned in %s  ·  %d technologies detected", elapsed.Round(time.Millisecond), len(result.Technologies))))
	fmt.Println()
	fmt.Println(dimStyle.Render("  Run 'autodev setup' to install all missing tools."))
	fmt.Println()
}
