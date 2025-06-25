package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
)

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Root Route")
	w.Write([]byte("Hello Root Route"))
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// RAW BODY
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println("Error parsing request body", err)
			http.Error(w, "Error parsing request body", http.StatusInternalServerError)
			return
		}

		fmt.Println("Raw body: ", body)
		fmt.Println("Raw body: ", string(body))

		// Unmarshalling as struct
		//var student Student
		//err = json.Unmarshal(body, &student)
		//if err != nil {
		//	fmt.Println("Error marshalling json as struct", err)
		//	return
		//}
		//
		//fmt.Println("Student as struct: ", student)

		// Unmarshalling as map
		studentMap := make(map[string]any)

		err = json.Unmarshal(body, &studentMap)
		if err != nil {
			fmt.Println("Error marshalling json as map", err)
			return
		}

		fmt.Println("Student as map", studentMap)
	}

	fmt.Println("Hello Students Route")
	w.Write([]byte("Hello Students Route"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)

	switch r.Method {
	case http.MethodGet:
		fmt.Println("Hello Get Method on Teachers Route")
		w.Write([]byte("Hello Get Method on Teachers Route"))
		return
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		fmt.Println("Form: ", r.Form)
		// Preparing response
		response := make(map[string]any)

		for k, v := range r.Form {
			response[k] = v[0]
		}

		fmt.Println("Processed response: ", response)

		fmt.Println("Hello Post Method on Teachers Route")
		w.Write([]byte("Hello Post Method on Teachers Route"))
		return
	case http.MethodPut:
		fmt.Println("Hello Put Method on Teachers Route")
		w.Write([]byte("Hello Put Method on Teachers Route"))
		return
	case http.MethodDelete:
		fmt.Println("Hello Delete Method on Teachers Route")
		w.Write([]byte("Hello Delete Method on Teachers Route"))
		return
	case http.MethodPatch:
		fmt.Println("Hello Patch Method on Teachers Route")
		w.Write([]byte("Hello Patch Method on Teachers Route"))
		return
	}

	fmt.Println("Hello Teachers Route")
	w.Write([]byte("Hello Teachers Route"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Execs Route")
	w.Write([]byte("Hello Execs Route"))
}

func main() {
	port := 3000
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/students", studentsHandler)
	mux.HandleFunc("/teachers", teachersHandler)
	mux.HandleFunc("/execs", execsHandler)

	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		//Handler: middlewares.SecurityHeaders(mux),
		Handler: mw.Compression(mw.ResponseTime(mw.Cors(mux))),
	}
	fmt.Println("Server running on port: ", port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
