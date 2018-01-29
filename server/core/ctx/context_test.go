package ctx_test

import (
	"errors"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
)

func createTestContext() *ctx.Context {
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	db := gmq.Db{}
	return ctx.NewContext(g, &db)
}

func TestNewContext(t *testing.T) {
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	db := &gmq.Db{}
	c := ctx.NewContext(g, db)
	assert.Equal(t, db, c.Db)
	assert.Equal(t, g, c.GinContext)
	assert.Nil(t, c.Result)
	assert.Empty(t, c.Errors)
	assert.False(t, c.HasErrors())
}

func TestAddError(t *testing.T) {
	c := createTestContext()
	assert.False(t, c.HasErrors())
	const msg = "test message"
	err := errs.NewClientError(msg)
	c.AddError(err)
	assert.True(t, c.HasErrors())
	assert.Len(t, c.Errors, 1)
	assert.Equal(t, err.GetHTTPCode(), c.Errors[0].GetHTTPCode())
	assert.Equal(t, msg, c.Errors[0].Error())
}

func TestAddErrorMultiple(t *testing.T) {
	c := createTestContext()
	assert.False(t, c.HasErrors())
	messages := []string{"test message 1", "test message 2", "test message 3"}
	addErrs := []errs.Error{
		errs.NewClientError(messages[0]),
		errs.NewInternalError(messages[1]),
		errs.NewDbError(errors.New(messages[2])),
	}
	for _, err := range addErrs {
		c.AddError(err)
	}
	assert.True(t, c.HasErrors())
	assert.Len(t, c.Errors, 3)
	for i, err := range addErrs {
		assert.Equal(t, err.GetHTTPCode(), c.Errors[i].GetHTTPCode())
		assert.Equal(t, err.Error(), c.Errors[i].Error())
	}
}
