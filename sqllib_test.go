package sqllib_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	sqllib "github.com/lestrrat-go/sqllib"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := NewMockDB(ctrl)
	lib := sqllib.New(db)

	db.EXPECT().Prepare("SELECT 1")
	var key = &struct{}{}
	var badkey = &struct{}{}
	if !assert.NoError(t, lib.Register(key, "SELECT 1"), "Library.Register should succeed") {
		return
	}

	stmt, err := lib.GetStmt(key)
	if !assert.NoError(t, err, "Library.GetStmt should succeed") {
		return
	}
	_ = stmt

	stmt, err = lib.GetStmt(badkey)
	if !assert.Error(t, err, "Library.GetStmt should fail") {
		return
	}
	_ = stmt

}
