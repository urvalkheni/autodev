package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/autodev-sh/autodev/core/promptcapture"
	"github.com/spf13/cobra"
)

func newCaptureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "capture [command] [args...]",
		Short: "Capture prompts from a terminal-based AI assistant",
		Long: `Intercepts prompts typed into another terminal-based AI assistant CLI
(like gemini, claude, copilot, etc.) and logs them to your .autodevs folder.`,
		Example: `  autodev capture gemini
  autodev capture claude`,
		DisableFlagParsing: true, // Pass all flags through to the wrapped process
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("please specify the AI CLI command to capture (e.g., gemini, claude)")
			}
			return runCapture(args[0], args[1:])
		},
	}
	return cmd
}

func runCapture(name string, args []string) error {
	subCmd := exec.Command(name, args...)

	// Get stdin pipe to write to
	stdinPipe, err := subCmd.StdinPipe()
	if err != nil {
		return err
	}

	subCmd.Stdout = os.Stdout
	subCmd.Stderr = os.Stderr

	// Find project root for prompt logging
	engine, err := promptcapture.NewEngine("")
	var root string
	if err != nil {
		root, _ = os.Getwd()
	} else {
		root = engine.Root
	}

	// Stdin proxy callback: log prompts when user inputs them
	onLine := func(line string) {
		_ = promptcapture.AppendToPromptsMD(root, line)
		_ = promptcapture.SyncWithDevMentor(root, line)
	}

	// Create stdin proxy wrapping os.Stdin
	stdinProxy := promptcapture.NewStdinProxy(os.Stdin, onLine)

	// Forward from proxy to stdinPipe in a goroutine
	go func() {
		defer stdinPipe.Close()
		_, _ = io.Copy(stdinPipe, stdinProxy)
	}()

	// Start subprocess
	if err := subCmd.Start(); err != nil {
		return fmt.Errorf("failed to start command %s: %w", name, err)
	}

	// Handle terminal interrupts/signals gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range sigChan {
			if subCmd.Process != nil {
				_ = subCmd.Process.Signal(sig)
			}
		}
	}()

	// Wait for process to exit
	return subCmd.Wait()
}
