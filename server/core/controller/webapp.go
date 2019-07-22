package controller

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"net/http"
)

const webappTemplateLink = "webapp_home.html"

// GetWebapp Render the admin panel
func GetWebapp(ctx *ctx.Context) errs.Error {
	ctx.GinContext.HTML(http.StatusOK, webappTemplateLink, &map[string]string{})
	return nil
}
