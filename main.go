package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/vladwithcode/qrcatalog/internal/auth"
	"github.com/vladwithcode/qrcatalog/internal/db"
	"github.com/vladwithcode/qrcatalog/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to set enviroment from file\n%v\n", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("required env var PORT is not set")
	}

	dbPool, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to DB:\n%v\n", err)
	}
	defer dbPool.Close()

	auth.SetAuthParameters()

	router := routes.NewRouter()
	fmt.Printf("Starting server on port http://localhost:%s\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatalf("failed to listen and serve %v\n", err)
	}
}
