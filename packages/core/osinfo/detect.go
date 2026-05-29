// Package osinfo provides OS, architecture, and hardware detection.
package osinfo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Info holds detected system information.
type Info struct {
	OS             string
	Arch           string
	Version        string
	PackageManager string
	RAMBytes       uint64
	DiskBytes      uint64
	CPUCores       int
	HasGPU         bool
}

// Detect returns a populated Info struct for the current system.
func Detect() (*Info, error) {
	info := &Info{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCores: runtime.NumCPU(),
	}

	info.Version = detectOSVersion()
	info.PackageManager = detectPackageManager(info.OS)
	info.RAMBytes = detectRAM()
	info.DiskBytes = detectDisk()
	info.HasGPU = detectGPU()

	return info, nil
}

func detectOSVersion() string {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/etc/os-release"); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					return strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), `"`)
				}
			}
		}
		return "Linux"
	case "darwin":
		if out, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
			return "macOS " + strings.TrimSpace(string(out))
		}
		return "macOS"
	case "windows":
		if out, err := exec.Command("cmd", "/c", "ver").Output(); err == nil {
			return strings.TrimSpace(string(out))
		}
		return "Windows"
	default:
		return runtime.GOOS
	}
}

func detectPackageManager(os string) string {
	switch os {
	case "darwin":
		if commandExists("brew") {
			return "homebrew"
		}
		return "none"
	case "linux":
		for _, pm := range []string{"apt", "dnf", "yum", "pacman", "zypper", "apk"} {
			if commandExists(pm) {
				return pm
			}
		}
		return "none"
	case "windows":
		if commandExists("winget") {
			return "winget"
		}
		if commandExists("choco") {
			return "chocolatey"
		}
		if commandExists("scoop") {
			return "scoop"
		}
		return "none"
	}
	return "none"
}

func detectRAM() uint64 {
	switch runtime.GOOS {
	case "linux":
		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					if len(fields) >= 2 {
						kb, _ := strconv.ParseUint(fields[1], 10, 64)
						return kb * 1024
					}
				}
			}
		}
	case "darwin":
		if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
			bytes, _ := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
			return bytes
		}
	case "windows":
		if out, err := exec.Command("wmic", "ComputerSystem", "get", "TotalPhysicalMemory", "/Value").Output(); err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				if strings.HasPrefix(line, "TotalPhysicalMemory=") {
					val := strings.TrimPrefix(line, "TotalPhysicalMemory=")
					bytes, _ := strconv.ParseUint(strings.TrimSpace(val), 10, 64)
					return bytes
				}
			}
		}
	}
	return 0
}

func detectDisk() uint64 {
	// Simplified: report available disk on current partition
	// A full implementation would use syscall.Statfs
	return 0
}

func detectGPU() bool {
	switch runtime.GOOS {
	case "linux":
		return commandExists("nvidia-smi") || fileExists("/dev/dri")
	case "darwin":
		// All modern Macs have a GPU
		return true
	case "windows":
		out, err := exec.Command("wmic", "path", "win32_VideoController", "get", "name").Output()
		return err == nil && len(out) > 10
	}
	return false
}

// FormatRAM returns a human-readable RAM string.
func FormatRAM(bytes uint64) string {
	if bytes == 0 {
		return "Unknown"
	}
	gb := float64(bytes) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.1f GB", gb)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
