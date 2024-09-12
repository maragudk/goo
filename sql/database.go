package sql

import (
	"context"
	"database/sql"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maragudk/goqite"
	_ "github.com/mattn/go-sqlite3"
	"maragu.dev/errors"
	"maragu.dev/migrate"
	"maragu.dev/snorkel"
)

type Database struct {
	DB        *sqlx.DB
	jobsQueue *goqite.Queue
	log       *snorkel.Logger
	path      string
}

type NewDatabaseOptions struct {
	Log  *snorkel.Logger
	Path string
}

// NewDatabase with the given options.
// If no logger is provided, logs are discarded.
func NewDatabase(opts NewDatabaseOptions) *Database {
	if opts.Log == nil {
		opts.Log = snorkel.New(snorkel.Options{W: io.Discard})
	}

	// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	opts.Path += "?_journal=WAL&_timeout=5000&_fk=true"

	return &Database{
		log:  opts.Log,
		path: opts.Path,
	}
}

func (d *Database) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.log.Event("Starting database", 1, "path", d.path)

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, "sqlite3", d.path)
	if err != nil {
		return err
	}

	return nil
}

// InTransaction runs callback in a transaction, and makes sure to handle rollbacks, commits etc.
func (d *Database) InTransaction(ctx context.Context, callback func(tx *sqlx.Tx) error) (err error) {
	tx, err := d.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
	}
	defer func() {
		if rec := recover(); rec != nil {
			err = rollback(tx, errors.Newf("panic: %v", rec))
		}
	}()
	if err := callback(tx); err != nil {
		return rollback(tx, err)
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing transaction")
	}

	return nil
}

// rollback a transaction, handling both the original error and any transaction rollback errors.
func rollback(tx *sqlx.Tx, err error) error {
	if txErr := tx.Rollback(); txErr != nil {
		return errors.Wrap(err, "error rolling back transaction after error (transaction error: %v), original error", txErr)
	}
	return err
}

func (d *Database) MigrateUp(ctx context.Context) error {
	// TODO some migrations should be in goo
	return migrate.Up(ctx, d.DB.DB, d.getMigrations())
}

func (d *Database) MigrateDown(ctx context.Context) error {
	return migrate.Down(ctx, d.DB.DB, d.getMigrations())
}

func (d *Database) getMigrations() fs.FS {
	for _, path := range []string{"sql/migrations", "../sql/migrations"} {
		migrations := os.DirFS(path)

		matches, err := fs.Glob(migrations, "*.sql")
		if err == nil && len(matches) > 0 {
			return migrations
		}
	}

	panic("no migrations found")
}

func (d *Database) Ping(ctx context.Context) error {
	return d.InTransaction(ctx, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `select 1`)
		return err
	})
}

func (d *Database) SetJobsQueue(q *goqite.Queue) {
	d.jobsQueue = q
}
