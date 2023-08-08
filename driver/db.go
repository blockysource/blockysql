package driver

import (
	"database/sql"

	"github.com/blockysource/blockysql/bserr"
)

// DB is a driver specific wrapper over the database/sql.DB.
type DB interface {
	// FamilyName returns the family name of the driver.
	// I.e. "postgres" for the postgres driver.
	FamilyName() string

	// DriverName returns the name of the driver.
	DriverName() string

	// ErrorCode returns the error code of the given error
	// if the driver supports it.
	ErrorCode(err error) bserr.Code

	// DB returns the underlying database/sql.DB.
	DB() *sql.DB

	// ErrorColumn returns the column name of the given error
	// if the driver doesn't support it, it should return an empty string.
	ErrorColumn(err error) string

	// ErrorTable returns the table name of the given error
	// if the driver doesn't support it, it should return an empty string.
	ErrorTable(err error) string

	// ErrorConstraint returns the constraint name of the given error
	// if the driver doesn't support it, it should return an empty string.
	ErrorConstraint(err error) string

	// HasErrorDetails returns true if the driver supports error details,
	// such as column, table and constraint name.
	HasErrorDetails() bool
}
