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
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func formatPostgresDSN(params *connectionParams) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		params.host, params.port, params.user, params.password, params.database, params.sslMode)
}

type connectionParams struct {
	host     string
	user     string
	password string
	sslMode  string

	database string

	retryTimeOut time.Duration
	port         uint16
	retryCount   uint8
	maxOpenConn  uint8
	maxIdleConn  uint8

	debug bool
}

// Connection struct to store and manipulate postgres database connection...
type Connection struct {
	l *slog.Logger
	e errorFormatterService

	Dbx *sqlx.DB

	params *connectionParams
}

func (c *Connection) IsHealed(ctx context.Context) bool {
	err := c.Dbx.PingContext(ctx)
	if err != nil {
		return false
	}

	err = c.checkConnectionByQuery(c.Dbx)

	return err == nil
}

func (c *Connection) checkConnectionByQuery(dbx *sqlx.DB) error {
	rows, err := dbx.Query("SELECT 1")
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	defer func() {
		_ = rows.Close()
	}()

	err = rows.Err()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

func (c *Connection) Close() error {
	err := c.Dbx.Close()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

// Connect to postgres database...
func (c *Connection) Connect() (*Connection, error) {
	retryDecValue := uint8(1)
	retryCount := c.params.retryCount

	if retryCount == 0 {
		retryDecValue = 0
		retryCount = 1
	}

	try := 0

	var err error

	for i := retryCount; i != 0; i -= retryDecValue {
		dbx, loopErr := c.tryConnect()
		if loopErr != nil {
			c.l.Error("unable to connect to database", loopErr,
				slog.Int(ConnectionRetryCountTag, try))

			err = loopErr

			time.Sleep(c.params.retryTimeOut)

			continue
		}

		dbx.SetMaxOpenConns(int(c.params.maxOpenConn))
		dbx.SetMaxIdleConns(int(c.params.maxIdleConn))

		c.Dbx = dbx

		return c, nil
	}

	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	return c, nil
}

func (c *Connection) tryConnect() (*sqlx.DB, error) {
	dbx, err := sqlx.Connect("postgres", formatPostgresDSN(c.params))
	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	err = dbx.Ping()
	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	err = c.checkConnectionByQuery(dbx)
	if err != nil {
		return nil, err
	}

	return dbx, nil
}

// NewConnection to postgres db...
func NewConnection(_ context.Context,
	logFactorySvc loggerService,
	errFormatterSvc errorFormatterService,
	cfgSvc DBConfigService,
) *Connection {
	conn := &Connection{
		e: errFormatterSvc,
		l: logFactorySvc.NewSlogNamedLoggerEntry("lib-postgres"),
		params: &connectionParams{
			host:     cfgSvc.GetDBHost(),
			port:     cfgSvc.GetDBPort(),
			user:     cfgSvc.GetDBUser(),
			password: cfgSvc.GetDBPassword(),
			database: cfgSvc.GetDBName(),

			retryCount:   cfgSvc.GetDBRetryCount(),
			retryTimeOut: time.Duration(cfgSvc.GetDBConnectTimeOut()) * time.Millisecond,

			maxOpenConn: cfgSvc.GetDBMaxOpenConns(),
			maxIdleConn: cfgSvc.GetDBMaxIdleConns(),

			debug: cfgSvc.IsDebug(),

			sslMode: cfgSvc.GetDBTLSMode(),
		},
		Dbx: nil,
	}

	return conn
}
