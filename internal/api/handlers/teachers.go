package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"reflect"
	"restapi/internal/models"
	"restapi/repository/sqlconnect"
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
	}
	filters := make(map[string]any)

	for param, field := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			filters[field] = value
		}
	}

	sortParams := r.URL.Query()["sort_by"]
	teachers, err := sqlconnect.GetTeachers(filters, sortParams)

	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	teacher, err := sqlconnect.GetTeacherByID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "success",
		Data:   teacher,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func CreateTeacherHandler(w http.ResponseWriter, r *http.Request) {

	var newTeacher models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeacher)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newTeacher, err = sqlconnect.CreateTeacher(newTeacher)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

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
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		return
	}
	defer db.Close()
	id := r.PathValue("id")

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		http.Error(w, "Teacher not found.", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Error querying data", http.StatusInternalServerError)
		return
	}

	//for k, v := range updates {
	//	switch k {
	//	case "first_name":
	//		existingTeacher.FirstName = v.(string)
	//	case "last_name":
	//		existingTeacher.LastName = v.(string)
	//	case "email":
	//		existingTeacher.Email = v.(string)
	//	case "class":
	//		existingTeacher.Class = v.(string)
	//	case "subject":
	//		existingTeacher.Subject = v.(string)
	//	}
	//}

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == fmt.Sprintf("%s,omitempty", k) {
				if teacherVal.Field(i).CanSet() {
					fieldValue := teacherVal.Field(i)
					fieldValue.Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)

	if err != nil {
		http.Error(w, "Error updating data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)

}
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	id := r.PathValue("id")

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erro deleting teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving delete result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	}{
		Status: "success",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}
