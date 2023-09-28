// Copyright 2023 The Blocky Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pgxblockysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/blockysource/blockysql"
	"github.com/blockysource/blockysql/bserr"
	"github.com/blockysource/blockysql/driver"
)

func init() {
	// Register the driver.
	blockysql.DefaultURLMux().RegisterDB(SchemaName, new(URLOpener))
}

// SchemaName is the name of the driver.
const SchemaName = "pgx"

var _ blockysql.DBURLOpener = (*URLOpener)(nil)

// URLOpener opens PostgreSQL URLs like "pgx://user:password@host/dbname?option=value".
// For details see https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING.
// This driver supports runtime configurable parameters as defined in
// https://www.postgresql.org/docs/current/runtime-config.html.
// I.e.:
//   - search_path (string) - Sets the schema search order for names that are not schema-qualified.
//   - timezone (string) - Sets the time zone for displaying and interpreting time stamps.
type URLOpener struct {
	// TraceOpts contains options for OpenCensus.
	TraceOpts []ocsql.TraceOption

	// OpenDBOpts contains options for the pgx driver.
	OpenDBOpts []stdlib.OptionOpenDB
}

// OpenDBURL opens a new database connection for the given URL.
func (o *URLOpener) OpenDBURL(ctx context.Context, u *url.URL) (*blockysql.DB, error) {
	// Make a copy of the url to avoid modifying the original.
	pqu := url.URL{
		Scheme:      "postgres", // postgres:// as the scheme is required by pq.
		Opaque:      u.Opaque,
		User:        u.User, // user:password
		Host:        u.Host, // host:port
		Path:        u.Path, // dbname
		RawPath:     u.RawPath,
		OmitHost:    u.OmitHost,
		ForceQuery:  u.ForceQuery,
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
	}

	// Parse query parameters.
	cfg, err := pgx.ParseConfig(pqu.String())
	if err != nil {
		return nil, fmt.Errorf("pgxblockysql: open database failed: %v", err)
	}

	// Open database.
	return OpenDB(ctx, *cfg, Options{
		TraceOpts:  o.TraceOpts,
		OpenDBOpts: o.OpenDBOpts,
	})
}

// Options contains options for configuring the database connection.
type Options struct {
	// TraceOpts contains options for OpenCensus.
	TraceOpts []ocsql.TraceOption

	// OpenDBOpts contains options for configuring the database connection.
	OpenDBOpts []stdlib.OptionOpenDB
}

// Compile time check if URLOpener implements blockysql.DBURLOpener.
var _ driver.DB = (*DB)(nil)

// OpenDB returns a new bsql.DB backed by a *sql.DB.
func OpenDB(ctx context.Context, c pgx.ConnConfig, opts Options) (*blockysql.DB, error) {
	// Open database connection.
	db := sql.OpenDB(ocsql.WrapConnector(
		stdlib.GetConnector(c, opts.OpenDBOpts...),
		opts.TraceOpts...),
	)

	// Check if the database is CockroachDB.
	row := db.QueryRowContext(ctx, "SELECT VERSION()")
	var version string
	if err := row.Scan(&version); err != nil {
		return nil, fmt.Errorf("pgxblockysql: open database failed: %v", err)
	}

	d := &DB{db: db}
	if isCockroach := strings.Contains(version, "CockroachDB"); isCockroach {
		d.dialect = driver.DialectCockroach
	} else if isYugaByte := strings.Contains(version, "-YB-"); isYugaByte {
		d.dialect = driver.DialectYugabyte
	} else {
		d.dialect = driver.DialectPostgres
	}

	return blockysql.NewDB(d)
}

// DB is the driver for the PostgreSQL database.
type DB struct {
	db      *sql.DB
	dialect string
}

// DriverName implements driver.DB
func (d *DB) DriverName() string {
	return "pgx"
}

// Dialect implements driver.DB
func (d *DB) Dialect() string {
	return d.dialect
}

// ErrorCode implements driver.DB.
func (d *DB) ErrorCode(err error) bserr.Code {
	if err == nil {
		return bserr.OK
	}

	if err == sql.ErrNoRows {
		return bserr.NotFound
	}

	if err == sql.ErrTxDone {
		return bserr.TxDone
	}

	if err == context.DeadlineExceeded {
		return bserr.Timeout
	}

	var pgerr *pgconn.PgError
	if !errors.As(err, &pgerr) {
		return bserr.Unknown
	}

	switch pgerr.Code {
	case "23505":
		return bserr.UniqueViolation
	case "23503":
		return bserr.ForeignKeyViolation
	case "23502":
		return bserr.ConstraintViolation
	case "23514":
		return bserr.ConstraintViolation
	case "42501":
		return bserr.PermissionDenied
	case "40001": // CockroachDB serialization failure.
		return bserr.ConcurrentUpdate
	case "CR000":
		return bserr.ConcurrentUpdate
	case "42P01":
		return bserr.TableNotFound
	case "53200":
		return bserr.OutOfMemory
	case "53100":
		return bserr.OutOfDisk
	case "53300":
		return bserr.TooManyConnections
	case "57000":
		return bserr.Timeout
	case "57P02":
		return bserr.Timeout
	case "XX000":
		return bserr.InternalError
	}
	if strings.HasPrefix(pgerr.Code, "42") {
		return bserr.InvalidInputSyntax
	}

	if strings.HasPrefix(pgerr.Code, "22") {
		return bserr.DataException
	}

	if strings.HasPrefix(pgerr.Code, "23") {
		return bserr.ConstraintViolation
	}

	if strings.HasPrefix(pgerr.Code, "54") {
		return bserr.InternalError // Program limit exceeded.
	}

	if strings.HasPrefix(pgerr.Code, "XX") {
		return bserr.InternalError
	}

	return bserr.Unknown
}

// HasErrorDetails returns true if the driver supports error details,
// such as column, table and constraint name.
func (d *DB) HasErrorDetails() bool {
	return true
}

// ErrorColumn returns the column name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorColumn(err error) string {
	if err == nil {
		return ""
	}

	var pgerr *pgconn.PgError
	if !errors.As(err, &pgerr) {
		return ""
	}
	return pgerr.ColumnName
}

// ErrorTable returns the table name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorTable(err error) string {
	if err == nil {
		return ""
	}

	var pgerr *pgconn.PgError
	if !errors.As(err, &pgerr) {
		return ""
	}
	return pgerr.TableName
}

// ErrorConstraint returns the constraint name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorConstraint(err error) string {
	if err == nil {
		return ""
	}

	var pgerr *pgconn.PgError
	if !errors.As(err, &pgerr) {
		return ""
	}
	return pgerr.ConstraintName
}

// DB returns the underlying database/sql.DB.
func (d *DB) DB() *sql.DB {
	return d.db
}
