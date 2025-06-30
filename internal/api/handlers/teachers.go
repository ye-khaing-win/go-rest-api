package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"restapi/internal/models"
	"restapi/repository/sqlconnect"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = sync.Mutex{}
	nextID   = 1
)

func getTeacherHandler(id int, w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	log.Println(err)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying data", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status,omitempty"`
		Data   models.Teacher `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   teacher,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeacher)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare(`INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)`)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error inserting into database.", http.StatusInternalServerError)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error getting last id", http.StatusInternalServerError)
		return
	}

	newTeacher.ID = int(lastID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string         `json:"status,omitempty"`
		Data   models.Teacher `json:"data,omitempty"`
	}{
		Status: "success",
		Data:   newTeacher,
	}
	json.NewEncoder(w).Encode(response)
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []any

	if firstName != "" {
		query += " AND first_name=?"
		args = append(args, firstName)
	}
	if lastName != "" {
		query += " AND last_name=?"
		args = append(args, lastName)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error querying data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error scanning data", http.StatusInternalServerError)
			return
		}
		teacherList = append(teacherList, teacher)
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
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
