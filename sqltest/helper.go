package sqltest

import (
	"context"
	"testing"

	"github.com/maragudk/goqite"

	"maragu.dev/goo/sql"
)

// NewHelper for testing.
func NewHelper(t *testing.T) *sql.Helper {
	t.Helper()

	sqlHelper := sql.NewHelper(sql.NewHelperOptions{
		Path: ":memory:",
	})
	if err := sqlHelper.Connect(); err != nil {
		t.Fatal(err)
	}

	if err := sqlHelper.MigrateUp(context.Background()); err != nil {
		t.Fatal(err)
	}

	q := goqite.New(goqite.NewOpts{
		DB:   sqlHelper.DB.DB,
		Name: "jobs",
	})
	sqlHelper.SetJobsQueue(q)

	return sqlHelper
}
