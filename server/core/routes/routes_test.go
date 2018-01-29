package routes

import (
	"testing"

	"fmt"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	code "net/http"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
)

func TestHandlerResult(t *testing.T) {
	db := &gmq.Db{}
	hw := handlerWrapper{db}
	msg := "test message"
	handler := hw.wrapHandler(func(c *ctx.Context) {
		c.Result = msg
	})
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	handler(g)
	assert.Equal(t, code.StatusOK, writer.StatusCode)
	assert.Equal(t, fmt.Sprintf(`{"Result":"%s"}`, msg), writer.Output)
}

func TestHandlerError(t *testing.T) {
	db := &gmq.Db{}
	hw := handlerWrapper{db}
	msg := "test error message"
	handler := hw.wrapHandler(func(c *ctx.Context) {
		c.AddError(errs.NewClientError(msg))
	})
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	handler(g)
	assert.Equal(t, code.StatusBadRequest, writer.StatusCode)
	assert.Equal(t, fmt.Sprintf(`{"Errors":[{"Code":400,"Message":"%s"}]}`, msg), writer.Output)
}

func TestHandlerMultipleErrors(t *testing.T) {
	db := &gmq.Db{}
	hw := handlerWrapper{db}
	msg1, msg2 := "test error message 1", "test error message 2"
	handler := hw.wrapHandler(func(c *ctx.Context) {
		c.AddError(errs.NewClientError(msg1))
		c.AddError(errs.NewInternalError(msg2))
	})
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	handler(g)
	assert.Equal(t, code.StatusInternalServerError, writer.StatusCode)
	assert.Equal(t, fmt.Sprintf(`{"Errors":[{"Code":400,"Message":"%s"},{"Code":500,"Message":"%s"}]}`, msg1, msg2), writer.Output)
}
