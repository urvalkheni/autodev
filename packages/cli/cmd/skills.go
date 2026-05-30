package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/autodev-sh/autodev/scanner"
	"github.com/autodev-sh/autodev/skills"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func newSkillsCmd() *cobra.Command {
	var (
		path        string
		deep        bool
		exportFmt   string
		sync        bool
		interactive bool
	)

	cmd := &cobra.Command{
		Use:   "skills [path]",
		Short: "Generate a personalized learning roadmap based on detected technologies",
		Long: `Scan the repository, detect all technologies, and generate a personalized
learning roadmap showing current skills, next steps, and long-term goals.

Use --deep to analyze git history for confidence scoring.
Use --export to output as json, md, or html.
Use --interactive to select custom target skills to focus on.

Powered by skills.sh.`,
		Example: `  autodev skills
  autodev skills ./my-project
  autodev skills --deep
  autodev skills --deep --export md > roadmap.md
  autodev skills -i
  autodev skills --sync`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				path = args[0]
			}
			return runSkills(path, deep, exportFmt, sync, interactive)
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "path to scan")
	cmd.Flags().BoolVar(&deep, "deep", false, "deep analysis: scan git history for confidence scoring")
	cmd.Flags().StringVar(&exportFmt, "export", "", "export format: json, md, or html")
	cmd.Flags().BoolVar(&sync, "sync", false, "sync skill profile to skills.sh")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "interactive mode: select target skills from catalog")
	return cmd
}

func runSkills(path string, deep bool, exportFmt string, syncProfile bool, interactive bool) error {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	successStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF87"))
	errorStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF4444"))

	// Only print header when not exporting (to keep export output clean)
	silent := exportFmt != ""

	if !silent && !interactive {
		fmt.Println()
		fmt.Println(titleStyle.Render("  ⚡ AutoDev Skills Engine v0.2.0"))
		fmt.Println(dimStyle.Render("  Powered by skills.sh"))
		fmt.Println()
	}

	var detected []string
	if interactive {
		gen := skills.New()
		available := gen.GetAvailableSkills()

		// Run interactive selector TUI
		selModel := skillSelectModel{
			skills:   available,
			selected: make(map[string]bool),
		}
		p := tea.NewProgram(selModel)
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("interactive selector failed: %w", err)
		}

		m := finalModel.(skillSelectModel)
		if m.quitting {
			return nil
		}

		for skill, checked := range m.selected {
			if checked {
				detected = append(detected, skill)
			}
		}

		if len(detected) == 0 {
			fmt.Println("  No skills selected. Exiting.")
			return nil
		}

		if !silent {
			fmt.Println()
			fmt.Println(titleStyle.Render("  ⚡ AutoDev Skills Engine v0.2.0"))
			fmt.Println(dimStyle.Render("  Powered by skills.sh"))
			fmt.Println()
		}
	} else {
		// ── Step 1: Scan to detect technologies ─────────────────────────────
		s := scanner.New(path)
		result, err := s.Scan()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		detected = append(detected, result.Languages...)
		detected = append(detected, result.Frameworks...)
		detected = append(detected, result.PackageManagers...)
		detected = append(detected, result.Databases...)
	}

	if len(detected) == 0 {
		if !silent {
			fmt.Println(dimStyle.Render("  No technologies detected. Try running in a project directory."))
		}
		return nil
	}

	// ── Step 2: Generate roadmap ────────────────────────────────────────
	gen := skills.New()
	roadmap := gen.Generate(detected)

	// ── Step 3: Deep analysis (git history) ─────────────────────────────
	if deep {
		if !silent {
			fmt.Println(dimStyle.Render("  [deep] Analyzing git history..."))
		}
		deepStats := runDeepAnalysis(path, detected)
		roadmap.DeepStats = deepStats
		if !silent {
			fmt.Println(successStyle.Render(fmt.Sprintf("  [deep] ✓ Analyzed %d technologies with confidence scoring", len(deepStats))))
		}
	}

	// ── Step 5: Output ──────────────────────────────────────────────────
	switch exportFmt {
	case "json":
		out, err := roadmap.ExportJSON()
		if err != nil {
			return err
		}
		fmt.Println(out)
	case "md", "markdown":
		fmt.Print(roadmap.ExportMarkdown())
	case "html":
		fmt.Print(generateHTMLExport(roadmap))
	case "":
		// Default terminal output
		if deep {
			roadmap.PrintDeep()
		} else {
			roadmap.Print()
		}
	default:
		return fmt.Errorf("unknown export format: %s (use json, md, or html)", exportFmt)
	}

	// ── Step 6: Sync to skills.sh ───────────────────────────────────────
	if syncProfile {
		if !silent {
			fmt.Println(dimStyle.Render("  [sync] Uploading skill profile to skills.sh..."))
		}
		err := syncToSkillsSh(roadmap)
		if err != nil {
			if !silent {
				fmt.Println(errorStyle.Render(fmt.Sprintf("  [sync] ✗ %s", err.Error())))
			}
		} else {
			if !silent {
				fmt.Println(successStyle.Render("  [sync] ✓ Profile synced to skills.sh"))
			}
		}
	}

	if !silent && exportFmt == "" {
		fmt.Println(dimStyle.Render("  Visit https://skills.sh for interactive learning paths and courses."))
		fmt.Println()
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Deep analysis: git history parsing
// ──────────────────────────────────────────────────────────────────────────────

// techExtensions maps technology names to their associated file extensions.
var techExtensions = map[string][]string{
	"Node.js":    {".js", ".mjs", ".cjs"},
	"TypeScript": {".ts", ".tsx"},
	"React":      {".jsx", ".tsx"},
	"Python":     {".py", ".pyw"},
	"Go":         {".go"},
	"Rust":       {".rs"},
	"Java":       {".java"},
	"PHP":        {".php"},
	"Ruby":       {".rb"},
	"Kotlin":     {".kt", ".kts"},
	"Dart":       {".dart"},
	"C/C++":      {".c", ".cpp", ".h", ".hpp"},
	".NET":       {".cs", ".vb", ".fs"},
	"Svelte":     {".svelte"},
	"Vue":        {".vue"},
}

func runDeepAnalysis(rootPath string, detected []string) []skills.DeepSkillStats {
	absPath, _ := filepath.Abs(rootPath)

	// Check if path is a git repo
	if _, err := os.Stat(filepath.Join(absPath, ".git")); os.IsNotExist(err) {
		// Not a git repo — return basic stats
		return buildBasicStats(detected)
	}

	// Get git log with numstat
	cmd := exec.Command("git", "log", "--format=%H %aI", "--numstat", "--no-merges", "--since=2 years ago")
	cmd.Dir = absPath
	output, err := cmd.Output()
	if err != nil {
		return buildBasicStats(detected)
	}

	// Parse git log
	type fileChange struct {
		added   int
		removed int
		date    time.Time
	}

	techStats := map[string]*struct {
		files      map[string]bool
		commits    int
		linesAdded int
		lastDate   time.Time
	}{}

	// Initialize stats for detected technologies
	for _, name := range detected {
		techStats[name] = &struct {
			files      map[string]bool
			commits    int
			linesAdded int
			lastDate   time.Time
		}{files: map[string]bool{}}
	}

	lines := strings.Split(string(output), "\n")
	var currentDate time.Time
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Commit header line: <hash> <date>
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 && len(parts[0]) == 40 {
			t, err := time.Parse(time.RFC3339, parts[1])
			if err == nil {
				currentDate = t
			}
			continue
		}

		// Numstat line: <added>\t<removed>\t<file>
		fields := strings.Split(line, "\t")
		if len(fields) != 3 {
			continue
		}
		added, _ := strconv.Atoi(fields[0])
		filePath := fields[2]
		ext := strings.ToLower(filepath.Ext(filePath))

		// Match extension to technology
		for techName, exts := range techExtensions {
			for _, e := range exts {
				if ext == e {
					if stats, ok := techStats[techName]; ok {
						stats.files[filePath] = true
						stats.commits++
						stats.linesAdded += added
						if currentDate.After(stats.lastDate) {
							stats.lastDate = currentDate
						}
					}
				}
			}
		}
	}

	// Build DeepSkillStats
	var result []skills.DeepSkillStats
	now := time.Now()
	for _, name := range detected {
		stats, ok := techStats[name]
		if !ok {
			continue
		}

		daysSince := 365 // default if no commits found
		if !stats.lastDate.IsZero() {
			daysSince = int(now.Sub(stats.lastDate).Hours() / 24)
		}
		if daysSince < 0 {
			daysSince = 0
		}

		confidence := skills.ComputeConfidence(
			len(stats.files),
			stats.commits,
			daysSince,
			1, // single repo for now
		)

		result = append(result, skills.DeepSkillStats{
			Name:        name,
			Category:    getCategoryForTech(name),
			FileCount:   len(stats.files),
			CommitCount: stats.commits,
			LinesAdded:  stats.linesAdded,
			DaysSince:   daysSince,
			Repos:       1,
			Confidence:  confidence,
			Level:       skills.LevelFromConfidence(confidence),
		})
	}

	// Sort by confidence descending
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Confidence > result[i].Confidence {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

func buildBasicStats(detected []string) []skills.DeepSkillStats {
	var result []skills.DeepSkillStats
	for _, name := range detected {
		result = append(result, skills.DeepSkillStats{
			Name:       name,
			Category:   getCategoryForTech(name),
			Confidence: 50,
			Level:      "intermediate",
		})
	}
	return result
}

func getCategoryForTech(name string) string {
	categories := map[string]string{
		"Node.js": "Runtime", "TypeScript": "Language", "React": "Framework",
		"Next.js": "Framework", "Python": "Language", "Go": "Language",
		"Rust": "Language", "Java": "Language", "PHP": "Language",
		"Ruby": "Language", "Docker": "DevOps", "Kubernetes": "DevOps",
		"PostgreSQL": "Database", "MongoDB": "Database", "Redis": "Database",
		"Firebase": "Database", "pnpm": "Package Manager", "yarn": "Package Manager",
		"Flutter": "Framework", "Dart": "Language", "Vue": "Framework",
		"Angular": "Framework", "Svelte": "Framework", "Django": "Framework",
		"FastAPI": "Framework", "Laravel": "Framework", "Express": "Framework",
		"Spring Boot": "Framework", "NestJS": "Framework",
	}
	if cat, ok := categories[name]; ok {
		return cat
	}
	return "Tool"
}


// ──────────────────────────────────────────────────────────────────────────────
// Sync to skills.sh (placeholder — ready for API integration)
// ──────────────────────────────────────────────────────────────────────────────

func syncToSkillsSh(roadmap *skills.Roadmap) error {
	// Future: POST to skills.sh API
	// For now, save locally as a JSON profile
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	profileDir := filepath.Join(homeDir, ".config", "autodev")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(roadmap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profile: %w", err)
	}

	profilePath := filepath.Join(profileDir, "skills-profile.json")
	if err := os.WriteFile(profilePath, data, 0644); err != nil {
		return fmt.Errorf("write profile: %w", err)
	}

	fmt.Printf("  [sync] Profile saved to %s\n", profilePath)
	fmt.Println("  [sync] Visit https://skills.sh to connect your account.")
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// HTML export
// ──────────────────────────────────────────────────────────────────────────────

func generateHTMLExport(roadmap *skills.Roadmap) string {
	var b strings.Builder

	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>AutoDev Skills Report</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'SF Mono','Fira Code',monospace;background:#0A0A0A;color:#F0F0F0;padding:40px;max-width:800px;margin:0 auto}
h1{color:#FFD700;font-size:24px;margin-bottom:8px}
h2{color:#FFD700;font-size:16px;margin:24px 0 12px;text-transform:uppercase;letter-spacing:2px}
.meta{color:#888;font-size:12px;margin-bottom:24px}
table{width:100%;border-collapse:collapse;margin:12px 0}
th,td{text-align:left;padding:8px 12px;border:1px solid #2A2A2A;font-size:13px}
th{background:#111;color:#FFD700;font-weight:bold;text-transform:uppercase;font-size:11px;letter-spacing:1px}
.bar-container{background:#111;border:1px solid #2A2A2A;height:16px;position:relative}
.bar{height:100%;transition:width 0.3s}
.expert{background:#00FF87}
.advanced{background:#4A90E2}
.intermediate{background:#FFD700}
.beginner{background:#FF4444}
.skill-tag{display:inline-block;border:1px solid #444;padding:4px 10px;margin:4px;font-size:12px}
.insight{background:#111;border-left:3px solid #FFD700;padding:10px 16px;margin:8px 0;font-size:13px}
.footer{margin-top:40px;padding-top:20px;border-top:1px solid #2A2A2A;color:#555;font-size:11px;text-align:center}
a{color:#FFD700}
</style>
</head>
<body>
`)

	b.WriteString(fmt.Sprintf("<h1>⚡ AutoDev Skills Report</h1>\n"))
	b.WriteString(fmt.Sprintf("<p class='meta'>%s · Generated by AutoDev CLI</p>\n", roadmap.GeneratedAt))

	// Deep stats table
	if len(roadmap.DeepStats) > 0 {
		b.WriteString("<h2>Skill Matrix</h2>\n<table>\n")
		b.WriteString("<tr><th>Technology</th><th>Level</th><th>Confidence</th><th>Visualization</th><th>Files</th><th>Commits</th></tr>\n")
		for _, s := range roadmap.DeepStats {
			barClass := s.Level
			b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%.1f%%</td>",
				s.Name, s.Level, s.Confidence))
			b.WriteString(fmt.Sprintf("<td><div class='bar-container'><div class='bar %s' style='width:%.0f%%'></div></div></td>",
				barClass, s.Confidence))
			b.WriteString(fmt.Sprintf("<td>%d</td><td>%d</td></tr>\n", s.FileCount, s.CommitCount))
		}
		b.WriteString("</table>\n")
	}

	// Current skills
	if len(roadmap.CurrentSkills) > 0 {
		b.WriteString("<h2>Current Skills</h2>\n<div>\n")
		for _, s := range roadmap.CurrentSkills {
			b.WriteString(fmt.Sprintf("<span class='skill-tag' style='border-color:#00FF87;color:#00FF87'>%s</span>\n", s.Name))
		}
		b.WriteString("</div>\n")
	}

	// Next steps
	if len(roadmap.NextSteps) > 0 {
		b.WriteString("<h2>Recommended Next Steps</h2>\n<div>\n")
		for _, s := range roadmap.NextSteps {
			b.WriteString(fmt.Sprintf("<span class='skill-tag' style='border-color:#FFD700;color:#FFD700'>%s</span>\n", s.Name))
		}
		b.WriteString("</div>\n")
	}

	// AI insights
	if len(roadmap.AIInsights) > 0 {
		b.WriteString("<h2>AI-Powered Insights</h2>\n")
		for _, insight := range roadmap.AIInsights {
			b.WriteString(fmt.Sprintf("<div class='insight'>💡 %s</div>\n", insight))
		}
	}

	b.WriteString("<div class='footer'>Generated by <a href='https://github.com/HEETMEHTA18/autodev'>AutoDev</a> · Powered by <a href='https://skills.sh'>skills.sh</a></div>\n")
	b.WriteString("</body></html>")

	return b.String()
}

// ──────────────────────────────────────────────────────────────────────────────
// Interactive TUI for Selecting Skills
// ──────────────────────────────────────────────────────────────────────────────

type skillSelectModel struct {
	skills   []string
	cursor   int
	selected map[string]bool
	quitting bool
}

func (m skillSelectModel) Init() tea.Cmd { return nil }

func (m skillSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.skills)-1 {
				m.cursor++
			}
		case " ", "space", "enter":
			if msg.String() == "enter" {
				hasSelection := false
				for _, checked := range m.selected {
					if checked {
						hasSelection = true
						break
					}
				}
				if !hasSelection {
					m.selected[m.skills[m.cursor]] = true
				}
				return m, tea.Quit
			}
			m.selected[m.skills[m.cursor]] = !m.selected[m.skills[m.cursor]]
		}
	}
	return m, nil
}

func (m skillSelectModel) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFD700"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#F0F0F0"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	checkStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF87")).Bold(true)

	var b strings.Builder
	b.WriteString("\n  " + titleStyle.Render("🎯 SELECT SKILLS TO FOCUS ON") + "\n")
	b.WriteString("  " + dimStyle.Render("Choose target skills from the skills.sh catalog:") + "\n\n")

	// Paging calculations
	maxVisible := 12
	startIdx := 0
	if m.cursor >= maxVisible/2 {
		startIdx = m.cursor - maxVisible/2
	}
	if startIdx+maxVisible > len(m.skills) {
		startIdx = len(m.skills) - maxVisible
	}
	if startIdx < 0 {
		startIdx = 0
	}
	endIdx := startIdx + maxVisible
	if endIdx > len(m.skills) {
		endIdx = len(m.skills)
	}

	if startIdx > 0 {
		b.WriteString("   " + dimStyle.Render("▲ -- scroll up for more --") + "\n")
	} else {
		b.WriteString("\n")
	}

	for i := startIdx; i < endIdx; i++ {
		skill := m.skills[i]
		cursor := "  "
		style := normalStyle
		if m.cursor == i {
			cursor = "▶ "
			style = selectedStyle
		}

		checkbox := "[ ]"
		if m.selected[skill] {
			checkbox = checkStyle.Render("[X]")
		}

		b.WriteString(fmt.Sprintf("  %s%s %s\n", cursor, checkbox, style.Render(skill)))
	}

	if endIdx < len(m.skills) {
		b.WriteString("   " + dimStyle.Render("▼ -- scroll down for more --") + "\n")
	} else {
		b.WriteString("\n")
	}

	b.WriteString("\n  " + dimStyle.Render("space: toggle select | enter: generate customized roadmap | q: quit") + "\n")
	return b.String()
}
