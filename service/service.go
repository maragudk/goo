package service

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/maragudk/goqite"
	qjobs "github.com/maragudk/goqite/jobs"
	"golang.org/x/sync/errgroup"
	"maragu.dev/env"
	"maragu.dev/snorkel"

	"maragu.dev/goo/http"
	"maragu.dev/goo/sql"
)

type Options struct {
	Migrate bool
}

func Start(opts Options) {
	log := snorkel.New(snorkel.Options{})

	if opts.Migrate {
		if err := migrate(log); err != nil {
			log.Event("Error migrating", 1, "error", err)
		}
	}

	if err := start(log); err != nil {
		log.Event("Error starting", 1, "error", err)
	}
}

func start(log *snorkel.Logger) error {
	log.Event("Starting app", 1)

	_ = env.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	db := sql.NewDatabase(sql.NewDatabaseOptions{
		Log:  log,
		Path: env.GetStringOrDefault("DATABASE_PATH", "app.db"),
	})
	if err := db.Connect(); err != nil {
		return err
	}

	q := goqite.New(goqite.NewOpts{
		DB:   db.DB.DB,
		Name: "jobs",
	})

	r := qjobs.NewRunner(qjobs.NewRunnerOpts{
		Log:   &logAdapter{log: log},
		Queue: q,
	})

	db.SetJobsQueue(q)

	s := http.NewServer(http.NewServerOptions{
		AdminPassword: env.GetStringOrDefault("ADMIN_PASSWORD", "correct horse battery staple"),
		BaseURL:       env.GetStringOrDefault("BASE_URL", "http://localhost:8080"),
		DB:            db,
		Log:           log,
		SecureCookie:  env.GetBoolOrDefault("SECURE_COOKIE", true),
	})

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.Start()
	})

	eg.Go(func() error {
		r.Start(ctx)
		return nil
	})

	<-ctx.Done()
	log.Event("Stopping app", 1)

	eg.Go(func() error {
		return s.Stop()
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	log.Event("Stopped app", 1)

	return nil
}

func migrate(log *snorkel.Logger) error {
	log.Event("Migrating", 1)

	_ = env.Load()

	db := sql.NewDatabase(sql.NewDatabaseOptions{
		Log:  log,
		Path: env.GetStringOrDefault("DATABASE_PATH", "app.db"),
	})
	if err := db.Connect(); err != nil {
		return err
	}

	if err := db.MigrateUp(context.Background()); err != nil {
		return err
	}

	log.Event("Migrated", 1)

	return nil
}

type logAdapter struct {
	log *snorkel.Logger
}

func (l *logAdapter) Info(msg string, args ...any) {
	l.log.Event(msg, 1, args...)
}
