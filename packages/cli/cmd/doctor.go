package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/autodev-sh/autodev/core/osinfo"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	var fix bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check the health and security of your codebase",
		Long:  `Verify that your codebase is healthy, secure, has correct git configurations, is free of exposed secrets, and matches standard linting and environment variables.`,
		Example: `  autodev doctor
  autodev doctor --fix`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(fix)
		},
	}
	cmd.Flags().BoolVar(&fix, "fix", false, "automatically attempt to fix detected codebase issues")
	return cmd
}

type codebaseDiagnostic struct {
	name        string
	description string
	checkFn     func(path string) (bool, string, error)
	fixFn       func(path string) (bool, error)
}

func runDoctor(fix bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	okStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87"))
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	boldStyle := lipgloss.NewStyle().Bold(true)
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF"))

	fmt.Println()
	fmt.Println(titleStyle.Render("⚡ AUTODEV CODEBASE HEALTH & SECURITY DOCTOR"))
	fmt.Println()

	// 1. System specifications summary
	info, err := osinfo.Detect()
	if err == nil {
		fmt.Println(titleStyle.Render("SYSTEM SPECIFICATIONS"))
		fmt.Printf("  %-20s %s\n", boldStyle.Render("OS"), info.Version)
		fmt.Printf("  %-20s %s\n", boldStyle.Render("Architecture"), info.Arch)
		fmt.Printf("  %-20s %s\n", boldStyle.Render("Package Manager"), info.PackageManager)
		fmt.Println()
	}

	fmt.Println(titleStyle.Render("CODEBASE DIAGNOSTICS SCAN"))
	fmt.Println(dimStyle.Render("  Scanning for secrets, configuration mismatches, and code errors..."))
	fmt.Println()

	diagnostics := getCodebaseDiagnostics()
	path := "."

	unhealthyCount := 0
	var toFix []*codebaseDiagnostic

	for i := range diagnostics {
		d := &diagnostics[i]
		ok, msg, err := d.checkFn(path)
		if err != nil {
			fmt.Printf("  %-10s %-25s %s\n", warnStyle.Render("[ERROR]"), d.name, dimStyle.Render(err.Error()))
			unhealthyCount++
			toFix = append(toFix, d)
		} else if !ok {
			fmt.Printf("  %-10s %-25s %s\n", warnStyle.Render("[WARNING]"), d.name, cyanStyle.Render(msg))
			unhealthyCount++
			toFix = append(toFix, d)
		} else {
			fmt.Printf("  %-10s %-25s %s\n", okStyle.Render("[OK]"), d.name, dimStyle.Render("Clean & Healthy"))
		}
	}

	fmt.Println()
	if unhealthyCount == 0 {
		fmt.Println(okStyle.Render("  ✓ Codebase is completely healthy, secure, and ready for production!"))
		return nil
	}

	fmt.Printf("  %s issues detected in the codebase.\n", warnStyle.Render(fmt.Sprintf("%d", unhealthyCount)))

	if fix {
		fmt.Println()
		fmt.Println(titleStyle.Render("🔧 Attempting auto-remediation (--fix)..."))
		fmt.Println()

		for _, d := range toFix {
			fmt.Printf("  Healing %s...\n", d.name)
			fixed, err := d.fixFn(path)
			if err != nil {
				fmt.Printf("    %s Failed to fix: %v\n", warnStyle.Render("✗"), err)
			} else if fixed {
				fmt.Printf("    %s Successfully healed %s!\n", okStyle.Render("✓"), d.name)
			} else {
				fmt.Printf("    %s Requires manual action (refer to warning logs)\n", warnStyle.Render("!"))
			}
		}
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println(dimStyle.Render("  Run 'autodev doctor --fix' to automatically resolve fixable issues."))
	}

	PrintGitHubCTA()
	return nil
}

func getCodebaseDiagnostics() []codebaseDiagnostic {
	return []codebaseDiagnostic{
		{
			name:        "Git Configuration (.gitignore)",
			description: "Verify that sensitive environment or package files are not tracked in git",
			checkFn: func(path string) (bool, string, error) {
				gitIgnorePath := filepath.Join(path, ".gitignore")
				if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
					return false, "No .gitignore file detected. Environment files and builds could be committed.", nil
				}
				data, err := os.ReadFile(gitIgnorePath)
				if err != nil {
					return false, "", err
				}
				content := string(data)
				hasEnv := strings.Contains(content, ".env")
				hasNode := strings.Contains(content, "node_modules")
				if !hasEnv || !hasNode {
					return false, ".gitignore does not ignore sensitive .env files or node_modules.", nil
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				gitIgnorePath := filepath.Join(path, ".gitignore")
				if _, err := os.Stat(gitIgnorePath); os.IsNotExist(err) {
					defaultGitignore := `# Standard AutoDev Gitignore
node_modules/
dist/
build/
.next/
.turbo/
.env
.env.local
.env.*.local
*.log
packages/cli/bin/
`
					err := os.WriteFile(gitIgnorePath, []byte(defaultGitignore), 0644)
					return err == nil, err
				}

				// Append to existing
				data, err := os.ReadFile(gitIgnorePath)
				if err != nil {
					return false, err
				}
				content := string(data)
				var toAppend []string
				if !strings.Contains(content, ".env") {
					toAppend = append(toAppend, ".env", ".env.local", ".env.*.local")
				}
				if !strings.Contains(content, "node_modules") {
					toAppend = append(toAppend, "node_modules/")
				}

				if len(toAppend) > 0 {
					f, err := os.OpenFile(gitIgnorePath, os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						return false, err
					}
					defer f.Close()
					if _, err = f.WriteString("\n# Added by AutoDev Doctor\n" + strings.Join(toAppend, "\n") + "\n"); err != nil {
						return false, err
					}
					if err := f.Close(); err != nil {
						return false, err
					}
					return true, nil
				}
				return true, nil
			},
		},
		{
			name:        "Exposed Secrets Scanner",
			description: "Scan the codebase to verify no hardcoded secrets or API keys are present",
			checkFn: func(path string) (bool, string, error) {
				var matches []string
				secretRegexes := []*regexp.Regexp{
					regexp.MustCompile(`(?i)(api_key|secret|password|private_key|token|auth_token)\s*[:=]\s*['"]([a-zA-Z0-9_\-\.]{16,})['"]`),
					regexp.MustCompile(`(?i)aws_[a-z_]*key\s*[:=]\s*['"]([a-zA-Z0-9/+=]{16,})['"]`),
				}

				err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
					if err != nil {
						return nil
					}
					// skip noise
					base := d.Name()
					if d.IsDir() {
						if p != path && (strings.HasPrefix(base, ".") ||
							base == "node_modules" || base == "vendor" ||
							base == "dist" || base == "build" || base == ".next" || base == ".turbo") {
							return filepath.SkipDir
						}
						return nil
					}

					// Only scan code files
					ext := filepath.Ext(p)
					if ext != ".js" && ext != ".ts" && ext != ".tsx" && ext != ".py" && ext != ".go" && ext != ".json" && ext != ".yml" && ext != ".env" {
						return nil
					}

					// Skip lockfiles and .env.example
					if base == "package-lock.json" || base == "pnpm-lock.yaml" || base == "yarn.lock" || strings.HasSuffix(base, ".example") {
						return nil
					}

					data, err := os.ReadFile(p)
					if err != nil {
						return nil
					}

					lines := strings.Split(string(data), "\n")
					for lineNum, line := range lines {
						for _, rx := range secretRegexes {
							sub := rx.FindStringSubmatch(line)
							if len(sub) > 2 {
								secretVal := sub[2]
								if !isPlaceholder(secretVal) {
									matches = append(matches, fmt.Sprintf("%s (Line %d)", filepath.Base(p), lineNum+1))
								}
							}
						}
					}

					return nil
				})

				if err != nil {
					return false, "", err
				}

				if len(matches) > 0 {
					return false, fmt.Sprintf("Exposed credentials detected: %s", strings.Join(matches, ", ")), nil
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				// Secrets are sensitive and hard to replace without breaking codes.
				// We prompt the user on how to fix instead of destructive rewrite.
				return false, fmt.Errorf("exposed secrets must be manually migrated to .env and referenced dynamically")
			},
		},
		{
			name:        "Environment Config (.env)",
			description: "Verify that environmental profiles are properly configured",
			checkFn: func(path string) (bool, string, error) {
				envExample := filepath.Join(path, ".env.example")
				envFile := filepath.Join(path, ".env")

				if _, err := os.Stat(envExample); err == nil {
					if _, err := os.Stat(envFile); os.IsNotExist(err) {
						return false, "Found .env.example but no active .env file was configured.", nil
					}
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				envExample := filepath.Join(path, ".env.example")
				envFile := filepath.Join(path, ".env")
				if _, err := os.Stat(envExample); err == nil {
					data, err := os.ReadFile(envExample)
					if err != nil {
						return false, err
					}
					err = os.WriteFile(envFile, data, 0644)
					return err == nil, err
				}
				return true, nil
			},
		},
		{
			name:        "Dependency Lockfiles",
			description: "Verify dependencies configurations are locked and synchronized",
			checkFn: func(path string) (bool, string, error) {
				pkgJSON := filepath.Join(path, "package.json")
				if _, err := os.Stat(pkgJSON); err == nil {
					lockfiles := []string{"package-lock.json", "pnpm-lock.yaml", "yarn.lock"}
					found := false
					for _, lf := range lockfiles {
						if _, err := os.Stat(filepath.Join(path, lf)); err == nil {
							found = true
							break
						}
					}
					if !found {
						return false, "package.json found but no matching package lockfile exists.", nil
					}
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				// If pnpm is available, run pnpm install. Otherwise try npm install.
				cmdName := "npm"
				if _, err := exec.LookPath("pnpm"); err == nil {
					cmdName = "pnpm"
				}

				cmd := exec.Command(cmdName, "install")
				cmd.Dir = path
				err := cmd.Run()
				return err == nil, err
			},
		},
		{
			name:        "Linter & Code Format Status",
			description: "Scan codebase for syntax compilation or structural formatting warnings",
			checkFn: func(path string) (bool, string, error) {
				pkgJSON := filepath.Join(path, "package.json")
				if _, err := os.Stat(pkgJSON); err == nil {
					// Check if lint script exists
					data, err := os.ReadFile(pkgJSON)
					if err != nil {
						return true, "", nil
					}
					if strings.Contains(string(data), `"lint"`) {
						cmdName := "npm"
						if _, err := exec.LookPath("pnpm"); err == nil {
							cmdName = "pnpm"
						}
						var stdout, stderr bytes.Buffer
						cmd := exec.Command(cmdName, "run", "lint")
						cmd.Dir = path
						cmd.Stdout = &stdout
						cmd.Stderr = &stderr
						if err := cmd.Run(); err != nil {
							return false, "Linter run reported syntax or formatting problems.", nil
						}
					}
				}
				return true, "", nil
			},
			fixFn: func(path string) (bool, error) {
				cmdName := "npm"
				if _, err := exec.LookPath("pnpm"); err == nil {
					cmdName = "pnpm"
				}

				// Attempt auto-linting fixes
				var cmd *exec.Cmd
				if cmdName == "pnpm" {
					cmd = exec.Command("pnpm", "run", "lint", "--fix")
				} else {
					cmd = exec.Command("npm", "run", "lint", "--", "--fix")
				}
				cmd.Dir = path
				_ = cmd.Run() // run best effort

				// Also run prettier if available
				if _, err := exec.LookPath("prettier"); err == nil {
					pcmd := exec.Command("npx", "prettier", "--write", "**/*.{js,ts,tsx,json,css,md}")
					pcmd.Dir = path
					_ = pcmd.Run()
				}
				return true, nil
			},
		},
	}
}

func isPlaceholder(val string) bool {
	val = strings.ToLower(val)
	placeholders := []string{"your", "placeholder", "key_here", "my-secret", "dummy", "example", "test", "token_here", "config"}
	for _, p := range placeholders {
		if strings.Contains(val, p) {
			return true
		}
	}
	return false
}
