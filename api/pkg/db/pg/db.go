package pg

import (
	"context"
	"database/sql"
)

// friendsDB wraps the *sql.DB provided
type friendsDB struct {
	*sql.DB
	info InstanceInfo
}

// BeginTx begins a transaction with the database in receiver and returns a Transactor
func (i friendsDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transactor, error) {
	return i.DB.BeginTx(ctx, opts)
}

// ExecContext wraps the base connector
func (i friendsDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.DB.ExecContext(ctx, query, args...)
}

// QueryContext wraps the base connector
func (i friendsDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return i.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext wraps the base connector
func (i friendsDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return i.DB.QueryRowContext(ctx, query, args...)
}

// InstanceInfo returns info about the DB
func (i friendsDB) InstanceInfo() InstanceInfo {
	return i.info
}

// friendsTx wraps the Transactor provided
type friendsTx struct {
	Transactor
	info InstanceInfo
	ctx  context.Context
}

// ExecContext wraps the base connector
func (i friendsTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.Transactor.ExecContext(ctx, query, args...)
}

// QueryContext wraps the base connector
func (i friendsTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return i.Transactor.QueryContext(ctx, query, args...)
}

// QueryRowContext wraps the base connector
func (i friendsTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return i.Transactor.QueryRowContext(ctx, query, args...)
}

// Commit commits the transaction
func (i friendsTx) Commit() error {
	return i.Transactor.Commit()
}

// Rollback aborts the transaction
func (i friendsTx) Rollback() error {
	return i.Transactor.Rollback()
}
