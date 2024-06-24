package middleware

import (
	"fmt"
	"net/http"
)

func MiddlewareLog(bidType string) func(http.Handler) http.Handler {
	fmt.Println("from middleware")
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("got request from %s for %s", r.URL.Path, bidType)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
