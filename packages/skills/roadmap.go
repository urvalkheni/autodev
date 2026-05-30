// Package skills provides learning roadmap generation, deep git analysis,
// confidence scoring, multi-format export, and AI-powered recommendations.
package skills

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// ──────────────────────────────────────────────────────────────────────────────
// Data structures
// ──────────────────────────────────────────────────────────────────────────────

// Skill represents a technology skill in a roadmap.
type Skill struct {
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	Level      string   `json:"level"` // beginner | intermediate | advanced | expert
	Resources  []string `json:"resources"`
	NextSkills []string `json:"next_skills"`
}

// DeepSkillStats holds git-derived metrics for a single technology.
type DeepSkillStats struct {
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	FileCount   int     `json:"file_count"`
	CommitCount int     `json:"commit_count"`
	LinesAdded  int     `json:"lines_added"`
	DaysSince   int     `json:"days_since_last_commit"`
	Repos       int     `json:"repo_count"`
	Confidence  float64 `json:"confidence"` // 0–100
	Level       string  `json:"level"`      // computed from confidence
}

// Roadmap represents a personalized learning roadmap.
type Roadmap struct {
	Title         string           `json:"title"`
	GeneratedAt   string           `json:"generated_at"`
	CurrentSkills []Skill          `json:"current_skills"`
	NextSteps     []Skill          `json:"next_steps"`
	LongTermGoals []Skill          `json:"long_term_goals"`
	DeepStats     []DeepSkillStats `json:"deep_stats,omitempty"`
	AIInsights    []string         `json:"ai_insights,omitempty"`
}

// RoadmapGenerator generates learning roadmaps from detected technologies.
type RoadmapGenerator struct {
	catalog map[string]Skill
}

// New creates a new RoadmapGenerator with the built-in skill catalog.
func New() *RoadmapGenerator {
	return &RoadmapGenerator{catalog: buildCatalog()}
}

// GetAvailableSkills returns a sorted list of all skill names in the catalog.
func (r *RoadmapGenerator) GetAvailableSkills() []string {
	var keys []string
	for k := range r.catalog {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ──────────────────────────────────────────────────────────────────────────────
// Roadmap generation
// ──────────────────────────────────────────────────────────────────────────────

// Generate produces a personalized roadmap from a list of detected technology names.
func (r *RoadmapGenerator) Generate(detected []string) *Roadmap {
	roadmap := &Roadmap{
		Title:       fmt.Sprintf("Personalized Roadmap — %s", strings.Join(detected, ", ")),
		GeneratedAt: time.Now().Format(time.RFC3339),
	}

	seen := map[string]bool{}
	for _, name := range detected {
		seen[name] = true
	}

	// Current skills from detected technologies
	for _, name := range detected {
		if skill, ok := r.catalog[name]; ok {
			roadmap.CurrentSkills = append(roadmap.CurrentSkills, skill)
		}
	}

	// Next steps: direct successors of current skills
	nextSeen := map[string]bool{}
	for _, skill := range roadmap.CurrentSkills {
		for _, next := range skill.NextSkills {
			if !seen[next] && !nextSeen[next] {
				if nextSkill, ok := r.catalog[next]; ok {
					roadmap.NextSteps = append(roadmap.NextSteps, nextSkill)
					nextSeen[next] = true
				}
			}
		}
	}

	// Long-term goals: two hops out
	longTermSeen := map[string]bool{}
	for _, step := range roadmap.NextSteps {
		for _, next := range step.NextSkills {
			if !seen[next] && !nextSeen[next] && !longTermSeen[next] {
				if nextSkill, ok := r.catalog[next]; ok {
					roadmap.LongTermGoals = append(roadmap.LongTermGoals, nextSkill)
					longTermSeen[next] = true
				}
			}
		}
	}

	sortSkills(roadmap.NextSteps)
	sortSkills(roadmap.LongTermGoals)

	return roadmap
}

// ──────────────────────────────────────────────────────────────────────────────
// Deep analysis: confidence scoring
// ──────────────────────────────────────────────────────────────────────────────

// ComputeConfidence calculates a weighted confidence score (0–100) for a skill.
//
// Algorithm:
//
//	confidence = (fileScore × 0.25) + (commitScore × 0.30) + (recencyScore × 0.25) + (diversityScore × 0.20)
//
// Each component is normalized to 0–100 before weighting.
func ComputeConfidence(fileCount, commitCount, daysSinceLastCommit, repoCount int) float64 {
	// File score: logarithmic scale, caps at ~200 files
	fileScore := math.Min(100, (math.Log2(float64(fileCount)+1)/math.Log2(201))*100)

	// Commit score: logarithmic scale, caps at ~500 commits
	commitScore := math.Min(100, (math.Log2(float64(commitCount)+1)/math.Log2(501))*100)

	// Recency score: exponential decay — recent code scores higher
	// 0 days = 100, 30 days = ~74, 90 days = ~41, 365 days = ~5
	recencyScore := 100 * math.Exp(-0.01*float64(daysSinceLastCommit))

	// Diversity score: number of repos (capped at 10)
	diversityScore := math.Min(100, float64(repoCount)*10)

	confidence := (fileScore * 0.25) + (commitScore * 0.30) + (recencyScore * 0.25) + (diversityScore * 0.20)
	return math.Round(confidence*10) / 10 // round to 1 decimal
}

// LevelFromConfidence maps a confidence score to a skill level label.
func LevelFromConfidence(confidence float64) string {
	switch {
	case confidence >= 80:
		return "expert"
	case confidence >= 55:
		return "advanced"
	case confidence >= 30:
		return "intermediate"
	default:
		return "beginner"
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Print (terminal output)
// ──────────────────────────────────────────────────────────────────────────────

// Print renders the roadmap to stdout (standard mode).
func (r *Roadmap) Print() {
	fmt.Println("\n╔══════════════════════════════════════════════════")
	fmt.Printf("║  %s\n", r.Title)
	fmt.Println("╚══════════════════════════════════════════════════")

	if len(r.CurrentSkills) > 0 {
		fmt.Println("\n  [CURRENT SKILLS DETECTED]")
		for _, s := range r.CurrentSkills {
			fmt.Printf("   - %-20s [%s] %s\n", s.Name, s.Level, s.Category)
		}
	}

	if len(r.NextSteps) > 0 {
		fmt.Println("\n  [RECOMMENDED NEXT STEPS]")
		for _, s := range r.NextSteps {
			fmt.Printf("   - %-20s [%s]\n", s.Name, s.Level)
			if len(s.Resources) > 0 {
				fmt.Printf("     Link: %s\n", s.Resources[0])
			}
		}
	}

	if len(r.LongTermGoals) > 0 {
		fmt.Println("\n  [LONG-TERM GOALS]")
		for _, s := range r.LongTermGoals {
			fmt.Printf("   - %-20s [%s]\n", s.Name, s.Level)
		}
	}
	fmt.Println()
}

// PrintDeep renders the deep analysis with confidence scores.
func (r *Roadmap) PrintDeep() {
	fmt.Println("\n╔══════════════════════════════════════════════════")
	fmt.Printf("║  %s\n", r.Title)
	fmt.Printf("║  Deep Analysis · %s\n", r.GeneratedAt)
	fmt.Println("╚══════════════════════════════════════════════════")

	if len(r.DeepStats) > 0 {
		fmt.Println("\n  [SKILL MATRIX — Confidence Scoring]")
		fmt.Println("  ┌─────────────────────┬───────────┬────────────┬─────────┬─────────┐")
		fmt.Println("  │ Technology          │ Level     │ Confidence │ Files   │ Commits │")
		fmt.Println("  ├─────────────────────┼───────────┼────────────┼─────────┼─────────┤")
		for _, s := range r.DeepStats {
			fmt.Printf("  │ %-19s │ %-9s │ %5.1f%%     │ %-7d │ %-7d │\n",
				s.Name, s.Level, s.Confidence, s.FileCount, s.CommitCount)
		}
		fmt.Println("  └─────────────────────┴───────────┴────────────┴─────────┴─────────┘")
	}

	if len(r.CurrentSkills) > 0 {
		fmt.Println("\n  [CURRENT SKILLS]")
		for _, s := range r.CurrentSkills {
			fmt.Printf("   - %-20s %s\n", s.Name, s.Category)
		}
	}

	if len(r.NextSteps) > 0 {
		fmt.Println("\n  [RECOMMENDED UPGRADES]")
		for _, s := range r.NextSteps {
			fmt.Printf("   → %-20s [%s]\n", s.Name, s.Level)
			if len(s.Resources) > 0 {
				fmt.Printf("     Learn: %s\n", s.Resources[0])
			}
		}
	}

	if len(r.AIInsights) > 0 {
		fmt.Println("\n  [AI-POWERED INSIGHTS — Perplexity]")
		for _, insight := range r.AIInsights {
			fmt.Printf("   💡 %s\n", insight)
		}
	}

	fmt.Println()
}

// ──────────────────────────────────────────────────────────────────────────────
// Export (multi-format)
// ──────────────────────────────────────────────────────────────────────────────

// ExportJSON returns the roadmap as pretty-printed JSON.
func (r *Roadmap) ExportJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(data), nil
}

// ExportMarkdown returns the roadmap as a Markdown document.
func (r *Roadmap) ExportMarkdown() string {
	var b strings.Builder

	b.WriteString("# AutoDev Skills Report\n\n")
	b.WriteString(fmt.Sprintf("> Generated: %s\n\n", r.GeneratedAt))
	b.WriteString(fmt.Sprintf("**%s**\n\n", r.Title))

	// Deep stats table
	if len(r.DeepStats) > 0 {
		b.WriteString("## Skill Matrix\n\n")
		b.WriteString("| Technology | Level | Confidence | Files | Commits | Last Active |\n")
		b.WriteString("|------------|-------|------------|-------|---------|-------------|\n")
		for _, s := range r.DeepStats {
			b.WriteString(fmt.Sprintf("| %s | %s | %.1f%% | %d | %d | %d days ago |\n",
				s.Name, s.Level, s.Confidence, s.FileCount, s.CommitCount, s.DaysSince))
		}
		b.WriteString("\n")
	}

	// Current skills
	if len(r.CurrentSkills) > 0 {
		b.WriteString("## Current Skills\n\n")
		for _, s := range r.CurrentSkills {
			b.WriteString(fmt.Sprintf("- **%s** — %s (%s)\n", s.Name, s.Category, s.Level))
		}
		b.WriteString("\n")
	}

	// Next steps
	if len(r.NextSteps) > 0 {
		b.WriteString("## Recommended Next Steps\n\n")
		for _, s := range r.NextSteps {
			b.WriteString(fmt.Sprintf("- **%s** [%s]", s.Name, s.Level))
			if len(s.Resources) > 0 {
				b.WriteString(fmt.Sprintf(" — [Learn](%s)", s.Resources[0]))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Long-term goals
	if len(r.LongTermGoals) > 0 {
		b.WriteString("## Long-Term Goals\n\n")
		for _, s := range r.LongTermGoals {
			b.WriteString(fmt.Sprintf("- **%s** [%s]\n", s.Name, s.Level))
		}
		b.WriteString("\n")
	}

	// AI insights
	if len(r.AIInsights) > 0 {
		b.WriteString("## AI-Powered Insights\n\n")
		for _, insight := range r.AIInsights {
			b.WriteString(fmt.Sprintf("- 💡 %s\n", insight))
		}
		b.WriteString("\n")
	}

	b.WriteString("---\n\n")
	b.WriteString("*Generated by [AutoDev](https://github.com/HEETMEHTA18/autodev) · Powered by [skills.sh](https://skills.sh)*\n")

	return b.String()
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

func sortSkills(skills []Skill) {
	sort.Slice(skills, func(i, j int) bool {
		levelOrder := map[string]int{"beginner": 0, "intermediate": 1, "advanced": 2, "expert": 3}
		return levelOrder[skills[i].Level] < levelOrder[skills[j].Level]
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// Skill catalog (expanded from 22 → 30+)
// ──────────────────────────────────────────────────────────────────────────────

func buildCatalog() map[string]Skill {
	return map[string]Skill{
		"Node.js": {
			Name: "Node.js", Category: "Runtime", Level: "beginner",
			Resources:  []string{"https://nodejs.org/en/learn", "https://nodeschool.io"},
			NextSkills: []string{"Express", "TypeScript", "Docker"},
		},
		"TypeScript": {
			Name: "TypeScript", Category: "Language", Level: "intermediate",
			Resources:  []string{"https://www.typescriptlang.org/docs", "https://execute-program.com/courses/typescript"},
			NextSkills: []string{"Next.js", "NestJS", "tRPC"},
		},
		"React": {
			Name: "React", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://react.dev", "https://ui.dev/react"},
			NextSkills: []string{"Next.js", "React Query", "Zustand"},
		},
		"Next.js": {
			Name: "Next.js", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://nextjs.org/learn", "https://masteringnextjs.com"},
			NextSkills: []string{"Docker", "CI/CD", "Kubernetes", "Vercel"},
		},
		"Python": {
			Name: "Python", Category: "Language", Level: "beginner",
			Resources:  []string{"https://docs.python.org/3/tutorial", "https://realpython.com"},
			NextSkills: []string{"FastAPI", "Django", "Docker", "PyTorch"},
		},
		"FastAPI": {
			Name: "FastAPI", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://fastapi.tiangolo.com"},
			NextSkills: []string{"Celery", "Docker", "PostgreSQL"},
		},
		"Django": {
			Name: "Django", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://docs.djangoproject.com/en/stable/intro/tutorial01"},
			NextSkills: []string{"Django REST Framework", "Celery", "Redis"},
		},
		"Go": {
			Name: "Go", Category: "Language", Level: "intermediate",
			Resources:  []string{"https://go.dev/tour", "https://gobyexample.com"},
			NextSkills: []string{"Kubernetes", "gRPC", "Terraform", "Docker"},
		},
		"Rust": {
			Name: "Rust", Category: "Language", Level: "advanced",
			Resources:  []string{"https://doc.rust-lang.org/book", "https://rustlings.cool"},
			NextSkills: []string{"WebAssembly", "tokio", "axum"},
		},
		"Docker": {
			Name: "Docker", Category: "DevOps", Level: "intermediate",
			Resources:  []string{"https://docs.docker.com/get-started", "https://labs.play-with-docker.com"},
			NextSkills: []string{"Kubernetes", "Docker Compose", "GitHub Actions"},
		},
		"Kubernetes": {
			Name: "Kubernetes", Category: "DevOps", Level: "advanced",
			Resources:  []string{"https://kubernetes.io/docs/tutorials", "https://killercoda.com"},
			NextSkills: []string{"Helm", "Istio", "ArgoCD"},
		},
		"CI/CD": {
			Name: "CI/CD", Category: "DevOps", Level: "intermediate",
			Resources:  []string{"https://docs.github.com/actions", "https://docs.gitlab.com/ee/ci"},
			NextSkills: []string{"ArgoCD", "Terraform", "Kubernetes"},
		},
		"Terraform": {
			Name: "Terraform", Category: "Infrastructure", Level: "advanced",
			Resources:  []string{"https://developer.hashicorp.com/terraform/tutorials"},
			NextSkills: []string{"Pulumi", "Ansible", "Vault"},
		},
		"Flutter": {
			Name: "Flutter", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://flutter.dev/learn", "https://docs.flutter.dev/get-started"},
			NextSkills: []string{"Dart", "Firebase", "Android SDK"},
		},
		"Java": {
			Name: "Java", Category: "Language", Level: "intermediate",
			Resources:  []string{"https://dev.java/learn", "https://openjdk.org"},
			NextSkills: []string{"Spring Boot", "Maven", "Docker"},
		},
		"Spring Boot": {
			Name: "Spring Boot", Category: "Framework", Level: "advanced",
			Resources:  []string{"https://spring.io/guides"},
			NextSkills: []string{"Kubernetes", "Kafka", "gRPC"},
		},
		"PHP": {
			Name: "PHP", Category: "Language", Level: "beginner",
			Resources:  []string{"https://www.php.net/manual/en/getting-started.php"},
			NextSkills: []string{"Laravel", "Composer", "Docker"},
		},
		"Laravel": {
			Name: "Laravel", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://laravel.com/docs"},
			NextSkills: []string{"Docker", "Redis", "Horizon"},
		},
		"Ruby": {
			Name: "Ruby", Category: "Language", Level: "beginner",
			Resources:  []string{"https://www.ruby-lang.org/en/documentation/quickstart"},
			NextSkills: []string{"Ruby on Rails", "Sidekiq", "Docker"},
		},
		"Ruby on Rails": {
			Name: "Ruby on Rails", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://guides.rubyonrails.org"},
			NextSkills: []string{"Docker", "PostgreSQL", "Redis"},
		},
		"Firebase": {
			Name: "Firebase", Category: "Database", Level: "beginner",
			Resources:  []string{"https://firebase.google.com/docs"},
			NextSkills: []string{"Supabase", "GCP", "Cloud Functions"},
		},
		"Supabase": {
			Name: "Supabase", Category: "Database", Level: "intermediate",
			Resources:  []string{"https://supabase.com/docs"},
			NextSkills: []string{"PostgreSQL", "Edge Functions", "Realtime"},
		},
		"PostgreSQL": {
			Name: "PostgreSQL", Category: "Database", Level: "intermediate",
			Resources:  []string{"https://www.postgresql.org/docs/current/tutorial.html"},
			NextSkills: []string{"TimescaleDB", "pgvector", "Prisma"},
		},
		// ── New additions ─────────────────────────────────────────
		"Redis": {
			Name: "Redis", Category: "Database", Level: "intermediate",
			Resources:  []string{"https://redis.io/learn", "https://university.redis.com"},
			NextSkills: []string{"Kafka", "RabbitMQ", "Celery"},
		},
		"MongoDB": {
			Name: "MongoDB", Category: "Database", Level: "intermediate",
			Resources:  []string{"https://learn.mongodb.com", "https://university.mongodb.com"},
			NextSkills: []string{"Mongoose", "Atlas", "Aggregation"},
		},
		"PyTorch": {
			Name: "PyTorch", Category: "AI/ML", Level: "advanced",
			Resources:  []string{"https://pytorch.org/tutorials", "https://fast.ai"},
			NextSkills: []string{"Hugging Face", "CUDA", "ONNX"},
		},
		"TensorFlow": {
			Name: "TensorFlow", Category: "AI/ML", Level: "advanced",
			Resources:  []string{"https://www.tensorflow.org/tutorials"},
			NextSkills: []string{"Keras", "TFLite", "TF Serving"},
		},
		"Svelte": {
			Name: "Svelte", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://svelte.dev/tutorial", "https://learn.svelte.dev"},
			NextSkills: []string{"SvelteKit", "Docker", "Vercel"},
		},
		"Vue": {
			Name: "Vue", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://vuejs.org/tutorial", "https://vueschool.io"},
			NextSkills: []string{"Nuxt", "Pinia", "Docker"},
		},
		"Angular": {
			Name: "Angular", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://angular.dev/tutorials", "https://angular.io/start"},
			NextSkills: []string{"RxJS", "NgRx", "Docker"},
		},
		"Express": {
			Name: "Express", Category: "Framework", Level: "beginner",
			Resources:  []string{"https://expressjs.com/en/starter/installing.html"},
			NextSkills: []string{"NestJS", "Docker", "PostgreSQL"},
		},
		"NestJS": {
			Name: "NestJS", Category: "Framework", Level: "intermediate",
			Resources:  []string{"https://docs.nestjs.com"},
			NextSkills: []string{"Kubernetes", "GraphQL", "TypeORM"},
		},
	}
}
