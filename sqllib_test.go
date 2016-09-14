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
	if !assert.NoError(t, lib.Register("select", "SELECT 1"), "Library.Register should succeed") {
		return
	}

	stmt, err := lib.GetStmt("select")
	if !assert.NoError(t, err, "Library.GetStmt should succeed") {
		return
	}
	_ = stmt
}
