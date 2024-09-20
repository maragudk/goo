package sql_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"maragu.dev/is"

	"maragu.dev/goo/sql"
	"maragu.dev/goo/sqltest"
)

func TestHelper_CreateJobInTx(t *testing.T) {
	t.Run("creates a job in goqite", func(t *testing.T) {
		sqlHelper := sqltest.NewHelper(t)

		err := sqlHelper.InTransaction(context.Background(), func(tx *sql.Tx) error {
			return sqlHelper.CreateJobInTx(context.Background(), tx, "test", map[string]string{"key": "value"})
		})
		is.NotError(t, err)

		m, err := sqlHelper.JobsQ.Receive(context.Background())
		is.NotError(t, err)

		m2 := decodeGob(m.Body)
		is.Equal(t, "test", m2.Name)
		is.Equal(t, `{"key":"value"}`, string(m2.Message))
	})
}

type message struct {
	Name    string
	Message []byte
}

func decodeGob(body []byte) message {
	var m message
	if err := gob.NewDecoder(bytes.NewReader(body)).Decode(&m); err != nil {
		panic(err)
	}
	return m
}
