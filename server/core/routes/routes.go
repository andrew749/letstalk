package routes

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/users"
	"net/http"

	"letstalk/server/core/login"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
)

type handlerWrapper struct {
	db *gmq.Db
}

func Register(db *gmq.Db) *gin.Engine {
	hw := handlerWrapper{db}

	router := gin.Default()

	router.OPTIONS("/test")
	router.GET("/test", hw.wrapHandler(GetTest))

	v1 := router.Group("/v1")

	v1.OPTIONS("/users")
	v1.GET("/users", hw.wrapHandler(users.GetUsers))
	v1.POST("/users", hw.wrapHandler(users.PostUser))

	v1.OPTIONS("/login")
	v1.GET("/login", hw.wrapHandler(login.GetLogin))

	v1.OPTIONS("/login_succeed")
	v1.POST("/login_succeed", hw.wrapHandler(login.PostLoginSucceed))

	return router
}

func (hw handlerWrapper) wrapHandler(handler func(*ctx.Context)) func(*gin.Context) {
	return func(gCtx *gin.Context) {
		c := &ctx.Context{GinContext: gCtx, Db: hw.db}
		handler(c)
		// TODO(aklen): handle errors in context
		if c.Result != nil {
			c.GinContext.JSON(http.StatusOK, c.Result) // Encode json response.
		}
	}
}
