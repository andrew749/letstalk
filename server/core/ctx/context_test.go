package ctx_test

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/sessions"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/http"
)

func TestNewContext(t *testing.T) {
	writer := http.TestResponseWriter{}
	g, _ := gin.CreateTestContext(&writer)
	db := &gorm.DB{}
	sm := sessions.CreateCompositeSessionManager()
	sessionData, _ := sessions.CreateSessionData(1, nil, time.Now())
	c := ctx.NewContext(g, db, sessionData, &sm)
	assert.Equal(t, db, c.Db)
	assert.Equal(t, g, c.GinContext)
	assert.Equal(t, sessionData, c.SessionData)
	assert.Nil(t, c.Result)
}
