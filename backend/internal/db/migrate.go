package db

import (
	"context"
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

type statementExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
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
	ctx := context.Background()
	pinnedConn, err := conn.Conn(ctx)
	if err != nil {
		return fmt.Errorf("connect for migrations: %w", err)
	}
	defer pinnedConn.Close()

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
		if err := runMigrationStatements(ctx, pinnedConn, name, SplitSQLStatements(string(body))); err != nil {
			return err
		}
	}
	return nil
}

func runMigrationStatements(ctx context.Context, exec statementExecutor, name string, statements []string) error {
	for i, stmt := range statements {
		if _, err := exec.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("run migration %s statement %d: %w", name, i+1, err)
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
