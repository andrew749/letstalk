package ctx

import (
	"letstalk/server/core/sessions"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
)

type Context struct {
	GinContext     *gin.Context
	Db             *gmq.Db
	SessionData    *sessions.SessionData
	SessionManager *sessions.ISessionManager
	Result         interface{}
}

func NewContext(
	g *gin.Context,
	db *gmq.Db,
	sessionData *sessions.SessionData,
	sm *sessions.ISessionManager,
) *Context {
	return &Context{
		GinContext:     g,
		Db:             db,
		SessionData:    sessionData,
		SessionManager: sm,
	}
}
