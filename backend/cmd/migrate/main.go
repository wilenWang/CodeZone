package main

import (
	"log"
	"path/filepath"

	"codezone/backend/internal/config"
	"codezone/backend/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
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
