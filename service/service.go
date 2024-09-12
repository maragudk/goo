package service

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/maragudk/goqite"
	qjobs "github.com/maragudk/goqite/jobs"
	"golang.org/x/sync/errgroup"
	"maragu.dev/env"
	"maragu.dev/snorkel"

	"maragu.dev/goo/http"
	"maragu.dev/goo/sql"
)

type Options struct {
	HTTPRouterInjector func(chi.Router)
	Log                *snorkel.Logger
	Migrate            bool
	SQLHelperInjector  func(*sql.Helper)
}

func Start(opts Options) {
	log := opts.Log

	if opts.Migrate {
		if err := migrate(opts); err != nil {
			log.Event("Error migrating", 1, "error", err)
			return
		}
	}

	if err := start(opts); err != nil {
		log.Event("Error starting", 1, "error", err)
	}
}

func start(opts Options) error {
	log := opts.Log

	log.Event("Starting app", 1)

	_ = env.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	sqlHelper := sql.NewHelper(sql.NewHelperOptions{
		Log:  log,
		Path: env.GetStringOrDefault("DATABASE_PATH", "app.db"),
	})
	if err := sqlHelper.Connect(); err != nil {
		return err
	}
	if opts.SQLHelperInjector != nil {
		opts.SQLHelperInjector(sqlHelper)
	}

	q := goqite.New(goqite.NewOpts{
		DB:   sqlHelper.DB.DB,
		Name: "jobs",
	})

	r := qjobs.NewRunner(qjobs.NewRunnerOpts{
		Log:   &logAdapter{log: log},
		Queue: q,
	})

	sqlHelper.SetJobsQueue(q)

	s := http.NewServer(http.NewServerOptions{
		AdminPassword:      env.GetStringOrDefault("ADMIN_PASSWORD", "correct horse battery staple"),
		BaseURL:            env.GetStringOrDefault("BASE_URL", "http://localhost:8080"),
		HTTPRouterInjector: opts.HTTPRouterInjector,
		Log:                log,
		SecureCookie:       env.GetBoolOrDefault("SECURE_COOKIE", true),
		SQLHelper:          sqlHelper,
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

func migrate(opts Options) error {
	log := opts.Log

	log.Event("Migrating", 1)

	_ = env.Load()

	db := sql.NewHelper(sql.NewHelperOptions{
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

func NewLogger() *snorkel.Logger {
	return snorkel.New(snorkel.Options{})
}
