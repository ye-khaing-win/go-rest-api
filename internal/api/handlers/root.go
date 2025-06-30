package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello Root Route")
	w.Write([]byte("Hello Root Route"))
}
