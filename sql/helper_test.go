package sql_test

import (
	"context"
	"testing"

	"maragu.dev/is"

	"maragu.dev/goo/sqltest"
)

func TestHelper_Ping(t *testing.T) {
	t.Run("can ping", func(t *testing.T) {
		db := sqltest.NewHelper(t)

		err := db.Ping(context.Background())
		is.NotError(t, err)
	})
}
