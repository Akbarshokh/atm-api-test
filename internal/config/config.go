package config

import (
	"fmt"
)

type Config struct {
	// Server
	LogLevel   string `env:"LOG_LEVEL, default=debug"`
	ServerIP   string `env:"SERVER_IP"`
	HTTPPort   string `env:"HTTP_PORT"`
	ServerHost string `env:"SERVER_HOST"`

	// Db connection
	Postgres *DB `env:",prefix=POSTGRES_"`
}

type DB struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
	SSLMode  string `env:"SSLMODE"`
}

func (c *DB) PostgresURL() string {
	if c.Port == 0 {
		c.Port = 5432
	}

	if c.Host == "" {
		c.Host = "192.168.100.1"
	}

	if c.Username == "" {
		c.Username = "prod_twowish"
	}

	if c.Password == "" {
		c.Password = "kVUj76KfUnjD1TYALUNF"
	}

	if c.Database == "" {
		c.Database = "prod_twowish_01"
	}

	if c.Username == "" {
		return fmt.Sprintf(
			"host=%s port=%d dbname=%s sslmode=disable",
			c.Host,
			c.Port,
			c.Database,
			//c.Timezone,
		)
	}

	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host,
			c.Port,
			c.Username,
			c.Database,
		)
	}

	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Database,
	)
}
