package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/students/", handlers.StudentsHandler)
	mux.HandleFunc("/teachers/", handlers.TeachersHandler)
	mux.HandleFunc("/execs/", handlers.ExecsHandler)
	mux.HandleFunc("/", handlers.RootHandler)

	return mux
}
