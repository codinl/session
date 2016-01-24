package session

import (
	"fmt"
	"net/http"
)

var (
	RedirectUrl      string = "/account/login"
	AdminRedirectUrl string = "/admin/account/login"
	RedirectParam    string = "next"
)

// User
type User interface {
	Login() error
	Logout() error
	IsAdmin() bool
	IsAuthenticated() bool
	UniqueId() interface{}
	GetById(id interface{}) (User, error)
}

func LoginRequired(user User, req *http.Request, resp http.ResponseWriter) {
	if !user.IsAuthenticated() {
		path := fmt.Sprintf("%s?%s=%s", RedirectUrl, RedirectParam, req.URL.Path)
		http.Redirect(resp, req, path, http.StatusFound)
	}
}

func AdminRequired(user User, req *http.Request, resp http.ResponseWriter) {
	if !user.IsAuthenticated() || !user.IsAdmin() {
		path := fmt.Sprintf("%s?%s=%s", AdminRedirectUrl, RedirectParam, req.URL.Path)
		http.Redirect(resp, req, path, http.StatusFound)
	}
}
