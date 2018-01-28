package ctx

import (
	"letstalk/server/core/errs"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
)

type Context struct {
	GinContext *gin.Context
	Db         *gmq.Db
	Result     interface{}
	Errors     []errs.Error
}

func NewContext(g *gin.Context, db *gmq.Db) *Context {
	return &Context{
		GinContext: g,
		Db:         db,
		Errors:     make([]errs.Error, 0),
	}
}

func (c *Context) AddError(e errs.Error) {
	log.Println("Added error: ", e.Error())
	c.Errors = append(c.Errors, e)
}

func (c *Context) HasErrors() bool {
	return len(c.Errors) > 0
}
