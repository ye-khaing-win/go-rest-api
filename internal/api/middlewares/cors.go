package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{
	"https://mydomain.com",
	"http://localhost:3000",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Cors Middleware starts...")
		origin := r.Header.Get("Origin")

		// Always send Vary: Origin so caches know
		w.Header().Add("Vary", "Origin")

		if !isOriginAllowed(origin) {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}
		// Set the CORS headers
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Max-Age", "600") // browser may cache preflight for 10 min

		// If this is a preflight request, we’re done
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Otherwise call the next handler
		next.ServeHTTP(w, r)
		fmt.Println("Cors Middleware ends...")
	})
}

func isOriginAllowed(origin string) bool {
	for _, ao := range allowedOrigins {
		if origin == ao {
			return true
		}
	}
	return false
}
