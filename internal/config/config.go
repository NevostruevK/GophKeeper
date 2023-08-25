package config

import "github.com/NevostruevK/GophKeeper/internal/config/duration"

type Config struct {
	TokenDuration duration.Duration `json:"token_duration"`
	Address       string            `json:"address"`
	TokenKey      string            `json:"token_key"`
	DSN           string            `json:"dsn"`
	UseTLS        bool              `json:"use_tls"`
}

func NewConfig(options ...func(*Config)) *Config {
	cfg := &Config{}
	for _, o := range options {
		o(cfg)
	}
	return cfg
}

func (o *Config) SetOption(set func(*Config)) {
	set(o)
}

func WithTokenDuration(v duration.Duration) func(*Config) {
	return func(o *Config) {
		o.TokenDuration = v
	}
}

func WithAddress(v string) func(*Config) {
	return func(o *Config) {
		o.Address = v
	}
}

func WithTokenKey(v string) func(*Config) {
	return func(o *Config) {
		o.TokenKey = v
	}
}

func WithDSN(v string) func(*Config) {
	return func(o *Config) {
		o.DSN = v
	}
}

func WithUseTLS(v bool) func(*Config) {
	return func(o *Config) {
		o.UseTLS = v
	}
}
