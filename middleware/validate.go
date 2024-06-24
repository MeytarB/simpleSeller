package middleware

import (
	"fmt"
	"net/http"
)

func Validate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				fmt.Println("middleware: Method not allowed")
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			fmt.Println("middleware: method is valid")
			next.ServeHTTP(w, r)
		})
	}
}
