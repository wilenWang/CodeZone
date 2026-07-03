package db

import "testing"

func TestMigrationFilesAreSorted(t *testing.T) {
	files := []string{"0002_seed_dev.sql", "0001_schema.sql"}
	got := SortMigrationFiles(files)
	want := []string{"0001_schema.sql", "0002_seed_dev.sql"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}
