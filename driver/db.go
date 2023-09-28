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

package driver

import (
	"database/sql"

	"github.com/blockysource/blockysql/bserr"
)

// DB is a driver specific wrapper over the database/sql.DB.
type DB interface {
	// Dialect returns the dialect of the driver.
	// I.e. if the driver implements a protocol that matches multiple
	// database dialects (such as postgres and cockroachdb), it should
	// return the dialect name.
	Dialect() string

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
