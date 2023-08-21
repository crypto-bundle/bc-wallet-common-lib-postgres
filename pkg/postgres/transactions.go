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

type txStmtCallback func(txStmtCtx context.Context) error

// BeginTx ....
func (c *Connection) BeginTx() (*sqlx.Tx, error) {
	return c.Dbx.Beginx()
}

// BeginTxWithRollbackOnError ....
func (c *Connection) BeginTxWithRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	return c.BeginReadCommittedTxRollbackOnError(ctx, callback)
}

func (c *Connection) BeginReadCommittedTxRollbackOnError(ctx context.Context,
	callback func(txStmtCtx context.Context) error,
) error {
	txStmt, err := c.Dbx.Beginx()
	if err != nil {
		return err
	}

	newCtx := context.WithValue(ctx, transactionKey, txStmt)
	err = callback(newCtx)
	if err != nil {
		rollbackErr := txStmt.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	err = txStmt.Commit()
	if err != nil {
		return err
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
		return err
	}

	newCtx := context.WithValue(ctx, transactionKey, txStmt)
	err = callback(newCtx)
	if err != nil {
		rollbackErr := txStmt.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	err = txStmt.Commit()
	if err != nil {
		return err
	}

	return nil
}

// BeginContextualTxStatement ....
func (c *Connection) BeginContextualTxStatement(ctx context.Context) (context.Context, error) {
	txStmt, err := c.Dbx.Beginx()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, transactionKey, txStmt), nil
}

// CommitContextualTxStatement ....
func (c *Connection) CommitContextualTxStatement(ctx context.Context) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if !inTransaction {
		return ErrNotInContextualTxStatement
	}

	return tx.Commit()
}

// RollbackContextualTxStatement ....
func (c *Connection) RollbackContextualTxStatement(ctx context.Context) error {
	tx, inTransaction := ctx.Value(transactionKey).(*sqlx.Tx)
	if !inTransaction {
		return ErrNotInContextualTxStatement
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

	return ErrUnableGetTransactionFromContext
}
