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
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type connectionParams struct {
	host     string
	port     uint16
	user     string
	password string
	database string

	retryTimeOut time.Duration
	retryCount   uint8

	readTimeOut  int
	writeTimeOut int

	maxOpenConn uint8
	maxIdleConn uint8

	sslMode string

	debug bool
}

// Connection struct to store and manipulate postgres database connection
type Connection struct {
	Dbx *sqlx.DB

	params *connectionParams

	l *log.Logger
}

func (c *Connection) IsHealed(ctx context.Context) bool {
	err := c.Dbx.PingContext(ctx)
	if err != nil {
		return false
	}

	return true
}

func (c *Connection) Close() error {
	err := c.Dbx.Close()
	if err != nil {
		return err
	}

	return nil
}

// Connect to postgres database
func (c *Connection) Connect() (*Connection, error) {
	retryDecValue := uint8(1)
	retryCount := c.params.retryCount
	if retryCount == 0 {
		retryDecValue = 0
		retryCount = 1
	}
	try := 0

	for i := retryCount; i != 0; i -= retryDecValue {
		dbx, err := sqlx.Connect("postgres", formatPostgresDSN(c.params))
		if err != nil {
			c.l.Printf("unable connect to postgres. reconnecting. error: %s. iteration: %d", err, try)
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}

		err = dbx.Ping()
		if err != nil {
			c.l.Printf("unable ping postgres. reconnecting. error: %s. iteration: %d", err, try)
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}

		rows, err := dbx.Query("SELECT 1")
		if err != nil {
			c.l.Printf("unable make sql request. reconnecting. error: %s. iteration: %d", err, try)
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}
		err = rows.Close()
		if err != nil {
			c.l.Printf("unable to close rows statement. reconnecting. error: %s. iteration: %d", err, try)
			try++
			time.Sleep(c.params.retryTimeOut)

			continue
		}

		dbx.SetMaxOpenConns(int(c.params.maxOpenConn))
		dbx.SetMaxIdleConns(int(c.params.maxIdleConn))

		c.Dbx = dbx
		return c, nil
	}

	return c, nil
}

// NewConnection to postgres db
func NewConnection(_ context.Context, cfg DbConfig, logger *log.Logger) *Connection {
	conn := &Connection{
		params: &connectionParams{
			host:     cfg.GetDbHost(),
			port:     cfg.GetDbPort(),
			user:     cfg.GetDbUser(),
			password: cfg.GetDbPassword(),
			database: cfg.GetDbName(),

			retryCount:   cfg.GetDbRetryCount(),
			retryTimeOut: time.Duration(cfg.GetDbConnectTimeOut()) * time.Millisecond,

			maxOpenConn: cfg.GetDbMaxOpenConns(),
			maxIdleConn: cfg.GetDbMaxIdleConns(),

			debug: cfg.IsDebug(),

			sslMode: cfg.GetDbTLSMode(),
		},
		l: logger,
	}

	return conn
}

func formatPostgresDSN(params *connectionParams) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		params.host, params.port, params.user, params.password, params.database, params.sslMode)
}
