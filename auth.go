package session

import (
	"net/http"
)

var (
	RedirectUrl      string = "/account/login"
	AdminRedirectUrl string = "/admin/account/login"
	RedirectParam    string = "next"
)

// User
type User interface {
	Login() int
	Logout() int
	IsAdmin() bool
	IsAuthenticated() bool
	UniqueId() interface{}
	GetById(id interface{}) (User, int)
}

func LoginRequired(user User, req *http.Request, resp http.ResponseWriter) {
	if !user.IsAuthenticated() {
		//path := fmt.Sprintf("%s?%s=%s", RedirectUrl, RedirectParam, req.URL.Path)
		//http.Redirect(resp, req, path, http.StatusFound)
		http.Error(resp, "", http.StatusUnauthorized)
	}
}

func AdminRequired(user User, req *http.Request, resp http.ResponseWriter) {
	if !user.IsAuthenticated() || !user.IsAdmin() {
		//path := fmt.Sprintf("%s?%s=%s", AdminRedirectUrl, RedirectParam, req.URL.Path)
		//http.Redirect(resp, req, path, http.StatusFound)
		http.Error(resp, "", http.StatusUnauthorized)
	}
}
