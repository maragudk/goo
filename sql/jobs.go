package sql

import (
	"context"
	"encoding/json"

	"github.com/maragudk/goqite/jobs"
)

type stringMap = map[string]string

func (h *Helper) CreateJobInTx(ctx context.Context, tx *Tx, name string, m stringMap) error {
	return jobs.CreateTx(ctx, tx.Tx.Tx, h.JobsQ, name, mustMarshalJSON(m))
}

func mustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
