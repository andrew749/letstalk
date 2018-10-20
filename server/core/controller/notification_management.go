package controller

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"net/http"
)

const templateLink = "notification_console.html"

// GetNotificationManagementConsole Render the admin console to send batches of notifications
func GetNotificationManagementConsole(ctx *ctx.Context) errs.Error {
	// TODO(acod): load data about ongoing notifications
	ctx.GinContext.HTML(http.StatusOK, templateLink, &map[string]string{})
	return nil
}
