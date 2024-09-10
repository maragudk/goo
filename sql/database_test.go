package sql_test

import (
	"context"
	"testing"

	"maragu.dev/is"

	"maragu.dev/goo/sqltest"
)

func TestDatabase_Migrate(t *testing.T) {
	t.Run("can migrate down and back up", func(t *testing.T) {
		db := sqltest.CreateDatabase(t)

		err := db.MigrateDown(context.Background())
		is.NotError(t, err)

		err = db.MigrateUp(context.Background())
		is.NotError(t, err)
	})
}
