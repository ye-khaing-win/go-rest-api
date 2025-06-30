package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Execs Route")
	w.Write([]byte("Hello Execs Route"))
}
