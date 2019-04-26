package main

import (
	"io"

	"github.com/mattn/go-mastodon"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Secrets ConfigSecrets `yaml:"secrets"`
	List    mastodon.ID   `yaml:"list"`
}

type ConfigSecrets struct {
	Server       string `yaml:"server"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	AccessToken  string `yaml:"access_token"`
}

func readConfig(r io.Reader) (*Config, error) {
	var cfg Config
	if err := yaml.NewDecoder(r).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func writeConfig(w io.Writer, cfg *Config) error {
	enc := yaml.NewEncoder(w)
	if err := enc.Encode(cfg); err != nil {
		return err
	}
	return enc.Close()
}

func createClient(cfg *Config) *mastodon.Client {
	return mastodon.NewClient((*mastodon.Config)(&cfg.Secrets))
}
