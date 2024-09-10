package sqltest

import (
	"context"
	"testing"

	"github.com/maragudk/goqite"

	"maragu.dev/goo/sql"
)

// CreateDatabase for testing.
func CreateDatabase(t *testing.T) *sql.Database {
	t.Helper()

	db := sql.NewDatabase(sql.NewDatabaseOptions{
		Path: ":memory:",
	})
	if err := db.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := db.MigrateUp(context.Background()); err != nil {
		t.Fatal(err)
	}

	q := goqite.New(goqite.NewOpts{
		DB:   db.DB.DB,
		Name: "jobs",
	})
	db.SetJobsQueue(q)

	return db
}
