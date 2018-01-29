package ctx

import (
	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
)

type Context struct {
	GinContext *gin.Context
	Db         *gmq.Db
	Result     interface{}
}

func NewContext(g *gin.Context, db *gmq.Db) *Context {
	return &Context{
		GinContext: g,
		Db:         db,
	}
}
