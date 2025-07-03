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
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = sync.Mutex{}
	nextID   = 1
)

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
				log.Println("Invalid ID")
				http.Error(w, "Invalid ID", http.StatusBadRequest)
				return
			}
			getTeacherHandler(id, w, r)
		}

		//return
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPatch:
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		idStr := strings.TrimSuffix(path, "/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Println("Invalid ID")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		updateTeacherHandler(id, w, r)
	case http.MethodDelete:
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		idStr := strings.TrimSuffix(path, "/")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			log.Println("Invalid ID")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		deleteTeacherHandler(id, w, r)

	}

	//fmt.Println("Hello Teachers Route")
	//w.Write([]byte("Hello Teachers Route"))

}

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

	//firstName := r.URL.Query().Get("first_name")
	//lastName := r.URL.Query().Get("last_name")

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []any

	query, args = addFilters(r, query, args)
	query = addSorting(r, query)

	fmt.Println("Query: ", query)

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

func updateTeacherHandler(id int, w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return
	}

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer db.Close()

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

func deleteTeacherHandler(id int, w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error connecting database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

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
		ID     int    `json:"id"`
	}{
		Status: "success",
		ID:     id,
	}
	json.NewEncoder(w).Encode(response)
}

func addFilters(r *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += fmt.Sprintf(" AND %s = ?", dbField)
			args = append(args, value)
		}
	}

	return query, args
}

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sort_by"]
	//["first_name:asc", "last_name:desc"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortOrder(order) || !isValidSortField(field) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf(" %s %s", field, strings.ToUpper(order))
		}
	}

	return query
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	return validFields[field]

}
