package main

import (
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

func main() {
	err := godotenv.Overload()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
	var (
		user = os.Getenv("ADMIN_USERNAME")
		pass = os.Getenv("ADMIN_PASSWORD")
	)
	if user == "" || pass == "" {
		log.Fatalf("missing ADMIN_USERNAME or ADMIN_PASSWORD")
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	id := uuid.Must(uuid.NewV7()).String()
	defaultUser := db.User{
		ID:       id,
		Fullname: "Administrador",
		Username: user,
		Password: pass,
		Role:     db.RoleAdmin,
		Email:    "admin@qrestrellas.com",
	}
	_, err = db.CreateUser(&defaultUser)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Printf("User created with id: %v", id)
}
