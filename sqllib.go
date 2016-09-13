//go:generate mockgen -source sqllib.go -destination mockgen_test.go -package sqllib_test


package sqllib

import (
	"crypto/sha256"
	"database/sql"
	"sync"

	"github.com/pkg/errors"
)

func makeKey(b []byte) Key {
	h := sha256.Sum256(b)
	k := Key{}
	for i := 0; i < sha256.Size; i++ {
		k[i] = h[i]
	}
	return Key(k)
}

type Key [sha256.Size]byte
type DB interface {
	Prepare(string) (*sql.Stmt, error)
}
type Entry struct {
	mutex sync.Mutex
	sql   string
	stmt  *sql.Stmt
}
type Library struct {
	db    DB
	mutex sync.RWMutex
	stmts map[Key]*Entry
}

func New(db DB) *Library {
	return &Library{
		db:    db,
		stmts: make(map[Key]*Entry),
	}
}

func (l *Library) Register(sql string) Key {
	k := makeKey([]byte(sql))

	l.mutex.Lock()
	l.stmts[k] = &Entry{sql: sql}
	l.mutex.Unlock()
	return k
}

func (l *Library) GetStmt(k Key) (*sql.Stmt, error) {
	l.mutex.RLock()
	e, ok := l.stmts[k]
	l.mutex.RUnlock()

	if !ok {
		return nil, errors.New("statement not found")
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.stmt == nil {
		stmt, err := l.db.Prepare(e.sql)
		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare statement")
		}

		e.stmt = stmt
	}
	return e.stmt, nil
}
