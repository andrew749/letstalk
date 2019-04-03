package controller

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"net/http"
)

const adminPanelTemplateLink = "admin_panel.html"

// GetAdminPanel Render the admin panel
func GetAdminPanel(ctx *ctx.Context) errs.Error {
	ctx.GinContext.HTML(http.StatusOK, adminPanelTemplateLink, &map[string]string{})
	return nil
}
