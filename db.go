package blockysql

import (
	"context"
	"database/sql"
	sqldriver "database/sql/driver"
	"errors"
	"time"

	"github.com/blockysource/blockysql/bserr"
	"github.com/blockysource/blockysql/driver"
)

// DB is a driver specific wrapper over the database/sql.DB.
type DB struct {
	driver driver.DB
	db     *sql.DB
}

// NewDB creates a new instance of DB.
var NewDB = newDB

// newDB creates a new instance of DB.
// The driver is used to determine the error code of the given error.
func newDB(d driver.DB) (*DB, error) {
	if d == nil {
		return nil, errors.New("blockysql: driver is nil")
	}

	if d.DB() == nil {
		return nil, errors.New("blockysql: driver.DB() is nil")
	}

	return &DB{
		driver: d,
		db:     d.DB(),
	}, nil
}

// DriverName returns the name of the driver.
func (d *DB) DriverName() string {
	return d.driver.DriverName()
}

// DatabaseFamily returns the family name of the driver.
// I.e. "postgres" for the postgres driver.
func (d *DB) DatabaseFamily() string {
	return d.driver.FamilyName()
}

// ErrorCode returns the error code of the given error
// if the driver supports it.
func (d *DB) ErrorCode(err error) bserr.Code {
	return d.driver.ErrorCode(err)
}

// HasErrorDetails returns true if the driver supports error details,
// such as column, table and constraint name.
func (d *DB) HasErrorDetails() bool {
	return d.driver.HasErrorDetails()
}

// ErrorColumn returns the column name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorColumn(err error) string {
	return d.driver.ErrorColumn(err)
}

// ErrorTable returns the table name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorTable(err error) string {
	return d.driver.ErrorTable(err)
}

// ErrorConstraint returns the constraint name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorConstraint(err error) string {
	return d.driver.ErrorConstraint(err)
}

// DB returns the underlying database/sql.DB.
func (d *DB) DB() *sql.DB {
	return d.db
}

// Begin starts a transaction.
func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

// BeginTx starts a transaction with the provided options.
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

// Close closes the database and prevents new queries from starting.
func (d *DB) Close() error {
	return d.db.Close()
}

// Conn returns a single connection by either opening a new connection
// or returning an existing connection from the connection pool.
func (d *DB) Conn(ctx context.Context) (*sql.Conn, error) {
	return d.db.Conn(ctx)
}

// SQLDriver returns the underlying database/sql/driver.Driver.
// This is useful for using the database/sql package directly.
func (d *DB) SQLDriver() sqldriver.Driver {
	return d.db.Driver()
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (d *DB) Exec(query string, args ...any) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (d *DB) Ping() error {
	return d.db.Ping()
}

// PingContext verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (d *DB) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method when the statement is no longer needed.
func (d *DB) Prepare(query string) (*sql.Stmt, error) {
	return d.db.Prepare(query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method when the statement is no longer needed.
func (d *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return d.db.PrepareContext(ctx, query)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (d *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (d *DB) QueryRow(query string, args ...any) *sql.Row {
	return d.db.QueryRow(query, args...)
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (d *DB) SetConnMaxLifetime(dur time.Duration) {
	d.db.SetConnMaxLifetime(dur)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func (d *DB) SetMaxIdleConns(n int) {
	d.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (d *DB) SetMaxOpenConns(n int) {
	d.db.SetMaxOpenConns(n)
}

// Stats returns database statistics.
func (d *DB) Stats() sql.DBStats {
	return d.db.Stats()
}

// RunInTransaction runs the given function in a transaction.
func (d *DB) RunInTransaction(ctx context.Context, opts *sql.TxOptions, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := d.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	if err = fn(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
