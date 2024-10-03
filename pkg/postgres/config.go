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

var _ CommonDBConfig = (*PostgresConfig)(nil)

type PostgresConfig struct {
	DBHost         string `envconfig:"POSTGRESQL_SERVICE_HOST"`
	DBPort         uint16 `envconfig:"POSTGRESQL_SERVICE_PORT"`
	DBName         string `envconfig:"POSTGRESQL_DATABASE_NAME" secret:"true"`
	DBUsername     string `envconfig:"POSTGRESQL_USERNAME" secret:"true"`
	DBPassword     string `envconfig:"POSTGRESQL_PASSWORD" secret:"true"`
	DBSSLMode      string `envconfig:"POSTGRESQL_SSL_MODE" default:"prefer"`
	DBMaxOpenConns uint8  `envconfig:"POSTGRESQL_MAX_OPEN_CONNECTIONS" default:"8"`
	DBMaxIdleConns uint8  `envconfig:"POSTGRESQL_MAX_IDLE_CONNECTIONS" default:"8" `
	// DBConnectRetryCount is the maximum number of reconnection tries. If 0 - infinite loop
	DBConnectRetryCount uint8 `envconfig:"POSTGRESQL_CONNECTION_RETRY_COUNT" default:"0"`
	// DBConnectTimeOut is the timeout in millisecond to connect between connection tries
	DBConnectTimeOut uint16 `envconfig:"POSTGRESQL_CONNECTION_RETRY_TIMEOUT" default:"5000"`
}

func (c *PostgresConfig) Prepare() error {
	return nil
}

func (c *PostgresConfig) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, c.DBName, c.DBSSLMode)
}

func (c *PostgresConfig) GetDBHost() string {
	return c.DBHost
}

func (c *PostgresConfig) GetDBPort() uint16 {
	return c.DBPort
}

func (c *PostgresConfig) GetDBName() string {
	return c.DBName
}

func (c *PostgresConfig) GetDBUser() string {
	return c.DBUsername
}

func (c *PostgresConfig) GetDBPassword() string {
	return c.DBPassword
}

func (c *PostgresConfig) GetDBTLSMode() string {
	return c.DBSSLMode
}

func (c *PostgresConfig) GetDBRetryCount() uint8 {
	return c.DBConnectRetryCount
}

func (c *PostgresConfig) GetDBConnectTimeOut() uint16 {
	return c.DBConnectTimeOut
}

func (c *PostgresConfig) GetDBMaxOpenConns() uint8 {
	return c.DBMaxOpenConns
}

func (c *PostgresConfig) GetDBMaxIdleConns() uint8 {
	return c.DBMaxIdleConns
}
