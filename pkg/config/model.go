package config

import (
	"time"
)

type Config struct {
	SMTP    SMTPConfig    `yaml:"smtp"`
	Email   EmailConfig   `yaml:"email"`
	Storage StorageConfig `yaml:"storage"`
	HTTP    HTTPConfig    `yaml:"http"`
	KunaAPI KunaAPIConfig `yaml:"kuna_api"`
}

type EmailConfig struct {
	From    string `yaml:"from"`
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type StorageConfig struct {
	Path string `yaml:"path"`
}

type HTTPConfig struct {
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout_in_seconds"`
}

type KunaAPIConfig struct {
	URL        string  `yaml:"url"`
	DefaltRate float32 `yaml:"default_rate"`
}
