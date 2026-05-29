// Package github provides a client for the GitHub REST API.
// It supports fetching public repositories, detecting languages,
// and building environment setup plans for any GitHub user.
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const baseURL = "https://api.github.com"

// Client is a GitHub API client.
type Client struct {
	httpClient *http.Client
	token      string // optional PAT for higher rate limits
}

// Repo is a minimal GitHub repository representation.
type Repo struct {
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	Language        string    `json:"language"`
	Topics          []string  `json:"topics"`
	StargazersCount int       `json:"stargazers_count"`
	Archived        bool      `json:"archived"`
	Fork            bool      `json:"fork"`
	UpdatedAt       time.Time `json:"updated_at"`
	HTMLURL         string    `json:"html_url"`
}

// UserScanResult is the result of scanning all public repos for a user.
type UserScanResult struct {
	Username    string         `json:"username"`
	TotalRepos  int            `json:"total_repos"`
	Languages   map[string]int `json:"languages"` // lang → repo count
	Topics      map[string]int `json:"topics"`
	TopRepos    []Repo         `json:"top_repos"`
	Recommended []string       `json:"recommended_environment"`
	SkillGaps   []string       `json:"skill_gaps"`
}

// New creates a new GitHub API client.
// Pass an empty token to use unauthenticated requests (60 req/hr limit).
func New(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ScanUser fetches all public repos for a GitHub username and produces a UserScanResult.
func (c *Client) ScanUser(username string) (*UserScanResult, error) {
	repos, err := c.fetchAllRepos(username)
	if err != nil {
		return nil, err
	}

	result := &UserScanResult{
		Username:   username,
		TotalRepos: len(repos),
		Languages:  map[string]int{},
		Topics:     map[string]int{},
	}

	// Aggregate languages and topics
	for _, repo := range repos {
		if repo.Fork {
			continue // skip forks
		}
		if repo.Language != "" {
			result.Languages[repo.Language]++
		}
		for _, topic := range repo.Topics {
			result.Topics[topic]++
		}
	}

	// Top repos by stars
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].StargazersCount > repos[j].StargazersCount
	})
	if len(repos) > 10 {
		result.TopRepos = repos[:10]
	} else {
		result.TopRepos = repos
	}

	result.Recommended = buildRecommendedEnv(result.Languages)
	result.SkillGaps = buildSkillGaps(result.Languages)

	return result, nil
}

// fetchAllRepos paginates through the GitHub API to get all public repos.
func (c *Client) fetchAllRepos(username string) ([]Repo, error) {
	var all []Repo
	page := 1

	for {
		url := fmt.Sprintf("%s/users/%s/repos?per_page=100&page=%d&type=public", baseURL, username, page)
		var repos []Repo
		if err := c.get(url, &repos); err != nil {
			return nil, err
		}
		if len(repos) == 0 {
			break
		}
		all = append(all, repos...)
		if len(repos) < 100 {
			break
		}
		page++
	}
	return all, nil
}

// get makes an authenticated GET request and decodes JSON into v.
func (c *Client) get(url string, v any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "autodev-cli/0.1.0")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("github API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("github user not found")
	}
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("github rate limit exceeded — set AUTODEV_GITHUB_TOKEN to increase limits")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github API error %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// buildRecommendedEnv maps detected languages to recommended tools.
func buildRecommendedEnv(langs map[string]int) []string {
	envMap := map[string]string{
		"JavaScript": "Node.js 22",
		"TypeScript": "Node.js 22 + TypeScript",
		"Python":     "Python 3.12 + pip",
		"Go":         "Go 1.22",
		"Rust":       "Rust (rustup)",
		"Java":       "OpenJDK 21 + Maven/Gradle",
		"Kotlin":     "Kotlin + OpenJDK 21",
		"PHP":        "PHP 8.3 + Composer",
		"Ruby":       "Ruby 3.3 + Bundler",
		"Dart":       "Flutter SDK + Dart",
		"Swift":      "Xcode + Swift Toolchain",
		"C#":         ".NET 8 SDK",
		"C++":        "GCC / Clang + CMake",
		"Shell":      "Bash / Zsh",
		"Dockerfile": "Docker Desktop",
	}

	seen := map[string]bool{}
	var result []string
	for lang := range langs {
		if env, ok := envMap[lang]; ok && !seen[env] {
			result = append(result, env)
			seen[env] = true
		}
	}
	sort.Strings(result)
	return result
}

// buildSkillGaps suggests next skills based on existing language portfolio.
func buildSkillGaps(langs map[string]int) []string {
	gaps := map[string][]string{
		"JavaScript": {"TypeScript", "Docker", "CI/CD"},
		"TypeScript": {"Docker", "Kubernetes", "Go"},
		"Python":     {"Docker", "FastAPI", "Celery"},
		"Go":         {"Kubernetes", "gRPC", "Terraform"},
		"Java":       {"Spring Boot", "Docker", "Kubernetes"},
		"Rust":       {"WebAssembly", "async-std", "tokio"},
	}

	seen := map[string]bool{}
	for lang := range langs {
		seen[lang] = true
	}

	var result []string
	for lang := range langs {
		if suggestions, ok := gaps[lang]; ok {
			for _, s := range suggestions {
				if !seen[s] {
					result = append(result, s)
					seen[s] = true
				}
			}
		}
	}
	sort.Strings(result)
	return result
}
