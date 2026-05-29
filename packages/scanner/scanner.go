// Package scanner implements the AutoDev repository detection engine.
// It scans a directory tree and identifies languages, frameworks,
// package managers, and infrastructure requirements.
package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Technology represents a detected technology with its confidence level.
type Technology struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"` // language | framework | package_manager | database | infra
	Version    string   `json:"version,omitempty"`
	ConfigFile string   `json:"config_file,omitempty"`
	Confidence float64  `json:"confidence"` // 0.0 – 1.0
	InstallCmd string   `json:"install_cmd,omitempty"`
	DocsURL    string   `json:"docs_url,omitempty"`
}

// ScanResult is the output of a repository scan.
type ScanResult struct {
	Path          string       `json:"path"`
	Technologies  []Technology `json:"technologies"`
	Languages     []string     `json:"languages"`
	Frameworks    []string     `json:"frameworks"`
	PackageManagers []string   `json:"package_managers"`
	Databases     []string     `json:"databases"`
	Infra         []string     `json:"infra"`
	HasDocker     bool         `json:"has_docker"`
	HasK8s        bool         `json:"has_kubernetes"`
	RecommendedSetup []string  `json:"recommended_setup"`
}

// indicator maps a filename / glob pattern to a Technology.
type indicator struct {
	file string
	tech Technology
}

// allIndicators is the master list of detection indicators.
var allIndicators = []indicator{
	// ── JavaScript / Node ecosystem ──────────────────────────────────────
	{file: "package.json", tech: Technology{Name: "Node.js", Type: "language", Confidence: 0.9, InstallCmd: "autodev install nodejs", DocsURL: "https://nodejs.org"}},
	{file: "bun.lockb", tech: Technology{Name: "Bun", Type: "package_manager", Confidence: 1.0, InstallCmd: "autodev install bun", DocsURL: "https://bun.sh"}},
	{file: "pnpm-lock.yaml", tech: Technology{Name: "pnpm", Type: "package_manager", Confidence: 1.0, InstallCmd: "npm install -g pnpm", DocsURL: "https://pnpm.io"}},
	{file: "yarn.lock", tech: Technology{Name: "yarn", Type: "package_manager", Confidence: 1.0, InstallCmd: "npm install -g yarn", DocsURL: "https://yarnpkg.com"}},
	{file: "next.config.js", tech: Technology{Name: "Next.js", Type: "framework", Confidence: 1.0, DocsURL: "https://nextjs.org"}},
	{file: "next.config.ts", tech: Technology{Name: "Next.js", Type: "framework", Confidence: 1.0, DocsURL: "https://nextjs.org"}},
	{file: "next.config.mjs", tech: Technology{Name: "Next.js", Type: "framework", Confidence: 1.0, DocsURL: "https://nextjs.org"}},
	{file: "angular.json", tech: Technology{Name: "Angular", Type: "framework", Confidence: 1.0, DocsURL: "https://angular.io"}},
	{file: "vue.config.js", tech: Technology{Name: "Vue", Type: "framework", Confidence: 1.0, DocsURL: "https://vuejs.org"}},
	{file: "svelte.config.js", tech: Technology{Name: "Svelte", Type: "framework", Confidence: 1.0, DocsURL: "https://svelte.dev"}},
	{file: "vite.config.ts", tech: Technology{Name: "Vite", Type: "framework", Confidence: 0.9, DocsURL: "https://vitejs.dev"}},
	{file: "vite.config.js", tech: Technology{Name: "Vite", Type: "framework", Confidence: 0.9, DocsURL: "https://vitejs.dev"}},

	// ── Python ───────────────────────────────────────────────────────────
	{file: "requirements.txt", tech: Technology{Name: "Python", Type: "language", Confidence: 0.9, InstallCmd: "autodev install python", DocsURL: "https://python.org"}},
	{file: "pyproject.toml", tech: Technology{Name: "Python", Type: "language", Confidence: 1.0, InstallCmd: "autodev install python", DocsURL: "https://python.org"}},
	{file: "setup.py", tech: Technology{Name: "Python", Type: "language", Confidence: 0.9, DocsURL: "https://python.org"}},
	{file: "poetry.lock", tech: Technology{Name: "Poetry", Type: "package_manager", Confidence: 1.0, InstallCmd: "curl -sSL https://install.python-poetry.org | python3 -", DocsURL: "https://python-poetry.org"}},
	{file: ".python-version", tech: Technology{Name: "pyenv", Type: "package_manager", Confidence: 1.0, InstallCmd: "autodev install pyenv", DocsURL: "https://github.com/pyenv/pyenv"}},
	{file: "Pipfile", tech: Technology{Name: "Pipenv", Type: "package_manager", Confidence: 1.0, DocsURL: "https://pipenv.pypa.io"}},
	{file: "manage.py", tech: Technology{Name: "Django", Type: "framework", Confidence: 0.9, DocsURL: "https://www.djangoproject.com"}},

	// ── Go ───────────────────────────────────────────────────────────────
	{file: "go.mod", tech: Technology{Name: "Go", Type: "language", Confidence: 1.0, InstallCmd: "autodev install go", DocsURL: "https://go.dev"}},
	{file: "go.sum", tech: Technology{Name: "Go", Type: "language", Confidence: 0.8, DocsURL: "https://go.dev"}},

	// ── Rust ─────────────────────────────────────────────────────────────
	{file: "Cargo.toml", tech: Technology{Name: "Rust", Type: "language", Confidence: 1.0, InstallCmd: "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh", DocsURL: "https://www.rust-lang.org"}},
	{file: "Cargo.lock", tech: Technology{Name: "Cargo", Type: "package_manager", Confidence: 1.0, DocsURL: "https://doc.rust-lang.org/cargo"}},

	// ── Java / Kotlin ─────────────────────────────────────────────────────
	{file: "pom.xml", tech: Technology{Name: "Maven", Type: "package_manager", Confidence: 1.0, InstallCmd: "autodev install maven", DocsURL: "https://maven.apache.org"}},
	{file: "build.gradle", tech: Technology{Name: "Gradle", Type: "package_manager", Confidence: 1.0, InstallCmd: "autodev install gradle", DocsURL: "https://gradle.org"}},
	{file: "build.gradle.kts", tech: Technology{Name: "Kotlin", Type: "language", Confidence: 1.0, InstallCmd: "autodev install kotlin", DocsURL: "https://kotlinlang.org"}},
	{file: "settings.gradle", tech: Technology{Name: "Java", Type: "language", Confidence: 0.8, InstallCmd: "autodev install java", DocsURL: "https://java.com"}},

	// ── PHP ───────────────────────────────────────────────────────────────
	{file: "composer.json", tech: Technology{Name: "PHP", Type: "language", Confidence: 0.9, InstallCmd: "autodev install php", DocsURL: "https://php.net"}},
	{file: "artisan", tech: Technology{Name: "Laravel", Type: "framework", Confidence: 1.0, DocsURL: "https://laravel.com"}},

	// ── Ruby ──────────────────────────────────────────────────────────────
	{file: "Gemfile", tech: Technology{Name: "Ruby", Type: "language", Confidence: 0.9, InstallCmd: "autodev install ruby", DocsURL: "https://ruby-lang.org"}},
	{file: "Gemfile.lock", tech: Technology{Name: "Bundler", Type: "package_manager", Confidence: 1.0, DocsURL: "https://bundler.io"}},
	{file: "config/routes.rb", tech: Technology{Name: "Ruby on Rails", Type: "framework", Confidence: 1.0, DocsURL: "https://rubyonrails.org"}},

	// ── C / C++ / .NET ────────────────────────────────────────────────────
	{file: "CMakeLists.txt", tech: Technology{Name: "C/C++", Type: "language", Confidence: 0.9, DocsURL: "https://cmake.org"}},
	{file: "Makefile", tech: Technology{Name: "Make", Type: "package_manager", Confidence: 0.5, DocsURL: "https://www.gnu.org/software/make"}},
	{file: "*.csproj", tech: Technology{Name: ".NET", Type: "language", Confidence: 1.0, InstallCmd: "autodev install dotnet", DocsURL: "https://dotnet.microsoft.com"}},
	{file: "*.sln", tech: Technology{Name: ".NET", Type: "language", Confidence: 1.0, DocsURL: "https://dotnet.microsoft.com"}},

	// ── Flutter / Dart ────────────────────────────────────────────────────
	{file: "pubspec.yaml", tech: Technology{Name: "Flutter", Type: "framework", Confidence: 0.8, InstallCmd: "autodev install flutter", DocsURL: "https://flutter.dev"}},
	{file: "pubspec.lock", tech: Technology{Name: "Dart", Type: "language", Confidence: 1.0, InstallCmd: "autodev install dart", DocsURL: "https://dart.dev"}},

	// ── Docker / Container ────────────────────────────────────────────────
	{file: "Dockerfile", tech: Technology{Name: "Docker", Type: "infra", Confidence: 1.0, InstallCmd: "autodev install docker", DocsURL: "https://docker.com"}},
	{file: "docker-compose.yml", tech: Technology{Name: "Docker Compose", Type: "infra", Confidence: 1.0, DocsURL: "https://docs.docker.com/compose"}},
	{file: "docker-compose.yaml", tech: Technology{Name: "Docker Compose", Type: "infra", Confidence: 1.0, DocsURL: "https://docs.docker.com/compose"}},

	// ── Kubernetes ────────────────────────────────────────────────────────
	{file: "k8s", tech: Technology{Name: "Kubernetes", Type: "infra", Confidence: 0.8, InstallCmd: "autodev install kubectl", DocsURL: "https://kubernetes.io"}},
	{file: "helm", tech: Technology{Name: "Helm", Type: "infra", Confidence: 0.8, DocsURL: "https://helm.sh"}},
	{file: "*.yaml", tech: Technology{Name: "Kubernetes", Type: "infra", Confidence: 0.2}}, // low confidence

	// ── Terraform ─────────────────────────────────────────────────────────
	{file: "main.tf", tech: Technology{Name: "Terraform", Type: "infra", Confidence: 1.0, InstallCmd: "autodev install terraform", DocsURL: "https://www.terraform.io"}},
	{file: "*.tf", tech: Technology{Name: "Terraform", Type: "infra", Confidence: 0.9, DocsURL: "https://www.terraform.io"}},

	// ── Databases / Services ──────────────────────────────────────────────
	{file: "firebase.json", tech: Technology{Name: "Firebase", Type: "database", Confidence: 1.0, DocsURL: "https://firebase.google.com"}},
	{file: ".firebaserc", tech: Technology{Name: "Firebase", Type: "database", Confidence: 1.0, DocsURL: "https://firebase.google.com"}},
	{file: "supabase/config.toml", tech: Technology{Name: "Supabase", Type: "database", Confidence: 1.0, DocsURL: "https://supabase.com"}},
	{file: ".supabase", tech: Technology{Name: "Supabase", Type: "database", Confidence: 0.9, DocsURL: "https://supabase.com"}},
	{file: "mongod.conf", tech: Technology{Name: "MongoDB", Type: "database", Confidence: 1.0, DocsURL: "https://mongodb.com"}},
	{file: "nginx.conf", tech: Technology{Name: "Nginx", Type: "infra", Confidence: 1.0, DocsURL: "https://nginx.org"}},
	{file: "httpd.conf", tech: Technology{Name: "Apache", Type: "infra", Confidence: 1.0, DocsURL: "https://httpd.apache.org"}},
	{file: ".ruby-version", tech: Technology{Name: "rbenv", Type: "package_manager", Confidence: 1.0, DocsURL: "https://github.com/rbenv/rbenv"}},
	{file: ".nvmrc", tech: Technology{Name: "nvm", Type: "package_manager", Confidence: 1.0, InstallCmd: "autodev install nvm", DocsURL: "https://github.com/nvm-sh/nvm"}},
	{file: ".tool-versions", tech: Technology{Name: "asdf", Type: "package_manager", Confidence: 1.0, DocsURL: "https://asdf-vm.com"}},
}

// Scanner scans a directory for technologies.
type Scanner struct {
	RootPath string
}

// New creates a new Scanner for the given directory.
func New(path string) *Scanner {
	abs, _ := filepath.Abs(path)
	return &Scanner{RootPath: abs}
}

// Scan walks the directory tree and returns a ScanResult.
func (s *Scanner) Scan() (*ScanResult, error) {
	result := &ScanResult{Path: s.RootPath}
	seen := map[string]bool{} // deduplicate

	err := filepath.WalkDir(s.RootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}
		// Skip hidden dirs and common noise
		base := d.Name()
		if d.IsDir() && (strings.HasPrefix(base, ".") ||
			base == "node_modules" || base == "vendor" ||
			base == ".git" || base == "dist" || base == "build" ||
			base == "__pycache__" || base == ".next") {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(s.RootPath, path)
		for _, ind := range allIndicators {
			if matchIndicator(rel, ind.file) {
				key := ind.tech.Name + "|" + ind.tech.Type
				if seen[key] {
					continue
				}
				tech := ind.tech
				tech.ConfigFile = rel
				result.Technologies = append(result.Technologies, tech)
				seen[key] = true
			}
		}

		// Check package.json for deeper React/framework detection
		if base == "package.json" {
			detectFromPackageJSON(path, result, seen)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Separate into categories
	for _, t := range result.Technologies {
		switch t.Type {
		case "language":
			result.Languages = appendUniq(result.Languages, t.Name)
		case "framework":
			result.Frameworks = appendUniq(result.Frameworks, t.Name)
		case "package_manager":
			result.PackageManagers = appendUniq(result.PackageManagers, t.Name)
		case "database":
			result.Databases = appendUniq(result.Databases, t.Name)
		case "infra":
			result.Infra = appendUniq(result.Infra, t.Name)
		}
		if t.Name == "Docker" || t.Name == "Docker Compose" {
			result.HasDocker = true
		}
		if t.Name == "Kubernetes" {
			result.HasK8s = true
		}
	}

	result.RecommendedSetup = buildSetupPlan(result)
	return result, nil
}

// matchIndicator returns true if the file path matches the indicator pattern.
func matchIndicator(filePath, pattern string) bool {
	base := filepath.Base(filePath)
	if strings.ContainsAny(pattern, "*?") {
		matched, _ := filepath.Match(pattern, base)
		return matched
	}
	// Exact basename or subpath match
	return base == pattern || strings.HasSuffix(filePath, pattern)
}

// detectFromPackageJSON reads package.json and detects React, Angular, Vue, etc.
func detectFromPackageJSON(path string, result *ScanResult, seen map[string]bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var pkg map[string]json.RawMessage
	if err := json.Unmarshal(data, &pkg); err != nil {
		return
	}

	allDeps := map[string]json.RawMessage{}
	for _, key := range []string{"dependencies", "devDependencies", "peerDependencies"} {
		if raw, ok := pkg[key]; ok {
			var deps map[string]json.RawMessage
			if err := json.Unmarshal(raw, &deps); err == nil {
				for k, v := range deps {
					allDeps[k] = v
				}
			}
		}
	}

	frameworkDeps := map[string]Technology{
		"react":        {Name: "React", Type: "framework", Confidence: 1.0, DocsURL: "https://react.dev"},
		"vue":          {Name: "Vue", Type: "framework", Confidence: 1.0, DocsURL: "https://vuejs.org"},
		"@angular/core": {Name: "Angular", Type: "framework", Confidence: 1.0, DocsURL: "https://angular.io"},
		"svelte":       {Name: "Svelte", Type: "framework", Confidence: 1.0, DocsURL: "https://svelte.dev"},
		"next":         {Name: "Next.js", Type: "framework", Confidence: 1.0, DocsURL: "https://nextjs.org"},
		"express":      {Name: "Express", Type: "framework", Confidence: 1.0, DocsURL: "https://expressjs.com"},
		"fastify":      {Name: "Fastify", Type: "framework", Confidence: 1.0, DocsURL: "https://fastify.dev"},
		"typescript":   {Name: "TypeScript", Type: "language", Confidence: 1.0, DocsURL: "https://typescriptlang.org"},
	}

	for dep, tech := range frameworkDeps {
		if _, ok := allDeps[dep]; ok {
			key := tech.Name + "|" + tech.Type
			if !seen[key] {
				result.Technologies = append(result.Technologies, tech)
				seen[key] = true
			}
		}
	}
}

// buildSetupPlan generates an ordered list of install commands.
func buildSetupPlan(result *ScanResult) []string {
	var plan []string
	for _, t := range result.Technologies {
		if t.InstallCmd != "" {
			plan = append(plan, t.InstallCmd)
		}
	}
	return plan
}

func appendUniq(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}
