package postgres

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const ConfigPrefix = "DB"

type PostgresConfig struct {
	DbHost         string `envconfig:"DB_HOST" json:"-"`
	DbPort         uint16 `envconfig:"DB_PORT" json:"-"`
	DbName         string `envconfig:"DB_DATABASE" json:"-"`
	DbUsername     string `envconfig:"DB_USERNAME" json:"DB_USERNAME"`
	DbPassword     string `envconfig:"DB_PASSWORD" json:"DB_PASSWORD"`
	DbSSLMode      string `envconfig:"DB_SSL_MODE" default:"prefer" json:"-"`
	DbMaxOpenConns uint8  `envconfig:"DB_MAX_OPEN_CONNECTIONS" default:"8" json:"-"`
	DbMaxIdleConns uint8  `envconfig:"DB_MAX_IDLE_CONNECTIONS" default:"8" json:"-"`
	// DbConnectRetryCount is the maximum number of reconnection tries. If 0 - infinite loop
	DbConnectRetryCount uint8 `envconfig:"DB_RETRY_COUNT" default:"0" json:"-"`
	// DbConnectTimeOut is the timeout in millisecond to connect between connection tries
	DbConnectTimeOut uint16 `envconfig:"DB_RETRY_TIMEOUT" default:"5000" json:"-"`
}

func (c *PostgresConfig) Prepare() error {
	err := envconfig.Process(ConfigPrefix, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *PostgresConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("mysql://%s:%s@%s/%s?sslmode=%t",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbName, c.DbSSLMode)
}

func (c *PostgresConfig) GetDbHost() string {
	return c.DbHost
}

func (c *PostgresConfig) GetDbPort() uint16 {
	return c.DbPort
}

func (c *PostgresConfig) GetDbName() string {
	return c.DbName
}

func (c *PostgresConfig) GetDbUser() string {
	return c.DbUsername
}

func (c *PostgresConfig) GetDbPassword() string {
	return c.DbPassword
}

func (c *PostgresConfig) GetDbTLSMode() string {
	return c.DbSSLMode
}

func (c *PostgresConfig) GetDbRetryCount() uint8 {
	return c.DbConnectRetryCount
}

func (c *PostgresConfig) GetDbConnectTimeOut() uint16 {
	return c.DbConnectTimeOut
}

func (c *PostgresConfig) GetDbMaxOpenConns() uint8 {
	return c.DbMaxOpenConns
}

func (c *PostgresConfig) GetDbMaxIdleConns() uint8 {
	return c.DbMaxIdleConns
}
