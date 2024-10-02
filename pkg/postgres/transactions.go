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
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUnableGetTransactionFromContext = errors.New("unable get transaction from context")
	ErrNotInContextualTxStatement      = errors.New("unable to commit transaction statement - not in tx statement")
)

type transactionCtxKey string

var transactionKey = transactionCtxKey("transaction")
var transactionCommittedKey = transactionCtxKey("is_committed")

// BeginTx ....
func (c *Connection) BeginTx() (*sqlx.Tx, error) {
	return c.Dbx.Beginx()
}

// BeginTxWithRollbackOnError ....
func (c *Connection) BeginTxWithRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	err := c.BeginReadCommittedTxRollbackOnError(ctx, callback)
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

func (c *Connection) BeginReadCommittedTxRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	txStmt, err := c.Dbx.Beginx()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	newCtx := context.WithValue(ctx, transactionKey, txStmt)
	err = callback(newCtx)
	if err != nil {
		rollbackErr := txStmt.Rollback()
		if rollbackErr != nil {
			return c.e.ErrorOnly(rollbackErr)
		}

		return c.e.ErrorOnly(err)
	}

	err = txStmt.Commit()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

func (c *Connection) BeginReadUncommittedTxRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	txStmt, err := c.Dbx.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	newCtx := context.WithValue(ctx, transactionKey, txStmt)
	err = callback(newCtx)
	if err != nil {
		rollbackErr := txStmt.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return c.e.ErrorOnly(err)
	}

	err = txStmt.Commit()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	return nil
}

// BeginContextualTxStatement ....
func (c *Connection) BeginContextualTxStatement(ctx context.Context) (context.Context, error) {
	txStmt, err := c.Dbx.Beginx()
	if err != nil {
		return nil, c.e.ErrorOnly(err)
	}

	return context.WithValue(ctx, transactionKey, txStmt), nil
}

// CommitContextualTxStatement ....
func (c *Connection) CommitContextualTxStatement(ctx context.Context) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if !inTransaction {
		return c.e.ErrorOnly(ErrNotInContextualTxStatement)
	}

	return tx.Commit()
}

// RollbackContextualTxStatement ....
func (c *Connection) RollbackContextualTxStatement(ctx context.Context) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if !inTransaction {
		return c.e.ErrorOnly(ErrNotInContextualTxStatement)
	}

	return tx.Rollback()
}

func (c *Connection) TryWithTransaction(ctx context.Context, fn func(stmt sqlx.Ext) error) error {
	stmt := sqlx.Ext(c.Dbx)

	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if inTransaction {
		stmt = tx
	}

	return fn(stmt)
}

func (c *Connection) MustWithTransaction(ctx context.Context, fn func(stmt *sqlx.Tx) error) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if inTransaction {
		return fn(tx)
	}

	return c.e.ErrorOnly(ErrUnableGetTransactionFromContext)
}
