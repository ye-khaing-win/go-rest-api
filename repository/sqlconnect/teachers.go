package sqlconnect

import (
	"database/sql"
	"errors"
	"fmt"
	"restapi/internal/models"
	"strings"
)

func GetTeachers(filters map[string]any, sortParams []string) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []any

	for k, v := range filters {
		query += fmt.Sprintf(" AND %s = ?", k)
		args = append(args, v)
	}

	query = sort(sortParams, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teachers := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func GetTeacherByID(id string) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, err
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id=?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Teacher{}, err
	} else if err != nil {
		return models.Teacher{}, err
	}

	return teacher, nil

}

func CreateTeacher(newTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)`)
	if err != nil {
		return models.Teacher{}, nil
	}
	defer stmt.Close()

	res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
	if err != nil {
		return models.Teacher{}, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return models.Teacher{}, err
	}

	newTeacher.ID = int(lastID)

	return newTeacher, nil

}

func sort(sortParams []string, query string) string {
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
