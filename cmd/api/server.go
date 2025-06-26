package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
	"strconv"
	"strings"
	"sync"
)

type Teacher struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

var (
	teachers = make(map[int]Teacher)
	mutex    = sync.Mutex{}
	nextID   = 1
)

func init() {
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10B",
		Subject:   "History",
	}
	nextID++
}

func getTeacherHandler(id int, w http.ResponseWriter, r *http.Request) {

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher not found.", http.StatusNotFound)
		return
	}
	response := struct {
		Status string  `json:"status"`
		Data   Teacher `json:"data"`
	}{
		Status: "success",
		Data:   teacher,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	teachersList := make([]Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
			teachersList = append(teachersList, teacher)
		}

	}
	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachersList),
		Data:   teachersList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeacher Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeacher)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	newTeacher.ID = nextID
	teachers[nextID] = newTeacher
	nextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string  `json:"status"`
		Data   Teacher `json:"data"`
	}{
		Status: "success",
		Data:   newTeacher,
	}

	json.NewEncoder(w).Encode(response)
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

	switch r.Method {
	case http.MethodGet:
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		idStr := strings.TrimSuffix(path, "/")

		if idStr == "" {
			getTeachersHandler(w, r)
		} else {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return
			}
			getTeacherHandler(id, w, r)
		}

		//return
	case http.MethodPost:
		addTeacherHandler(w, r)
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

	//fmt.Println("Hello Teachers Route")
	//w.Write([]byte("Hello Teachers Route"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Execs Route")
	w.Write([]byte("Hello Execs Route"))
}

func main() {
	port := 3000
	mux := http.NewServeMux()

	mux.HandleFunc("/students/", studentsHandler)
	mux.HandleFunc("/teachers/", teachersHandler)
	mux.HandleFunc("/execs/", execsHandler)
	mux.HandleFunc("/", rootHandler)

	//rl := mw.NewRateLimiter(5, 10*time.Second)
	//hpp := mw.HPP{
	//	CheckQuery:      true,
	//	CheckBody:       true,
	//	BodyContentType: "application/x-www-form-urlencoded",
	//	Whitelist:       []string{"name", "age", "gender"},
	//}

	secureMux := mw.SecurityHeaders(mux)
	//secureMux := rl.Middleware(mw.ResponseTime(mw.SecurityHeaders(mw.Compression(hpp.Middleware()(mux)))))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: secureMux,
	}
	fmt.Println("Server running on port: ", port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
