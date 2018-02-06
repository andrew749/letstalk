package ctx

import (
	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
)

type Context struct {
	GinContext  *gin.Context
	Db          *gmq.Db
	SessionData *SessionData
	Result      interface{}
}

func NewContext(g *gin.Context, db *gmq.Db, sessionData *SessionData) *Context {
	return &Context{
		GinContext:  g,
		Db:          db,
		SessionData: sessionData,
	}
}
