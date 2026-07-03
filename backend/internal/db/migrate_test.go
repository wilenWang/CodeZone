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

func TestSplitSQLStatements(t *testing.T) {
	body := "  CREATE TABLE a (id BIGINT);\n\nINSERT INTO a VALUES (1);\n ; \n"
	got := SplitSQLStatements(body)
	want := []string{
		"CREATE TABLE a (id BIGINT)",
		"INSERT INTO a VALUES (1)",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d statements want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d: got %q want %q", i, got[i], want[i])
		}
	}
}

func TestFilterMigrationFilesHonorsDevSeed(t *testing.T) {
	files := []string{"0002_seed_dev.sql", "0003_feature.sql", "0001_schema.sql"}

	withSeed := FilterMigrationFiles(files, true)
	wantWithSeed := []string{"0001_schema.sql", "0002_seed_dev.sql", "0003_feature.sql"}
	for i := range wantWithSeed {
		if withSeed[i] != wantWithSeed[i] {
			t.Fatalf("with seed index %d: got %q want %q", i, withSeed[i], wantWithSeed[i])
		}
	}

	withoutSeed := FilterMigrationFiles(files, false)
	wantWithoutSeed := []string{"0001_schema.sql", "0003_feature.sql"}
	if len(withoutSeed) != len(wantWithoutSeed) {
		t.Fatalf("without seed got %d files want %d: %#v", len(withoutSeed), len(wantWithoutSeed), withoutSeed)
	}
	for i := range wantWithoutSeed {
		if withoutSeed[i] != wantWithoutSeed[i] {
			t.Fatalf("without seed index %d: got %q want %q", i, withoutSeed[i], wantWithoutSeed[i])
		}
	}
}
