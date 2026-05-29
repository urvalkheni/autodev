package cmd

import (
	"fmt"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/autodev-sh/autodev/skills"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newSkillsCmd() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:   "skills [path]",
		Short: "Generate a personalized learning roadmap based on detected technologies",
		Long: `Scan the repository, detect all technologies, and generate a personalized
learning roadmap showing current skills, next steps, and long-term goals.
Powered by skills.sh.`,
		Example: `  autodev skills
  autodev skills ./my-project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			return runSkills(path)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to scan")
	return cmd
}

func runSkills(path string) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(titleStyle.Render("  AutoDev Skills — Learning Roadmap"))
	fmt.Println(dimStyle.Render("  Powered by skills.sh"))
	fmt.Println()

	// Scan to detect technologies
	s := scanner.New(path)
	result, err := s.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Combine all detected tech names
	var detected []string
	detected = append(detected, result.Languages...)
	detected = append(detected, result.Frameworks...)
	detected = append(detected, result.PackageManagers...)
	detected = append(detected, result.Databases...)

	if len(detected) == 0 {
		fmt.Println(dimStyle.Render("  No technologies detected. Try running in a project directory."))
		return nil
	}

	// Generate roadmap
	gen := skills.New()
	roadmap := gen.Generate(detected)
	roadmap.Print()

	fmt.Println(dimStyle.Render("  Visit https://skills.sh for interactive learning paths and courses."))
	fmt.Println()

	return nil
}
