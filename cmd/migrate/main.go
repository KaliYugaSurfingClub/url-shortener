package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	fmt.Println("start migrate")

	dbURL := "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"
	migratePath := "file://C:/Users/leono/Desktop/prog/go/shortener/tools/migrate"

	m, err := migrate.New(migratePath, dbURL)
	if err != nil {
		log.Fatal("migrate.New ", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Up ", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal("get version ", err)
	}

	fmt.Printf("Applied migration %d, Dirty: %t\n", version, dirty)
}
