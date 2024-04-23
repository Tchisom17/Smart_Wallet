package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tchisom17/internal/app/handler/accounthand"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	http.HandleFunc("/create-account", accounthand.HandleCreateAccount)
	fmt.Println("started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
