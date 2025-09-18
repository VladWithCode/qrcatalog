package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

type ImportData struct {
	Categories []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Slug            string `json:"slug"`
		Description     string `json:"description"`
		LongDescription string `json:"long_description"`
	} `json:"categories"`
	Subcategories []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Slug            string `json:"slug"`
		Description     string `json:"description"`
		LongDescription string `json:"long_description"`
		CategorySlug    string `json:"category_slug"`
	} `json:"subcategories"`
	Products []struct {
		ID              string  `json:"id"`
		Name            string  `json:"name"`
		Slug            string  `json:"slug"`
		Description     string  `json:"description"`
		LongDescription string  `json:"long_description"`
		Price           float64 `json:"price"`
		Unit            string  `json:"unit"`
		Quantity        int     `json:"quantity"`
		CategorySlug    string  `json:"category_slug"`
		SubcategorySlug string  `json:"subcategory_slug"`
		Available       bool    `json:"available"`
	} `json:"products"`
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

	ctgSlugToID := make(map[string]string)
	ctgBatch := pgx.Batch{}
	for _, category := range data.Categories {
		id := uuid.Must(uuid.NewV7()).String()

		ctgArgs := pgx.NamedArgs{
			"id":               id,
			"name":             category.Name,
			"slug":             category.Slug,
			"description":      category.Description,
			"long_description": category.LongDescription,
		}
		ctgBatch.Queue(
			`INSERT INTO categories
				(id, name, slug, description, long_description)
			VALUES (@id, @name, @slug, @description, @long_description)
			`,
			ctgArgs,
		)

		slug := ctgArgs["slug"].(string)
		ctgSlugToID[slug] = id
	}
	ctgResults := tx.SendBatch(ctx, &ctgBatch)

	_, err = ctgResults.Exec()
	ctgResults.Close()
	if err != nil {
		fmt.Printf("failed to create categories: %v\n", err)
		return
	}

	sctgSlugToID := make(map[string]string)
	sctgBatch := pgx.Batch{}
	for _, subcategory := range data.Subcategories {
		id := uuid.Must(uuid.NewV7()).String()
		sctg := pgx.NamedArgs{
			"id":               id,
			"name":             subcategory.Name,
			"slug":             subcategory.Slug,
			"description":      subcategory.Description,
			"long_description": subcategory.LongDescription,
		}
		if subcategory.CategorySlug != "" {
			sctg["category_id"] = ctgSlugToID[subcategory.CategorySlug]
		}

		sctgBatch.Queue(
			`INSERT INTO subcategories
				(id, name, slug, description, long_description, category_id)
			VALUES (@id, @name, @slug, @description, @long_description, @category_id)
			`,
			sctg,
		)

		sctgSlugToID[sctg["slug"].(string)] = id
	}
	sctgResults := tx.SendBatch(ctx, &sctgBatch)
	_, err = sctgResults.Exec()
	sctgResults.Close()
	if err != nil {
		fmt.Printf("failed to create subcategories: %v\n", err)
		return
	}

	prodBatch := pgx.Batch{}
	for i, product := range data.Products {
		id := uuid.Must(uuid.NewV7()).String()
		prod := pgx.NamedArgs{
			"id":               id,
			"name":             product.Name,
			"slug":             product.Slug,
			"description":      product.Description,
			"long_description": product.LongDescription,
			"price":            product.Price,
			"unit":             product.Unit,
			"quantity":         product.Quantity,
			"available":        product.Available,
		}
		if product.CategorySlug != "" {
			prod["category_id"] = ctgSlugToID[product.CategorySlug]
		}
		if product.SubcategorySlug != "" {
			prod["subcategory_id"] = sctgSlugToID[product.SubcategorySlug]
		}
		prodBatch.Queue(
			`INSERT INTO products
				(id, name, slug, description, long_description, price, unit, quantity, available, category_id, subcategory_id)
			VALUES (@id, @name, @slug, @description, @long_description, @price, @unit, @quantity, @available, @category_id, @subcategory_id)
			`,
			prod,
		)
		data.Products[i].ID = id
	}
	prodResults := tx.SendBatch(ctx, &prodBatch)
	_, err = prodResults.Exec()
	prodResults.Close()
	if err != nil {
		fmt.Printf("failed to create products: %v\n", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		fmt.Printf("failed to commit transaction: %v\n", err)
		return
	}
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
