package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HPP struct {
	CheckBody       bool
	CheckQuery      bool
	BodyContentType string
	Whitelist       []string
}

func (hpp *HPP) isCorrectContentType(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), hpp.BodyContentType)
}

func (hpp *HPP) isWhitelisted(param string) bool {
	for _, v := range hpp.Whitelist {
		if param == v {
			return true
		}
	}
	return false
}

func (hpp *HPP) filterBodyParams(r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}

	for k, v := range r.Form {
		if len(v) > 1 {
			r.Form.Set(k, v[0])
		}
		if !hpp.isWhitelisted(k) {
			r.Form.Del(k)
		}
	}
}

func (hpp *HPP) filterQueryParams(r *http.Request) {
	query := r.URL.Query()

	for k, v := range query {
		if len(v) > 1 {
			query.Set(k, v[0])
		}
		if !hpp.isWhitelisted(k) {
			query.Del(k)
		}
	}

	r.URL.RawQuery = query.Encode()
}

func (hpp *HPP) Middleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Hpp Middleware starts...")
			if hpp.CheckBody && r.Method == http.MethodPost && hpp.isCorrectContentType(r) {
				hpp.filterBodyParams(r)
			}

			if hpp.CheckQuery && r.URL.Query() != nil {
				hpp.filterQueryParams(r)
			}
			next.ServeHTTP(w, r)

			fmt.Println("Hpp Middleware ends...")
		})
	}
}
