package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

type ExportData struct {
	Sections []*db.Section `json:"sections"`
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

	file, err := os.Create(flags.File)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	ctx := context.Background()
	var data ExportData

	sections, err := db.FindAllSections(ctx)
	if err != nil {
		fmt.Printf("failed to find all sections: %v\n", err)
		os.Exit(1)
	}

	data.Sections = sections
	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		fmt.Printf("failed to encode json: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Sections exported successfully")
}

type Flags struct {
	File string `json:"file"`
}

func parseFlags() Flags {
	var flags Flags

	flag.StringVar(&flags.File, "f", "", "File to export")
	flag.StringVar(&flags.File, "file", "", "File to export")
	flag.Parse()

	if flags.File == "" {
		panic("specify either -f or -file")
	}

	return flags
}
