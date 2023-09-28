# The Go database/sql wrapper (Blocky SQL)

_Write once, use with any SQL driver, without driver specific types_

[![GoDoc](https://godoc.org/github.com/blockysource/blockysql?status.svg)](https://godoc.org/github.com/blockysource/blockysql)
[![Go Report Card](https://goreportcard.com/badge/github.com/blockysource/blockysql)](https://goreportcard.com/report/github.com/blockysource/blockysql)
[![Build Status](https://travis-ci.org/blockysource/blockysql.svg?branch=master)](https://travis-ci.org/blockysource/blockysql)
[![codecov](https://codecov.io/gh/blockysource/blockysql/branch/master/graph/badge.svg)](https://codecov.io/gh/blockysource/blockysql)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

The Golang SQL wrapper with driver-specific parameters.

This implementation allows to use the same code for different SQL drivers.
It also provides a way to get driver-specific information, such as error codes and details.

The drivers could be opened using the OpenDB function, which returns a *blockysql.DB.

The *blockysql.DB is simply a wrapper around the *sql.DB,
which provides the same interface. It also provides a way to get driver-specific information.

What's more it allows to handle query errors without casting
the error to the driver-specific error type.

The *blockysql.DB also provides a way to run a function in a transaction,

## Usage

```go
package main

import (
    "context"
	
    "github.com/blockysource/blockysql"
    _ "github.com/blockysource/blockysql/pgxblockysql"
)

func main() {
    ctx := context.Background()
    db, err := blockysql.OpenDB(ctx, "mysql://user:password@host:port/dbname?param=value")
    if err != nil {
        panic(err)
    }
    
    // use db as usual
    rows, err := db.QueryContext(ctx, "SELECT * FROM table")
    if err != nil {
        panic(err)
    }
	defer rows.Close()
    
    // ...
    
    // close db
    db.Close()    
}
```

### Getting driver information.

```go main.go
package main 

import (
	"context"
	"fmt"
	
	"github.com/blockysource/blockysql"
    _ "github.com/blockysource/blockysql/pgxblockysql"      
)

func main() {
    ctx := context.Background()
    db, err := blockysql.OpenDB(ctx, "mysql://user:password@host:port/dbname?param=value")
    if err != nil {
        panic(err)
    }
    
    // Get the name of the driver.
	// Prints: 'pgx'
    fmt.Println(db.DriverName()) 

	
    // Get the dialect name from the driver.
	// It could be used to switch implementations based on the database family.
	// I.e. postgres specific queries.
    // Prints: 'postgres' or 'cockroach' (if database is cockroachdb).
    fmt.Println(db.Dialect())
    
    // close db
    db.Close()    
}
```

### Generic error codes and details.

```go main.go  
package main

import (
    "context"
    "fmt"
	
    "github.com/blockysource/blockysql"
    "github.com/blockysource/blockysql/bserr"
    _ "github.com/blockysource/blockysql/pgxblockysql"
)

func main() {
    ctx := context.Background()
    db, err := blockysql.OpenDB(ctx, "mysql://user:password@host:port/dbname?param=value")
    if err != nil {
        panic(err)
    }
    
    // use db as usual
    _, err = db.ExecContext(ctx, "INSERT INTO table (id, name) VALUES ($1, $2)", 1, "name")
    if err != nil {
        if db.ErrorCode(err) == bserr.UniqueViolation {
            // handle unique violation
            // i.e. return error to the user.
            if db.HasErrorDetails() {
                // get error details
                // i.e. get the column name that violated the unique constraint.                
                fmt.Println(db.ErrorColumn(err))
            }
        }        
    }  
    
    // close db
    db.Close()    
}
```

#### Getting *sql.DB from *blockysql.DB

```go main.go
package main

import (
	"context"
	
	"github.com/blockysource/blockysql"
	_ "github.com/blockysource/blockysql/pgxblockysql"
)

func main() {
    ctx := context.Background()
    db, err := blockysql.OpenDB(ctx, "mysql://user:password@host:port/dbname?param=value")
    if err != nil {
        panic(err)
    }
    
    // get *sql.DB
    sqlDB := db.DB()
    
    // use sqlDB as usual
	
	// close DB
    db.Close() 	
}
```

### Auto Commit/Rollback of transactions.

```go main.go
package main

import (
	"context"
	"database/sql"
	
	"github.com/blockysource/blockysql"
	_ "github.com/blockysource/blockysql/pgxblockysql"
)

func main() {
    ctx := context.Background()
    db, err := blockysql.OpenDB(ctx, "mysql://user:password@host:port/dbname?param=value")
    if err != nil {
        panic(err)
    }
    
    // The execFn function will be executed in a transaction.
	// It could alternatively be a closure, or a method.
    err = db.RunInTransaction(ctx, nil, execFn)
    if err != nil {
        // An error occurred, it can be handled here.
        // The transaction is already rolled back.
        // i.e. db.ErrorCode(err) == bserr.UniqueViolation...
    }
    
    // close db
    db.Close()
}

// execFn is the function that will be executed in a transaction.
func execFn(ctx context.Context, tx *sql.Tx) error {
    // use tx as usual
    _, err := tx.ExecContext(ctx, "INSERT INTO table (id, name) VALUES ($1, $2)", 1, "name")
    if err != nil {
        // Do not handle the transaction as rollback or commit will be called automatically.
        return err
    }
    
    // ...
    return nil
}
```