package routes

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/login"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
)

type handlerWrapper struct {
	db *gmq.Db
}

type handlerFunc func(*ctx.Context) errs.Error

func Register(db *gmq.Db) *gin.Engine {
	hw := handlerWrapper{db}

	router := gin.Default()

	router.OPTIONS("/test")
	router.GET("/test", hw.wrapHandler(GetTest))

	v1 := router.Group("/v1")

	// create a new user
	v1.OPTIONS("/signup")
	v1.POST("/signup", hw.wrapHandler(login.SignupUser))

	// create a new session for an existing user
	v1.OPTIONS("/login")
	v1.GET("/login", hw.wrapHandler(login.GetLogin))

	// for fb_authentication
	v1.OPTIONS("/login_redirect")
	v1.GET("/login_redirect", hw.wrapHandler(login.GetLoginResponse))

	return router
}

func (hw handlerWrapper) wrapHandler(handler handlerFunc) gin.HandlerFunc {
	return func(g *gin.Context) {
		c := ctx.NewContext(g, hw.db)
		err := handler(c)

		if err != nil {
			rlog.Infof("Returning error: %s\n", err)
			c.GinContext.JSON(err.GetHTTPCode(), gin.H{"Error": convertError(err)})
			return
		}
		rlog.Infof("Returning result: %s\n", c.Result)
		c.GinContext.JSON(http.StatusOK, gin.H{"Result": c.Result})
	}
}

func convertError(e errs.Error) api.Error {
	return api.Error{
		Code:    e.GetHTTPCode(),
		Message: e.Error(),
	}
}
