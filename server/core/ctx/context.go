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
