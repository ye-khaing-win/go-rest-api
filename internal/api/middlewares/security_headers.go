package middlewares

import (
	"fmt"
	"net/http"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SecurityHeaders Middleware starts...")
		w.Header().Set("X-Powered-By", "Express")
		next.ServeHTTP(w, r)
		fmt.Println("SecurityHeaders Middleware ends...")

	})
}
