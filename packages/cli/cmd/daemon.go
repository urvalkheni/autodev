package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/spf13/cobra"
)

func newDaemonCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Start background AI session prompt capture daemon",
		Long: `Runs a background loop to automatically detect active AI CLI sessions (like gemini,
claude, copilot, etc.) and log command line prompts to your .autodevs folder.`,
		Example: `  autodev daemon`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDaemon()
		},
	}
	return cmd
}

func runDaemon() error {
	engine, err := promptcapture.NewEngine("")
	var root string
	if err != nil {
		root, _ = os.Getwd()
	} else {
		root = engine.Root
	}

	fmt.Printf("Starting AutoDevs Daemon in root: %s\n", root)
	fmt.Println("Monitoring active sessions for: gemini, claude, copilot, agy, codex...")

	stopChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		close(stopChan)
	}()

	promptcapture.StartDaemon(root, stopChan)
	fmt.Println("Daemon stopped.")
	return nil
}
