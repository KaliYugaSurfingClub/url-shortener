package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	dbURL, ok := os.LookupEnv("POSTGRES_URL")
	if !ok {
		log.Fatal("POSTGRES_URL environment variable not set")
	}

	migratePath := "file://migrations"

	m, err := migrate.New(migratePath, dbURL)
	if err != nil {
		log.Fatal("migrate.New: ", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Up: ", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal("get version: ", err)
	}

	fmt.Printf("Applied migration %d, Dirty: %t\n", version, dirty)
}
