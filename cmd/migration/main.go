package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/yaBliznyk/newsportal/migrations"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	dbURL := getEnv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	goose.SetBaseFS(migrations.Files)

	switch command {
	case "up":
		if err := goose.Up(db, "."); err != nil {
			log.Fatalf("goose up: %v", err)
		}
		fmt.Println("Migrations applied successfully")

	case "down":
		if err := goose.Down(db, "."); err != nil {
			log.Fatalf("goose down: %v", err)
		}
		fmt.Println("Migration rolled back successfully")

	case "reset":
		if err := goose.Reset(db, "."); err != nil {
			log.Fatalf("goose reset: %v", err)
		}
		fmt.Println("All migrations rolled back successfully")

	case "status":
		if err := goose.Status(db, "."); err != nil {
			log.Fatalf("goose status: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: migration <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up      Apply all pending migrations")
	fmt.Println("  down    Roll back the last migration")
	fmt.Println("  reset   Roll back all migrations")
	fmt.Println("  status  Show migration status")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
