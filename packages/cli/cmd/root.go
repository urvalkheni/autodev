// Package cmd defines all AutoDev CLI commands using Cobra.
package cmd

import (
	"fmt"
	"os"

	"github.com/autodev-sh/autodev/catalog"
	"github.com/autodev-sh/autodev/cli/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	noColor bool
	dryRun  bool
	jsonOut bool
)

// rootCmd — running `autodev` with no args opens the interactive TUI.
var rootCmd = &cobra.Command{
	Use:     "autodev",
	Short:   "Clone. Scan. Install. Build. — The App Store for Developers",
	Version: "0.1.0",
	Long: `
  █████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ███████╗██╗   ██╗
  ██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝██║   ██║
  ███████║██║   ██║   ██║   ██║   ██║██║  ██║█████╗  ██║   ██║
  ██╔══██║██║   ██║   ██║   ██║   ██║██║  ██║██╔══╝  ╚██╗ ██╔╝
  ██║  ██║╚██████╔╝   ██║   ╚██████╔╝██████╔╝███████╗ ╚████╔╝ 
  ╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚══════╝ ╚═════╝ ╚══════╝  ╚═══╝ 

  The App Store for Developers.
  Run with no arguments to open the interactive installer.`,
	// When called with no subcommand → open the TUI
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := catalog.Load()
		if err != nil {
			return fmt.Errorf("failed to load catalog: %w", err)
		}
		return tui.Run(c)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: .autodev.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "preview actions without executing")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "output results as JSON")

	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("no_color", rootCmd.PersistentFlags().Lookup("no-color"))

	rootCmd.AddCommand(
		newScanCmd(),
		newSetupCmd(),
		newGitHubCmd(),
		newDoctorCmd(),
		newReportCmd(),
		newInstallCmd(),
		newUpdateCmd(),
		newCleanCmd(),
		newSkillsCmd(),
		newExportCmd(),
		newProfileCmd(),
		newUICmd(),
		newCloneCmd(),
	)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".autodev")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(home + "/.config/autodev")
	}
	viper.SetEnvPrefix("AUTODEV")
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()
}
