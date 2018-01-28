package routes

import (
	"net/http"
)

type LoginCallbackController struct {
	Controller
}

func (lc LoginCallbackController) GetPath() string {
	return "/login_succeed"
}

type CallbackResponse struct {
	Status string
	Code   string
	State  string
}

func (lc LoginCallbackController) Handler(
	res http.ResponseWriter,
	req *http.Request) {

	status := req.URL.Query().Get("status")
	code := req.URL.Query().Get("code")
	state := req.URL.Query().Get("state")

	// authenticated data from the provider
	_ = CallbackResponse{status, code, state}
}
