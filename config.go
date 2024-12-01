package main

import "github.com/joeshaw/envdecode"

type Config struct {
	Port             string `env:"PORT" envDefault:"8080"`
	ConnectionString string `env:"CONNECTION_STRING" envDefault:""`
}

func NewConfigFromEnv() (Config, error) {
	cfg := Config{}
	if err := envdecode.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
