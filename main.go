package main

import (
	"fmt"
	"log"
	"os"

	"github.com/khabirovar/crud/backend"
)

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, name, user, pass)
	backend, err := backend.NewBackend(dsn, "8080")
	if err != nil {
		log.Fatal(err)
	}
	backend.Run()
}
