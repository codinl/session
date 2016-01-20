package session

import (
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Handler(t *testing.T) {
	m := martini.Classic()

	m.Use(Handler())

	m.Get("(/)", func() {})

	handler := func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		for _, c := range cookies {
			if c.Name == COOKIE_NAME {
				fmt.Println(c.Value)
			}
		}
	}

	req, err := http.NewRequest("GET", "http://localhost:3000", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	m.ServeHTTP(w, req)
}
