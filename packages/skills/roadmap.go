// Package skills provides learning roadmap generation and skills.sh integration.
package skills

import (
	"fmt"
	"sort"
	"strings"
)

// Skill represents a technology skill in a roadmap.
type Skill struct {
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	Level      string   `json:"level"` // beginner | intermediate | advanced
	Resources  []string `json:"resources"`
	NextSkills []string `json:"next_skills"`
}

// Roadmap represents a personalized learning roadmap.
type Roadmap struct {
	Title         string  `json:"title"`
	CurrentSkills []Skill `json:"current_skills"`
	NextSteps     []Skill `json:"next_steps"`
	LongTermGoals []Skill `json:"long_term_goals"`
}

// RoadmapGenerator generates learning roadmaps from detected technologies.
type RoadmapGenerator struct {
	catalog map[string]Skill
}

// New creates a new RoadmapGenerator with the built-in skill catalog.
func New() *RoadmapGenerator {
	return &RoadmapGenerator{catalog: buildCatalog()}
}

// Generate produces a personalized roadmap from a list of detected technology names.
func (r *RoadmapGenerator) Generate(detected []string) *Roadmap {
	roadmap := &Roadmap{
		Title: fmt.Sprintf("Personalized Roadmap — %s", strings.Join(detected, ", ")),
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

// PrintRoadmap prints the roadmap to stdout.
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

func sortSkills(skills []Skill) {
	sort.Slice(skills, func(i, j int) bool {
		levelOrder := map[string]int{"beginner": 0, "intermediate": 1, "advanced": 2}
		return levelOrder[skills[i].Level] < levelOrder[skills[j].Level]
	})
}

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
			NextSkills: []string{"Docker", "CI/CD", "Kubernetes"},
		},
		"Python": {
			Name: "Python", Category: "Language", Level: "beginner",
			Resources:  []string{"https://docs.python.org/3/tutorial", "https://realpython.com"},
			NextSkills: []string{"FastAPI", "Django", "Docker"},
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
			NextSkills: []string{"Kubernetes", "gRPC", "Terraform"},
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
	}
}
