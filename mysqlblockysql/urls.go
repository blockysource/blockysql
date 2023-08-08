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

package mysqlblockysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/go-sql-driver/mysql"

	"github.com/blockysource/blockysql"
	"github.com/blockysource/blockysql/bserr"
	"github.com/blockysource/blockysql/driver"
)

func init() {
	// Register the driver.
	blockysql.DefaultURLMux().RegisterDB(SchemaName, new(URLOpener))
}

// SchemaName is the name of the driver.
const SchemaName = "mysql"

var _ blockysql.DBURLOpener = (*URLOpener)(nil)

// URLOpener is a mysql driver url opener.
type URLOpener struct {
	// TraceOpts contains options for OpenCensus.
	TraceOpts []ocsql.TraceOption
}

// OpenDBURL opens a new database connection for the given URL.
func (o *URLOpener) OpenDBURL(ctx context.Context, u *url.URL) (*blockysql.DB, error) {
	cfg, err := mysql.ParseDSN(u.String())
	if err != nil {
		return nil, fmt.Errorf("mysqlblockysql: open database failed: %v", err)
	}

	return OpenDB(ctx, cfg, Options{TraceOpts: o.TraceOpts})
}

type Options struct {
	// TraceOpts contains options for OpenCensus.
	TraceOpts []ocsql.TraceOption
}

// OpenDB opens a new database connection and returns
// a *blockysql.DB wrapper.
func OpenDB(ctx context.Context, cfg *mysql.Config, opts Options) (*blockysql.DB, error) {
	c, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, fmt.Errorf("mysqlblockysql: open database failed: %v", err)
	}

	db := sql.OpenDB(ocsql.WrapConnector(c, opts.TraceOpts...))

	d := &DB{db: db}

	return blockysql.NewDB(d)
}

// Compile time check if URLOpener implements blockysql.DBURLOpener.
var _ driver.DB = (*DB)(nil)

// DB is the driver for the PostgreSQL database.
type DB struct {
	db *sql.DB
}

// DriverName implements driver.DB
func (d *DB) DriverName() string {
	return "mysql"
}

// FamilyName implements driver.DB
func (d *DB) FamilyName() string {
	return driver.FamilyMySQL
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

	var derr *mysql.MySQLError
	if !errors.As(err, &derr) {
		return bserr.Unknown
	}

	switch derr.Number {
	case 1040:
		return bserr.TooManyConnections
	case 1045:
		return bserr.PermissionDenied
	case 1022, 1062, 1169:
		return bserr.UniqueViolation
	case 1064, 1065, 1067, 1072, 1087, 1090:
		return bserr.InvalidInputSyntax
	case 1146:
		return bserr.TableNotFound
		// Foreign key violation
	case 1212, 1216, 1217, 1451, 1452, 1453, 1557, 1825, 1826:
		return bserr.ForeignKeyViolation
	case 3819, 3820, 3822, 3823:
		return bserr.ConstraintViolation
	}

	sqlState := string(derr.SQLState[:])
	switch sqlState {
	case "23000":
		return bserr.UniqueViolation
	case "23502":
		return bserr.ConstraintViolation
	case "23503":
		return bserr.ForeignKeyViolation
	case "23505":
		return bserr.UniqueViolation
	case "23514":
		return bserr.ConstraintViolation
	case "42000":
		return bserr.InvalidInputSyntax
	case "42S02":
		return bserr.TableNotFound
	case "HY000":
		return bserr.InternalError
	case "HY001":
		return bserr.TooManyConnections
	case "HYT00":
		return bserr.Timeout
	case "HYT01":
		return bserr.TooManyConnections
	case "HYT02":
		return bserr.TooManyConnections
	case "42102":
		return bserr.TableNotFound
	}

	if strings.HasPrefix(sqlState, "22") {
		return bserr.ConstraintViolation
	}

	if strings.HasPrefix(sqlState, "23") {
		return bserr.ConstraintViolation
	}

	return bserr.Unknown
}

// HasErrorDetails returns true if the driver supports error details,
// such as column, table and constraint name.
func (d *DB) HasErrorDetails() bool {
	return false
}

// ErrorColumn returns the column name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorColumn(err error) string {
	if err == nil {
		return ""
	}
	return ""
}

// ErrorTable returns the table name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorTable(err error) string {
	if err == nil {
		return ""
	}

	return ""
}

// ErrorConstraint returns the constraint name of the given error
// if the driver doesn't support it, it should return an empty string.
func (d *DB) ErrorConstraint(err error) string {
	if err == nil {
		return ""
	}

	return ""
}

// DB returns the underlying database/sql.DB.
func (d *DB) DB() *sql.DB {
	return d.db
}
