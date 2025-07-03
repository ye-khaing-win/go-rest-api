package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("POST /teachers/", handlers.CreateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	mux.HandleFunc("/execs/", handlers.ExecsHandler)
	mux.HandleFunc("/", handlers.RootHandler)

	return mux
}
