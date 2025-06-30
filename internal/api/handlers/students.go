package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
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
