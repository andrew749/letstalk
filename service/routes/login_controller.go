package routes

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"
	"uwletstalk/service/secrets"
)

type LoginController struct {
	Controller
}

var redirectLink = `https://www.facebook.com/v2.11/dialog/oauth?
  client_id={{.AppId}}
  &redirect_uri={{.RedirectUrl}}
  &state=`

func (lc LoginController) GetPath() string {
	return "/login"
}

type LoginTemplateData struct {
	AppId       string
	RedirectUrl string
	CSRFToken   string
}

func (lc LoginController) Handler(res http.ResponseWriter, req *http.Request) {
	data := LoginTemplateData{
		secrets.GetSecrets().AppId,
		secrets.GetSecrets().RedirectUrl,
		"TODO CSRF TOKEN GENERATION",
	}
	var link bytes.Buffer
	t, _ := template.New("login_template").Parse(redirectLink)
	t.Execute(&link, data)
	redirectUrl := link.String()
	fmt.Println(redirectUrl)
	http.Redirect(res, req, redirectUrl, 301)
}
