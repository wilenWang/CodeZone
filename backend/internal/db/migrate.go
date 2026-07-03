package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type MigrationOptions struct {
	IncludeDevSeed bool
}

func SortMigrationFiles(files []string) []string {
	out := append([]string(nil), files...)
	sort.Strings(out)
	return out
}

func RunMigrations(conn *sql.DB, dir string) error {
	return RunMigrationsWithOptions(conn, dir, MigrationOptions{IncludeDevSeed: true})
}

func RunMigrationsWithOptions(conn *sql.DB, dir string, opts MigrationOptions) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir %s: %w", dir, err)
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	for _, name := range FilterMigrationFiles(files, opts.IncludeDevSeed) {
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		for i, stmt := range SplitSQLStatements(string(body)) {
			if _, err := conn.Exec(stmt); err != nil {
				return fmt.Errorf("run migration %s statement %d: %w", name, i+1, err)
			}
		}
	}
	return nil
}

func FilterMigrationFiles(files []string, includeDevSeed bool) []string {
	var out []string
	for _, name := range SortMigrationFiles(files) {
		if !includeDevSeed && strings.Contains(name, "_seed_dev") {
			continue
		}
		out = append(out, name)
	}
	return out
}

func SplitSQLStatements(body string) []string {
	parts := strings.Split(body, ";")
	statements := make([]string, 0, len(parts))
	for _, part := range parts {
		stmt := strings.TrimSpace(part)
		if stmt == "" {
			continue
		}
		statements = append(statements, stmt)
	}
	return statements
}
