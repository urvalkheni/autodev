// Package config manages AutoDev configuration via Viper.
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Version    string `mapstructure:"version"`
	LogLevel   string `mapstructure:"log_level"`
	OutputDir  string `mapstructure:"output_dir"`
	GitHubToken string `mapstructure:"github_token"`
	SkillsAPIURL string `mapstructure:"skills_api_url"`
	NoColor    bool   `mapstructure:"no_color"`
	NoTelemetry bool  `mapstructure:"no_telemetry"`
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version:      "0.1.0",
		LogLevel:     "info",
		OutputDir:    "./autodev-reports",
		SkillsAPIURL: "https://www.skills.sh/api",
		NoColor:      false,
		NoTelemetry:  true,
	}
}

// Load reads config from file and environment, returning a merged Config.
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	def := DefaultConfig()
	v.SetDefault("version", def.Version)
	v.SetDefault("log_level", def.LogLevel)
	v.SetDefault("output_dir", def.OutputDir)
	v.SetDefault("skills_api_url", def.SkillsAPIURL)
	v.SetDefault("no_color", def.NoColor)
	v.SetDefault("no_telemetry", def.NoTelemetry)

	// Config file
	v.SetConfigName(".autodev")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath(configDir())
	_ = v.ReadInConfig() // ignore if not found

	// Environment variables
	v.SetEnvPrefix("AUTODEV")
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// configDir returns the platform config directory.
func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".config", "autodev")
}
