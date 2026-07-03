package main

import (
	"log"
	"path/filepath"

	"vibework-chat/backend/internal/config"
	"vibework-chat/backend/internal/db"
)

func main() {
	cfg := config.Load()
	conn, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := db.RunMigrations(conn, filepath.Join("migrations")); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}
