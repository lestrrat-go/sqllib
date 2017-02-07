//go:generate mockgen -source sqllib.go -destination mockgen_test.go -package sqllib_test

package sqllib

import (
	"database/sql"
	"sync"

	"github.com/pkg/errors"
)

// DB is declared so we can mock-test. You do not really have to
// think about this. Just assume it's the value returned from
// database/sql.Open
type DB interface {
	Prepare(string) (*sql.Stmt, error)
}

// Entry is an SQL statement registered to the library. This holds
// a reference to the prepared SQL statement (*sql.Stmt), but the
// statement is only prepared lazily.
type Entry struct {
	mutex sync.Mutex
	sql   string
	stmt  *sql.Stmt
}

// Library represents the top-level structure that holds the
// SQL statements.
type Library struct {
	db    DB
	mutex sync.RWMutex
	stmts map[interface{}]*Entry
}

func New(db DB) *Library {
	return &Library{
		db:    db,
		stmts: make(map[interface{}]*Entry),
	}
}

func (l *Library) Register(key interface{}, sql string) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if _, ok := l.stmts[key]; ok {
		return errors.New("duplicate key found")
	}
	l.stmts[key] = &Entry{sql: sql}
	return nil
}

func (l *Library) GetStmt(key interface{}) (*sql.Stmt, error) {
	l.mutex.RLock()
	e, ok := l.stmts[key]
	l.mutex.RUnlock()

	if !ok {
		return nil, errors.New("statement not found")
	}

	stmt, err := e.prepare(l.db)
	if err != nil {
		return nil, errors.Wrap(err, "lazy-prepare failed")
	}

	return stmt, nil
}

func (e *Entry) prepare(db DB) (*sql.Stmt, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.stmt == nil {
		stmt, err := db.Prepare(e.sql)
		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare statement")
		}

		e.stmt = stmt
	}
	return e.stmt, nil
}
