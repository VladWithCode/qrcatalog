package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

type ImportData struct {
	Sections []db.Section `json:"sections"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("failed to load .env file: %v\n", err)
		return
	}

	flags := parseFlags()

	conn, err := db.Connect()
	if err != nil {
		fmt.Printf("failed to connect to db: %v\n", err)
		return
	}
	defer conn.Close()

	file, err := os.Open(flags.File)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	var data ImportData
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		fmt.Printf("failed to decode json: %v\n", err)
		return
	}
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		fmt.Printf("failed to begin transaction: %v\n", err)
		return
	}
	defer tx.Rollback(ctx)

	for _, section := range data.Sections {
		id := uuid.Must(uuid.NewV7()).String()
		section.ID = id
		err = db.CreateSection(ctx, &section)
		if err != nil {
			fmt.Printf("failed to create section \"%s\": %v\n", section.Name, err)
			return
		}
	}

	fmt.Println("Sections imported successfully")
}

type Flags struct {
	File string `json:"file"`
}

func parseFlags() Flags {
	var flags Flags

	flag.StringVar(&flags.File, "f", "", "File to import")
	flag.StringVar(&flags.File, "file", "", "File to import")
	flag.Parse()

	if flags.File == "" {
		panic("specify either -f or -file")
	}

	return flags
}
