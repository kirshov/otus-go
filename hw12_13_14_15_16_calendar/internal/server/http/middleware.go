package internalhttp

import (
	"net/http"
	"strings"
)

func loggingMiddleware(app Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := strings.Builder{}
		msg.WriteString("HTTP request: ")
		msg.WriteString(r.RemoteAddr + " ")
		msg.WriteString(r.Method + " ")
		msg.WriteString(r.URL.String() + " ")
		msg.WriteString(r.Proto + " ")
		msg.WriteString(r.Header.Get("User-Agent"))
		app.GetLogger().Info(msg.String())

		next.ServeHTTP(w, r)
	})
}
