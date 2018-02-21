# sqllib

Maintain a library of prepared SQL statements (\*sql.Stmt)

[![Build Status](https://travis-ci.org/lestrrat-go/sqllib.png?branch=master)](https://travis-ci.org/lestrrat-go/sqllib)

[![GoDoc](https://godoc.org/github.com/lestrrat-go/sqllib?status.svg)](https://godoc.org/github.com/lestrrat-go/sqllib)

# SYNOPSIS

```go
import (
  "github.com/lestrrat-go/sqllib"
  "github.com/pkg/errors"
)

var lib *sqllib.Library
var db *sql.DB

func InitializeDB() {
  db, _ = sql.Open(...)

  lib = sqllib.New(db)

  // Register some SQL queries by name
  lib.Register("Simple Select", "SELECT foo FROM bar WHERE a = ?")
}

func SomeFunc(tx *sql.Tx, arg string) error {
  // When you access the SQL query again, you can ask for an
  // already prepared statement.
  stmt, err := lib.GetStmt("Simple Select")
  if err != nil {
    return errors.Wrap(err, "failed to get statement")
  }

  // Don't forget to call (*sql.Tx).Stmt on it to make a 
  // transaction-specific statement
  rows, err := tx.Stmt(stmt).Query(arg)
  ...
}
```

# DESCRIPTION

Using prepared statements repeatedly is usually better for performance.

Keeping prepared statements around for reuse is fairly painful. This library
is a very small utility to store SQL queries and refer to them by name to
get back already prepared statement.