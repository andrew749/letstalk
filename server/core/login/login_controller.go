package login

import (
	"bytes"
	"fmt"
	"letstalk/server/core/secrets"
	"net/http"
	"text/template"
	"letstalk/server/core/ctx"
)

var redirectLink = `https://www.facebook.com/v2.11/dialog/oauth?
  client_id={{.AppId}}
  &redirect_uri={{.RedirectUrl}}
  &state=`

type redirectData struct {
	AppId       string
	RedirectUrl string
	CSRFToken   string
}

func GetLogin(c *ctx.Context) {
	data := redirectData{
		AppId: secrets.GetSecrets().AppId,
		RedirectUrl: secrets.GetSecrets().RedirectUrl,
		CSRFToken: "TODO CSRF TOKEN GENERATION",
	}
	var link bytes.Buffer
	t, _ := template.New("login_template").Parse(redirectLink)
	t.Execute(&link, data)
	redirectUrl := link.String()
	fmt.Println(redirectUrl)

	c.GinContext.Redirect(http.StatusSeeOther, redirectUrl)
}
