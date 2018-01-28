package routes

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/users"
	"net/http"

	"letstalk/server/core/login"

	"letstalk/server/core/api"
	"letstalk/server/core/errs"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
)

type handlerWrapper struct {
	db *gmq.Db
}

type handlerFunc func(*ctx.Context)

func Register(db *gmq.Db) *gin.Engine {
	hw := handlerWrapper{db}

	router := gin.Default()

	router.OPTIONS("/test")
	router.GET("/test", hw.wrapHandler(GetTest))

	v1 := router.Group("/v1")

	v1.OPTIONS("/users")
	v1.POST("/users", hw.wrapHandler(users.PostUser))

	v1.OPTIONS("/login")
	v1.GET("/login", hw.wrapHandler(login.GetLogin))

	v1.OPTIONS("/login_redirect")
	v1.GET("/login_redirect", hw.wrapHandler(login.GetLoginResponse))

	return router
}

func (hw handlerWrapper) wrapHandler(handler handlerFunc) gin.HandlerFunc {
	return func(g *gin.Context) {
		c := ctx.NewContext(g, hw.db)
		handler(c)

		if c.HasErrors() {
			errors := make([]api.Error, 0, len(c.Errors))
			code := http.StatusOK
			for _, err := range c.Errors {
				code = err.GetHTTPCode()
				errors = append(errors, convertError(err))
			}
			c.GinContext.JSON(code, gin.H{"Errors": errors})
			return
		}
		c.GinContext.JSON(http.StatusOK, gin.H{"Result": c.Result})
	}
}

func convertError(e errs.Error) api.Error {
	return api.Error{
		Code:    e.GetHTTPCode(),
		Message: e.Error(),
	}
}
