package postgres

import (
	"fmt"
)

type PostgresConfig struct {
	DbHost         string `envconfig:"POSTGRESQL_SERVICE_HOST"`
	DbPort         uint16 `envconfig:"POSTGRESQL_SERVICE_PORT"`
	DbName         string `envconfig:"POSTGRESQL_DATABASE_NAME" secret:"true"`
	DbUsername     string `envconfig:"POSTGRESQL_USERNAME" secret:"true"`
	DbPassword     string `envconfig:"POSTGRESQL_PASSWORD" secret:"true"`
	DbSSLMode      string `envconfig:"POSTGRESQL_SSL_MODE" default:"prefer"`
	DbMaxOpenConns uint8  `envconfig:"POSTGRESQL_MAX_OPEN_CONNECTIONS" default:"8"`
	DbMaxIdleConns uint8  `envconfig:"POSTGRESQL_MAX_IDLE_CONNECTIONS" default:"8" `
	// DbConnectRetryCount is the maximum number of reconnection tries. If 0 - infinite loop
	DbConnectRetryCount uint8 `envconfig:"POSTGRESQL_CONNECTION_RETRY_COUNT" default:"0"`
	// DbConnectTimeOut is the timeout in millisecond to connect between connection tries
	DbConnectTimeOut uint16 `envconfig:"POSTGRESQL_CONNECTION_RETRY_TIMEOUT" default:"5000"`
}

func (c *PostgresConfig) Prepare() error {
	return nil
}

func (c *PostgresConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSSLMode)
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
