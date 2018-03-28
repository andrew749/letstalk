package routes

import (
	"testing"
	"time"

	"fmt"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/sessions"
	code "net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func setupSessionManager(t *testing.T) (*gorm.DB, sessions.ISessionManagerBase) {
	dbMock, _, err := sqlmock.New()
	if err != nil {
		assert.Fail(t, "Unable to create mock db.")
	}
	db, _ := gorm.Open("mysql", dbMock)
	sm := sessions.CreateSessionManager(db)
	return db, sm
}

func TestHandlerResult(t *testing.T) {
	db, sm := setupSessionManager(t)
	hw := handlerWrapper{db, &sm}
	msg := "test message"
	handler := hw.wrapHandler(func(c *ctx.Context) errs.Error {
		c.Result = msg
		return nil
	}, false)
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	handler(g)
	assert.Equal(t, code.StatusOK, writer.StatusCode)
	assert.Equal(t, fmt.Sprintf(`{"Result":"%s"}`, msg), writer.Output)
}

func TestHandlerAuthBad(t *testing.T) {
	db, sm := setupSessionManager(t)
	hw := handlerWrapper{db, &sm}
	msg := "test message"
	handler := hw.wrapHandler(func(c *ctx.Context) errs.Error {
		c.Result = msg
		return nil
	}, true)
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	g.Request, _ = code.NewRequest("POST", "/test", nil)
	g.Request.Header.Set("sessionId", "1")
	handler(g)
	assert.Equal(t, code.StatusUnauthorized, writer.StatusCode)
}

func TestHandlerAuthGood(t *testing.T) {
	db, sm := setupSessionManager(t)

	session, err := sm.CreateNewSessionForUserId(1, nil)
	assert.Nil(t, err)

	hw := handlerWrapper{db, &sm}
	msg := "test message"
	handler := hw.wrapHandler(func(c *ctx.Context) errs.Error {
		c.Result = msg
		return nil
	}, true)
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	g.Request, _ = code.NewRequest("POST", "/test", nil)
	g.Request.Header.Set("sessionId", *session.SessionId)
	handler(g)
	assert.Equal(t, code.StatusOK, writer.StatusCode)
	assert.Equal(t, fmt.Sprintf(`{"Result":"%s"}`, msg), writer.Output)
}

func TestHandlerExpiredToken(t *testing.T) {
	db, sm := setupSessionManager(t)

	session, err := sm.CreateNewSessionForUserIdWithExpiry(1, nil, time.Unix(0, 0))
	assert.Nil(t, err)

	hw := handlerWrapper{db, &sm}
	msg := "test message"
	handler := hw.wrapHandler(func(c *ctx.Context) errs.Error {
		c.Result = msg
		return nil
	}, true)
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	g.Request, _ = code.NewRequest("POST", "/test", nil)
	g.Request.Header.Set("sessionId", *session.SessionId)
	handler(g)
	assert.Equal(t, code.StatusUnauthorized, writer.StatusCode)
}

func TestHandlerClientError(t *testing.T) {
	db, sm := setupSessionManager(t)

	hw := handlerWrapper{db, &sm}
	msg := "test error message"
	handler := hw.wrapHandler(func(c *ctx.Context) errs.Error {
		return errs.NewClientError(msg)
	}, false)
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	handler(g)
	assert.Equal(t, code.StatusBadRequest, writer.StatusCode)
	assert.Equal(t, fmt.Sprintf(`{"Error":{"Code":400,"Message":"%s"}}`, msg), writer.Output)
}

func TestHandlerInternalError(t *testing.T) {
	db, sm := setupSessionManager(t)

	hw := handlerWrapper{db, &sm}
	msg := "test error message"
	handler := hw.wrapHandler(func(c *ctx.Context) errs.Error {
		return errs.NewInternalError(msg)
	}, false)
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	handler(g)
	assert.Equal(t, code.StatusInternalServerError, writer.StatusCode)
}
