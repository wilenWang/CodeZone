package main

import (
	"log"
	"path/filepath"

	"codezone/backend/internal/config"
	"codezone/backend/internal/db"
)

func main() {
	cfg := config.Load()
	conn, err := db.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err := db.RunMigrationsWithOptions(conn, filepath.Join("migrations"), db.MigrationOptions{IncludeDevSeed: cfg.DevSeed}); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied")
}
