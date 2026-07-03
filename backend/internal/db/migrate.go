package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func SortMigrationFiles(files []string) []string {
	out := append([]string(nil), files...)
	sort.Strings(out)
	return out
}

func RunMigrations(conn *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	for _, name := range SortMigrationFiles(files) {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if _, err := conn.Exec(string(body)); err != nil {
			return err
		}
	}
	return nil
}
