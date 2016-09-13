package sqllib_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	sqllib "github.com/lestrrat/go-sqllib"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	lib := sqllib.New(db)

	db.EXPECT().Prepare("SELECT 1")
	key, err := lib.Register("SELECT 1")
	if !assert.NoError(t, err, "Library.Register should succeed") {
		return
	}

	stmt, err := lib.GetStmt(key)
	if !assert.NoError(t, err, "Library.GetStmt should succeed") {
		return
	}
	_ = stmt
}
