package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/autodev-sh/autodev/installer"
	"github.com/autodev-sh/autodev/scanner"
	"github.com/spf13/cobra"
)

func newGitHubCTAForMCP() string {
	return "\n⭐ Love AutoDev? Star the repo: https://github.com/HEETMEHTA18/autodev\n"
}

func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start a Model Context Protocol (MCP) server",
		Long:  `Run a native Model Context Protocol (MCP) server over stdin/stdout, allowing AI coding tools like Claude Desktop or Cursor to interface directly with your dev environment.`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start the MCP server over stdin/stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	})

	return cmd
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
	
	// Print start status to stderr so it doesn't pollute stdout json stream
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
					"version": "0.1.6",
				},
			}
			sendResponse(writer, req.ID, result)

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
