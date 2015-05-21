package middle

import (
	hand "github.com/squiidz/stamp/Web/module/handler"
	"net/http"
)

func AuthMiddle(next http.Handler, cookie string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		username := hand.CookieValue(req, cookie)
		if hand.UserExist(username) {
			next.ServeHTTP(rw, req)
		} else {
			http.Redirect(rw, req, "/", http.StatusFound)
		}
	})
}
