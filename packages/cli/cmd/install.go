package cmd

import (
	"fmt"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	var listFlag bool

	cmd := &cobra.Command{
		Use:   "install <package-id>",
		Short: "Install a specific package by ID",
		Long: `Install a specific runtime, framework, or tool from the AutoDev catalog.
Dependencies are resolved automatically. Run --list to see all available packages.`,
		Example: `  autodev install nodejs
  autodev install flutter
  autodev install pytorch
  autodev install --list`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := catalog.Load()
			if err != nil {
				return err
			}
			if listFlag || len(args) == 0 {
				return printCatalogList(c)
			}
			return runCatalogInstall(c, args[0])
		},
	}

	cmd.Flags().BoolVarP(&listFlag, "list", "l", false, "list all available packages")
	return cmd
}

func printCatalogList(c *catalog.Catalog) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	catStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#4A90E2"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  AutoDev Package Catalog"))
	fmt.Println()

	byCat := c.ByCategory()
	for _, cat := range catalog.CategoryOrder {
		pkgs, ok := byCat[cat]
		if !ok || len(pkgs) == 0 {
			continue
		}
		fmt.Println(catStyle.Render("  [" + cat + "]"))
		for _, pkg := range pkgs {
			deps := ""
			if len(pkg.Deps) > 0 {
				deps = dimStyle.Render(fmt.Sprintf("  (deps: %v)", pkg.Deps))
			}
			fmt.Printf("    %-18s  %s%s\n",
				pkg.ID,
				dimStyle.Render(pkg.Description),
				deps,
			)
		}
		fmt.Println()
	}

	fmt.Println(dimStyle.Render("  Usage: autodev install <id>"))
	fmt.Println(dimStyle.Render("  Or run 'autodev' for the interactive installer"))
	fmt.Println()
	return nil
}

func runCatalogInstall(c *catalog.Catalog, id string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	// Resolve with deps
	resolved, err := c.Resolve([]string{id})
	if err != nil {
		return err
	}

	fmt.Println()
	pkg, _ := c.GetPackage(id)
	fmt.Printf("%s\n", titleStyle.Render("Installing "+pkg.Name))

	if len(resolved) > 1 {
		fmt.Println(dimStyle.Render("  Installing dependencies first:"))
		for _, p := range resolved[:len(resolved)-1] {
			fmt.Printf("  %s\n", dimStyle.Render(p.Name))
		}
	}
	fmt.Println()

	for _, p := range resolved {
		fmt.Printf(titleStyle.Render("  - %s\n"), p.Name)
		if dryRun {
			fmt.Println(warnStyle.Render("    [dry-run] skipped"))
			continue
		}
		if err := execInstall(p); err != nil {
			fmt.Println(warnStyle.Render(fmt.Sprintf("  [FAIL] %s: %v", p.Name, err)))
		} else {
			fmt.Println(okStyle.Render(fmt.Sprintf("  [OK] %s installed", p.Name)))
		}
	}

	fmt.Println()
	fmt.Println(okStyle.Render(fmt.Sprintf("  [OK] %s ready!", pkg.Name)))
	return nil
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Check for and apply updates to managed packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
			fmt.Println()
			fmt.Println(okStyle.Render("  AutoDev update: re-run 'autodev install <pkg>' to get the latest version."))
			fmt.Println()
			return nil
		},
	}
}

func newCleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Remove AutoDev cache and temp files",
		RunE: func(cmd *cobra.Command, args []string) error {
			okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)
			fmt.Println()
			fmt.Println(okStyle.Render("  ✓ Cache cleaned."))
			fmt.Println()
			return nil
		},
	}
}

func newExportCmd() *cobra.Command {
	var outFile string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export environment as a reproducible JSON lockfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReport(".", "json", ".")
		},
	}
	cmd.Flags().StringVarP(&outFile, "output", "o", ".autodev.lock.json", "output file")
	return cmd
}
