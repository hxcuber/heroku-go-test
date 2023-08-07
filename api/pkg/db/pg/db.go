package pg

import (
	"context"
	"database/sql"
)

// hollywoodDB wraps the *sql.DB provided
type hollywoodDB struct {
	*sql.DB
	info InstanceInfo
}

// BeginTx begins a transaction with the database in receiver and returns a Transactor
func (i hollywoodDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transactor, error) {
	return i.DB.BeginTx(ctx, opts)
}

// ExecContext wraps the base connector
func (i hollywoodDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.DB.ExecContext(ctx, query, args...)
}

// QueryContext wraps the base connector
func (i hollywoodDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return i.DB.QueryContext(ctx, query, args...)
}

// QueryRowContext wraps the base connector
func (i hollywoodDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return i.DB.QueryRowContext(ctx, query, args...)
}

// InstanceInfo returns info about the DB
func (i hollywoodDB) InstanceInfo() InstanceInfo {
	return i.info
}

// hollywoodTx wraps the Transactor provided
type hollywoodTx struct {
	Transactor
	info InstanceInfo
	ctx  context.Context
}

// ExecContext wraps the base connector
func (i hollywoodTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return i.Transactor.ExecContext(ctx, query, args...)
}

// QueryContext wraps the base connector
func (i hollywoodTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return i.Transactor.QueryContext(ctx, query, args...)
}

// QueryRowContext wraps the base connector
func (i hollywoodTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return i.Transactor.QueryRowContext(ctx, query, args...)
}

// Commit commits the transaction
func (i hollywoodTx) Commit() error {
	return i.Transactor.Commit()
}

// Rollback aborts the transaction
func (i hollywoodTx) Rollback() error {
	return i.Transactor.Rollback()
}
