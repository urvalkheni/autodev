package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/autodev-sh/autodev/installer"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newGitHubCTAForMCP() string {
	return "\n⭐ Love AutoDev? Star the repo: https://github.com/HEETMEHTA18/autodev\n"
}

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Model Context Protocol (MCP) server integration",
		Long:  `Run a native Model Context Protocol (MCP) server over stdin/stdout, allowing AI coding tools like Claude Desktop or Cursor to interface directly with your dev environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPGuide()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start the MCP server over stdin/stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "setup",
		Short: "Automatically configure Claude Desktop to use AutoDev MCP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPSetup()
		},
	})

	return cmd
}

func runMCPGuide() error {
	goldStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	cyanStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00E5FF"))
	greenStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	borderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#00E5FF")).Padding(1, 2)

	fmt.Printf("\n  %s\n", goldStyle.Render("🔌 AUTODEV MODEL CONTEXT PROTOCOL (MCP) INTEGRATION"))
	fmt.Println("  Connect your AI coding tools directly to your local development environment.")
	fmt.Println()

	toolsText := "🛠️  Available AI Tools through AutoDev MCP:\n\n" +
		fmt.Sprintf("• %s: Scan repo structure and technology stack.\n", cyanStyle.Render("autodev_scan")) +
		fmt.Sprintf("• %s: Check system compilers & runtime health.\n", cyanStyle.Render("autodev_doctor")) +
		fmt.Sprintf("• %s: Install missing developer runtimes locally.\n", cyanStyle.Render("autodev_install")) +
		fmt.Sprintf("• %s: Audit package dependencies for security issues.\n", cyanStyle.Render("autodev_audit"))

	fmt.Println(borderStyle.Render(toolsText))
	fmt.Println()

	fmt.Println(goldStyle.Render("  🚀 QUICK CONNECT WITH CLAUDE DESKTOP:"))
	fmt.Println("  You can automatically write configuration for Claude Desktop with one command:")
	fmt.Println("    " + greenStyle.Render("autodev mcp setup"))
	fmt.Println()

	fmt.Println(goldStyle.Render("  📋 MANUAL CONFIGURATION GUIDE:"))
	fmt.Println("  To connect with Cursor, Windsurf, or custom AI clients, add a command-based MCP server:")

	executablePath, err := os.Executable()
	if err != nil || strings.Contains(executablePath, "go-build") || strings.Contains(executablePath, "exe/main") || strings.Contains(executablePath, "/tmp") {
		executablePath = "autodev"
	}

	fmt.Printf("  • Server Command: %s\n", greenStyle.Render(executablePath))
	fmt.Printf("  • Arguments:      %s\n", greenStyle.Render("mcp start"))
	fmt.Println()

	fmt.Println(dimStyle.Render("  Run 'autodev mcp start' to start the server over stdin/stdout manual pipe."))
	fmt.Println()
	return nil
}

func getClaudeDesktopConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var path string
	switch runtime.GOOS {
	case "windows":
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			appdata = filepath.Join(home, "AppData", "Roaming")
		}
		path = filepath.Join(appdata, "Claude", "claude_desktop_config.json")
	case "darwin":
		path = filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	default: // linux
		path = filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}
	return path, nil
}

func runMCPSetup() error {
	configPath, err := getClaudeDesktopConfigPath()
	if err != nil {
		return fmt.Errorf("failed to determine Claude config path: %w", err)
	}

	fmt.Printf("🔍 Locating Claude Desktop configuration...\n")
	fmt.Printf("   Config file path: %s\n", configPath)

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var configData map[string]interface{}
	fileBytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			configData = make(map[string]interface{})
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		if err := json.Unmarshal(fileBytes, &configData); err != nil {
			backupPath := configPath + ".bak"
			_ = os.WriteFile(backupPath, fileBytes, 0644)
			fmt.Printf("⚠️  Existing config was invalid JSON. Backed up to %s and creating new config.\n", backupPath)
			configData = make(map[string]interface{})
		}
	}

	mcpServersRaw, ok := configData["mcpServers"]
	var mcpServers map[string]interface{}
	if !ok {
		mcpServers = make(map[string]interface{})
	} else {
		mcpServers, ok = mcpServersRaw.(map[string]interface{})
		if !ok {
			mcpServers = make(map[string]interface{})
		}
	}

	executablePath, err := os.Executable()
	if err != nil || strings.Contains(executablePath, "go-build") || strings.Contains(executablePath, "exe/main") || strings.Contains(executablePath, "/tmp") {
		executablePath = "autodev"
	}

	mcpServers["autodev"] = map[string]interface{}{
		"command": executablePath,
		"args":    []string{"mcp", "start"},
	}

	configData["mcpServers"] = mcpServers

	newBytes, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize updated config: %w", err)
	}

	if err := os.WriteFile(configPath, newBytes, 0644); err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	// Helper to update any .mcp.json file format
	updateMcpJsonFile := func(filePath string, execPath string) bool {
		if _, err := os.Stat(filePath); err != nil {
			return false
		}
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			return false
		}
		var data map[string]interface{}
		if err := json.Unmarshal(fileBytes, &data); err != nil {
			return false
		}
		mcpServersRaw, ok := data["mcpServers"]
		var mcpServers map[string]interface{}
		if !ok {
			mcpServers = make(map[string]interface{})
		} else {
			mcpServers, ok = mcpServersRaw.(map[string]interface{})
			if !ok {
				mcpServers = make(map[string]interface{})
			}
		}
		mcpServers["autodev"] = map[string]interface{}{
			"type":    "stdio",
			"command": execPath,
			"args":    []string{"mcp", "start"},
		}
		data["mcpServers"] = mcpServers
		newBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return false
		}
		return os.WriteFile(filePath, newBytes, 0644) == nil
	}

	// Update local and global config files
	home, homeErr := os.UserHomeDir()
	configPaths := []string{".mcp.json"}
	if homeErr == nil {
		configPaths = append(configPaths,
			filepath.Join(home, ".mcp.json"),
			filepath.Join(home, ".cursor", "mcp.json"),
		)
	}

	mcpUpdatedCount := 0
	for _, p := range configPaths {
		if updateMcpJsonFile(p, executablePath) {
			mcpUpdatedCount++
		}
	}

	greenStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	goldStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	cyanStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00E5FF"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	fmt.Println()
	fmt.Println(greenStyle.Render("  ✓ Successfully added AutoDev MCP Server to Claude Desktop!"))
	if mcpUpdatedCount > 0 {
		fmt.Printf("  ✓ Successfully configured %d Cursor Agent config file(s)!\n", mcpUpdatedCount)
	}
	fmt.Println(dimStyle.Render("  Please restart your AI client (Claude or Cursor) to load the new tools."))
	fmt.Println()
	fmt.Println(goldStyle.Render("  👉 TO CONNECT CURSOR OR COGNITION WINDSURF MANUALLY:"))
	fmt.Println("  1. Open Cursor Settings > Features > MCP.")
	fmt.Println("  2. Click '+ Add New MCP Server'.")
	fmt.Println("  3. Set Name to: " + cyanStyle.Render("autodev"))
	fmt.Println("  4. Set Type to: " + cyanStyle.Render("command"))
	fmt.Printf("  5. Set Command to: %s\n", cyanStyle.Render(fmt.Sprintf("%s mcp start", executablePath)))
	fmt.Println()

	return nil
}

// JSON-RPC basic types
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type mcpTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type mcpToolCallResult struct {
	Content []mcpTextContent `json:"content"`
	IsError bool             `json:"isError,omitempty"`
}

func sendResponse(w io.Writer, id interface{}, result interface{}) {
	resp := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(data))
}

func sendError(w io.Writer, id interface{}, code int, message string) {
	resp := jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(data))
}

func runMCPServer() error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	fmt.Fprintln(os.Stderr, "AutoDev MCP Server starting...")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var req jsonRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON-RPC: %v\n", err)
			sendError(writer, nil, -32700, "Parse error")
			continue
		}

		switch req.Method {
		case "initialize":
			result := map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "autodev-mcp",
					"version": "0.2.0",
				},
			}
			sendResponse(writer, req.ID, result)

		case "ping":
			sendResponse(writer, req.ID, map[string]interface{}{})

		case "tools/list":
			tools := []map[string]interface{}{
				{
					"name":        "autodev_scan",
					"description": "Scan the current workspace directory to detect languages, frameworks, databases, and infra config.",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"path": map[string]interface{}{
								"type":        "string",
								"description": "Optional directory path to scan. Defaults to the current directory.",
							},
						},
					},
				},
				{
					"name":        "autodev_doctor",
					"description": "Diagnose the health of common system compilers and development runtimes.",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"fix": map[string]interface{}{
								"type":        "boolean",
								"description": "Automatically attempt to auto-remediate and install missing runtimes.",
							},
						},
					},
				},
				{
					"name":        "autodev_install",
					"description": "Install a developer tool or language SDK runtime.",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"tool": map[string]interface{}{
								"type":        "string",
								"description": "Name of tool to install (nodejs, go, python, rust, docker, bun, pnpm, java, terraform, kubectl, php, ruby).",
							},
						},
						"required": []string{"tool"},
					},
				},
				{
					"name":        "autodev_audit",
					"description": "Scan codebase dependencies and find security vulnerabilities using OSV database.",
					"inputSchema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"path": map[string]interface{}{
								"type":        "string",
								"description": "Optional path to the project repository. Defaults to current directory.",
							},
						},
					},
				},
			}
			sendResponse(writer, req.ID, map[string]interface{}{
				"tools": tools,
			})

		case "tools/call":
			var callParams struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			if err := json.Unmarshal(req.Params, &callParams); err != nil {
				sendError(writer, req.ID, -32602, "Invalid params")
				continue
			}

			result := handleToolCall(callParams.Name, callParams.Arguments)
			sendResponse(writer, req.ID, result)

		default:
			if req.ID != nil {
				sendError(writer, req.ID, -32601, "Method not found: "+req.Method)
			}
		}
	}

	return nil
}

func handleToolCall(name string, args map[string]interface{}) mcpToolCallResult {
	switch name {
	case "autodev_scan":
		path := "."
		if p, ok := args["path"].(string); ok && p != "" {
			path = p
		}

		s := scanner.New(path)
		res, err := s.Scan()
		if err != nil {
			return mcpToolCallResult{
				Content: []mcpTextContent{{Type: "text", Text: fmt.Sprintf("Scan failed: %v", err)}},
				IsError: true,
			}
		}

		output, _ := json.MarshalIndent(res, "", "  ")
		return mcpToolCallResult{
			Content: []mcpTextContent{
				{Type: "text", Text: "Scan Results:\n" + string(output) + newGitHubCTAForMCP()},
			},
		}

	case "autodev_doctor":
		var fix bool
		if f, ok := args["fix"].(bool); ok {
			fix = f
		}

		diagnostics := runDoctorMCP(fix)
		return mcpToolCallResult{
			Content: []mcpTextContent{
				{Type: "text", Text: diagnostics + newGitHubCTAForMCP()},
			},
		}

	case "autodev_install":
		tool, ok := args["tool"].(string)
		if !ok || tool == "" {
			return mcpToolCallResult{
				Content: []mcpTextContent{{Type: "text", Text: "Missing 'tool' parameter"}},
				IsError: true,
			}
		}

		inst := installer.New(false)
		err := inst.Install(tool)
		if err != nil {
			return mcpToolCallResult{
				Content: []mcpTextContent{{Type: "text", Text: fmt.Sprintf("Installation failed for %s: %v", tool, err)}},
				IsError: true,
			}
		}

		return mcpToolCallResult{
			Content: []mcpTextContent{
				{Type: "text", Text: fmt.Sprintf("Successfully installed runtime: %s", tool) + newGitHubCTAForMCP()},
			},
		}

	case "autodev_audit":
		path := "."
		if p, ok := args["path"].(string); ok && p != "" {
			path = p
		}
		resList, err := scanner.AuditRepository(path)
		if err != nil {
			return mcpToolCallResult{
				Content: []mcpTextContent{{Type: "text", Text: fmt.Sprintf("Audit failed: %v", err)}},
				IsError: true,
			}
		}

		var output strings.Builder
		output.WriteString("AutoDev Supply-Chain Safety Audit\n\n")
		if len(resList) == 0 {
			output.WriteString("✓ No known security vulnerabilities found! All dependencies are safe.\n")
		} else {
			output.WriteString(fmt.Sprintf("✗ Found vulnerabilities across %d packages:\n\n", len(resList)))
			for _, res := range resList {
				output.WriteString(fmt.Sprintf("📦 %s@%s (%s)\n", res.Package.Name, res.Package.Version, res.Package.Ecosystem))
				for _, v := range res.Vulnerabilities {
					alias := ""
					if len(v.Aliases) > 0 {
						alias = " (" + v.Aliases[0] + ")"
					}
					output.WriteString(fmt.Sprintf("  - [%s] %s%s: %s\n", v.Severity, v.ID, alias, v.Summary))
				}
				output.WriteString("\n")
			}
		}
		return mcpToolCallResult{
			Content: []mcpTextContent{
				{Type: "text", Text: output.String() + newGitHubCTAForMCP()},
			},
		}

	default:
		return mcpToolCallResult{
			Content: []mcpTextContent{{Type: "text", Text: "Unknown tool name: " + name}},
			IsError: true,
		}
	}
}

func runDoctorMCP(fix bool) string {
	var report strings.Builder
	report.WriteString("AUTODEV DOCTOR - ENVIRONMENT DIAGNOSTICS\n\n")

	checkToRuntime := map[string]string{
		"Node.js":   "nodejs",
		"pnpm":      "pnpm",
		"Bun":       "bun",
		"Go":        "go",
		"Python 3":  "python",
		"Rust":      "rust",
		"Docker":    "docker",
		"kubectl":   "kubectl",
		"Terraform": "terraform",
		"Flutter":   "flutter",
		"Java":      "java",
		"PHP":       "php",
		"Ruby":      "ruby",
	}

	type mcpCheck struct {
		name string
		cmd  string
		hint string
	}

	mcpChecks := []mcpCheck{
		{name: "Git", cmd: "git", hint: "https://git-scm.com/downloads"},
		{name: "Node.js", cmd: "node", hint: "autodev install nodejs"},
		{name: "npm", cmd: "npm", hint: "Comes with Node.js"},
		{name: "pnpm", cmd: "pnpm", hint: "npm install -g pnpm"},
		{name: "yarn", cmd: "yarn", hint: "npm install -g yarn"},
		{name: "Bun", cmd: "bun", hint: "https://bun.sh"},
		{name: "Go", cmd: "go", hint: "autodev install go"},
		{name: "Python 3", cmd: "python3", hint: "autodev install python"},
		{name: "pip", cmd: "pip3", hint: "Comes with Python 3"},
		{name: "Rust", cmd: "rustc", hint: "autodev install rust"},
		{name: "Docker", cmd: "docker", hint: "autodev install docker"},
		{name: "docker compose", cmd: "docker", hint: "Upgrade Docker Desktop"},
		{name: "kubectl", cmd: "kubectl", hint: "autodev install kubectl"},
		{name: "Terraform", cmd: "terraform", hint: "autodev install terraform"},
		{name: "Flutter", cmd: "flutter", hint: "autodev install flutter"},
		{name: "Java", cmd: "java", hint: "autodev install java"},
		{name: "PHP", cmd: "php", hint: "autodev install php"},
		{name: "Ruby", cmd: "ruby", hint: "autodev install ruby"},
	}

	installed := 0
	missing := 0
	var missingRuntimes []string

	for _, c := range mcpChecks {
		_, err := exec.LookPath(c.cmd)
		if err != nil {
			report.WriteString(fmt.Sprintf("[MISSING] %-15s - %s\n", c.name, c.hint))
			missing++
			if rt, ok := checkToRuntime[c.name]; ok {
				missingRuntimes = append(missingRuntimes, rt)
			}
		} else {
			report.WriteString(fmt.Sprintf("[OK]      %-15s\n", c.name))
			installed++
		}
	}

	report.WriteString(fmt.Sprintf("\nSummary: %d installed, %d missing\n", installed, missing))

	if missing > 0 {
		if fix && len(missingRuntimes) > 0 {
			report.WriteString("\n🔧 Remediation (fix mode active):\n")
			inst := installer.New(false)
			for _, rtName := range missingRuntimes {
				rt, _ := installer.GetRuntime(rtName)
				report.WriteString(fmt.Sprintf("  → Installing %s...\n", rt.Name))
				if err := inst.Install(rtName); err != nil {
					report.WriteString(fmt.Sprintf("  ✗ Failed to install %s: %v\n", rt.Name, err))
				} else {
					report.WriteString(fmt.Sprintf("  ✓ %s installed successfully\n", rt.Name))
				}
			}
		} else {
			report.WriteString("\nHint: Call tool with fix=true to automatically repair missing dependencies.\n")
		}
	}

	return report.String()
}
