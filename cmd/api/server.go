package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	router2 "restapi/internal/api/router"
	"restapi/repository/sqlconnect"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Error...", err)
	}
	port := os.Getenv("HTTP_PORT")

	router := router2.Router()

	//rl := mw.NewRateLimiter(5, 10*time.Second)
	//hpp := mw.HPP{
	//	CheckQuery:      true,
	//	CheckBody:       true,
	//	BodyContentType: "application/x-www-form-urlencoded",
	//	Whitelist:       []string{"name", "age", "gender"},
	//}

	secureMux := mw.SecurityHeaders(router)
	//secureMux := rl.Middleware(mw.ResponseTime(mw.SecurityHeaders(mw.Compression(hpp.Middleware()(mux)))))

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: secureMux,
	}
	fmt.Println("Server running on port: ", port)

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting the server", err)
	}
}
