package sql_test

import (
	"context"
	"testing"

	"maragu.dev/is"

	"maragu.dev/goo/sqltest"
)

func TestHelper_Migrate(t *testing.T) {
	t.Run("can migrate down and back up", func(t *testing.T) {
		sqlHelper := sqltest.NewHelper(t)

		err := sqlHelper.MigrateDown(context.Background())
		is.NotError(t, err)

		err = sqlHelper.MigrateUp(context.Background())
		is.NotError(t, err)

		var version string
		err = sqlHelper.Get(context.Background(), &version, `select version from migrations`)
		is.NotError(t, err)
		is.True(t, len(version) > 0)
	})
}
