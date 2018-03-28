package ctx

import (
	"letstalk/server/core/sessions"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Context struct {
	GinContext     *gin.Context
	Db             *gorm.DB
	SessionData    *sessions.SessionData
	SessionManager *sessions.ISessionManagerBase
	Result         interface{}
}

func NewContext(
	g *gin.Context,
	db *gorm.DB,
	sessionData *sessions.SessionData,
	sm *sessions.ISessionManagerBase,
) *Context {
	return &Context{
		GinContext:     g,
		Db:             db,
		SessionData:    sessionData,
		SessionManager: sm,
	}
}
