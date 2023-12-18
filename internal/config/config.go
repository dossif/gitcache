package config

import (
	"encoding/json"
	"fmt"
	"github.com/dossif/secretvar"
	"github.com/kelseyhightower/envconfig"
	"os"
	"time"
)

type Config struct {
	Listen    string `default:"127.0.0.1:8080"`
	LogLevel  string `default:"info"`
	LogGelf   string `default:""`
	GitSource GitSource
	GitMirror GitMirror
}

type GitSource struct {
	Url           string        `required:"true"`
	Path          string        `default:"/tmp/"`
	CheckInterval time.Duration `default:"5m"`
	Auth          bool          `default:"False"`
	BasicAuth     BasicAuth
}

type GitMirror struct {
	Path      string `default:""`
	Auth      bool   `default:"False"`
	BasicAuth BasicAuth
}

type BasicAuth struct {
	Username secretvar.SecretString `default:""`
	Password secretvar.SecretString `default:""`
}

func NewConfig(prefix string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(prefix, &cfg)
	if err != nil {
		return &cfg, err
	}
	err = envconfig.CheckDisallowed(prefix, &cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, nil
}

func PrintUsage(prefix string) {
	var cfg Config
	_ = envconfig.Usage(prefix, &cfg)
	os.Exit(128)
}

func DebugConfig(cfg *Config) {
	fmt.Println("debug config:")
	j, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(j))
}
