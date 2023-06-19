package config

type Config struct {
	SMTP  SMTPConfig  `yaml:"smtp"`
	Email EmailConfig `yaml:"email"`
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
