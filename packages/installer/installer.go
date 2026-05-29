// Package installer handles runtime and dependency installation.
package installer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Runtime represents a installable runtime or tool.
type Runtime struct {
	Name        string
	CheckCmd    string   // command to check if installed (e.g. "node --version")
	LinuxCmd    []string // install commands for Linux
	MacCmd      []string // install commands for macOS
	WindowsCmd  []string // install commands for Windows
	PostInstall string   // shell command to run after install
	Version     string   // preferred version (empty = latest)
}

// Status represents the installation status of a runtime.
type Status struct {
	Runtime   Runtime
	Installed bool
	Version   string
	Error     error
}

// Installer manages runtime installations.
type Installer struct {
	DryRun  bool
	Verbose bool
}

// New creates a new Installer.
func New(dryRun bool) *Installer {
	return &Installer{DryRun: dryRun}
}

// runtimes is the built-in catalog of supported runtimes.
var runtimes = map[string]Runtime{
	"nodejs": {
		Name:       "Node.js",
		CheckCmd:   "node --version",
		LinuxCmd:   []string{"curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -", "sudo apt-get install -y nodejs"},
		MacCmd:     []string{"brew install node@22"},
		WindowsCmd: []string{"winget install OpenJS.NodeJS.LTS"},
	},
	"go": {
		Name:       "Go",
		CheckCmd:   "go version",
		LinuxCmd:   []string{"curl -fsSL https://go.dev/dl/go1.22.2.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -", "echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc"},
		MacCmd:     []string{"brew install go"},
		WindowsCmd: []string{"winget install GoLang.Go"},
	},
	"python": {
		Name:       "Python",
		CheckCmd:   "python3 --version",
		LinuxCmd:   []string{"sudo apt-get install -y python3 python3-pip python3-venv"},
		MacCmd:     []string{"brew install python@3.12"},
		WindowsCmd: []string{"winget install Python.Python.3"},
	},
	"rust": {
		Name:       "Rust",
		CheckCmd:   "rustc --version",
		LinuxCmd:   []string{"curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y"},
		MacCmd:     []string{"curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y"},
		WindowsCmd: []string{"winget install Rustlang.Rust.MSVC"},
	},
	"docker": {
		Name:       "Docker",
		CheckCmd:   "docker --version",
		LinuxCmd:   []string{"curl -fsSL https://get.docker.com | sh", "sudo usermod -aG docker $USER"},
		MacCmd:     []string{"brew install --cask docker"},
		WindowsCmd: []string{"winget install Docker.DockerDesktop"},
	},
	"bun": {
		Name:       "Bun",
		CheckCmd:   "bun --version",
		LinuxCmd:   []string{"curl -fsSL https://bun.sh/install | bash"},
		MacCmd:     []string{"brew install oven-sh/bun/bun"},
		WindowsCmd: []string{"powershell -c \"irm bun.sh/install.ps1 | iex\""},
	},
	"pnpm": {
		Name:       "pnpm",
		CheckCmd:   "pnpm --version",
		LinuxCmd:   []string{"npm install -g pnpm"},
		MacCmd:     []string{"npm install -g pnpm"},
		WindowsCmd: []string{"npm install -g pnpm"},
	},
	"java": {
		Name:       "Java (OpenJDK 21)",
		CheckCmd:   "java -version",
		LinuxCmd:   []string{"sudo apt-get install -y openjdk-21-jdk"},
		MacCmd:     []string{"brew install openjdk@21"},
		WindowsCmd: []string{"winget install Microsoft.OpenJDK.21"},
	},
	"kotlin": {
		Name:       "Kotlin",
		CheckCmd:   "kotlin -version",
		LinuxCmd:   []string{"sudo apt-get install -y kotlin"},
		MacCmd:     []string{"brew install kotlin"},
		WindowsCmd: []string{"winget install JetBrains.Kotlin"},
	},
	"flutter": {
		Name:       "Flutter SDK",
		CheckCmd:   "flutter --version",
		LinuxCmd:   []string{"sudo snap install flutter --classic"},
		MacCmd:     []string{"brew install --cask flutter"},
		WindowsCmd: []string{"winget install Flutter.Flutter"},
	},
	"terraform": {
		Name:       "Terraform",
		CheckCmd:   "terraform version",
		LinuxCmd:   []string{"sudo apt-get install -y gnupg software-properties-common", "wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg", "echo \"deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main\" | sudo tee /etc/apt/sources.list.d/hashicorp.list", "sudo apt update && sudo apt-get install terraform"},
		MacCmd:     []string{"brew install terraform"},
		WindowsCmd: []string{"winget install Hashicorp.Terraform"},
	},
	"kubectl": {
		Name:       "kubectl",
		CheckCmd:   "kubectl version --client",
		LinuxCmd:   []string{"curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl", "sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl"},
		MacCmd:     []string{"brew install kubectl"},
		WindowsCmd: []string{"winget install Kubernetes.kubectl"},
	},
	"php": {
		Name:       "PHP",
		CheckCmd:   "php --version",
		LinuxCmd:   []string{"sudo apt-get install -y php php-cli php-fpm php-mysql php-zip php-gd php-mbstring php-curl php-xml php-bcmath"},
		MacCmd:     []string{"brew install php"},
		WindowsCmd: []string{"winget install PHP.PHP"},
	},
	"ruby": {
		Name:       "Ruby",
		CheckCmd:   "ruby --version",
		LinuxCmd:   []string{"sudo apt-get install -y ruby-full"},
		MacCmd:     []string{"brew install ruby"},
		WindowsCmd: []string{"winget install RubyInstallerTeam.Ruby"},
	},
}

// CheckStatus checks whether a runtime is installed and returns its version.
func (i *Installer) CheckStatus(name string) Status {
	rt, ok := runtimes[name]
	if !ok {
		return Status{Error: fmt.Errorf("unknown runtime: %s", name)}
	}

	parts := strings.Fields(rt.CheckCmd)
	out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	if err != nil {
		return Status{Runtime: rt, Installed: false}
	}
	return Status{Runtime: rt, Installed: true, Version: strings.TrimSpace(string(out))}
}

// Install installs the given runtime using the platform-appropriate command.
func (i *Installer) Install(name string) error {
	rt, ok := runtimes[name]
	if !ok {
		return fmt.Errorf("unknown runtime: %s", name)
	}

	var cmds []string
	switch runtime.GOOS {
	case "linux":
		cmds = rt.LinuxCmd
	case "darwin":
		cmds = rt.MacCmd
	case "windows":
		cmds = rt.WindowsCmd
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	for _, cmdStr := range cmds {
		if i.DryRun {
			fmt.Printf("[dry-run] Would run: %s\n", cmdStr)
			continue
		}

		fmt.Printf("  → Running: %s\n", cmdStr)
		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("command failed (%q): %w", cmdStr, err)
		}
	}

	return nil
}

// AllRuntimeNames returns a sorted list of known runtime names.
func AllRuntimeNames() []string {
	names := make([]string, 0, len(runtimes))
	for k := range runtimes {
		names = append(names, k)
	}
	return names
}

// GetRuntime returns a runtime by name.
func GetRuntime(name string) (Runtime, bool) {
	rt, ok := runtimes[name]
	return rt, ok
}
