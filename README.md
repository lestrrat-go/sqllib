# sqllib

Maintain a library of prepared SQL statements (\*sql.Stmt)

# SYNOPSIS

```go
import (
  "github.com/lestrrat/go-sqllib"
  "github.com/pkg/errors"
)

var lib *sqllib.Library
var db *sql.DB

var queryKey sqllib.Key
func InitializeDB() {
  db, _ = sql.Open(...)

  lib = sqllib.New(db)

  queryKey = lib.Register("SELECT foo FROM bar WHERE a = ?")
}

func SomeFunc(tx *sql.Tx, arg string) error {
  stmt, err := lib.Get(queryKey)
  if err != nil {
    return errors.Wrap(err, "failed to get statement")
  }

  rows, err := tx.Stmt(stmt).Query(arg)
  ...
}
```