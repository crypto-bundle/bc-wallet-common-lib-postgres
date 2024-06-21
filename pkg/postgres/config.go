/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DbHost, c.DbPort, c.DbUsername, c.DbPassword, c.DbName, c.DbSSLMode)
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
