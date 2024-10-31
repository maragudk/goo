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

	"maragu.dev/goo/email"
	"maragu.dev/goo/html"
	"maragu.dev/goo/http"
	"maragu.dev/goo/jobs"
	"maragu.dev/goo/llm"
	"maragu.dev/goo/sql"
)

type Options struct {
	HTMLPage            html.PageFunc
	HTTPRouterInjector  func(*http.Router)
	LLMPrompterInjector func(llm.Prompter)
	Log                 *snorkel.Logger
	Migrate             bool
	SQLHelperInjector   func(*sql.Helper)
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

	baseURL := env.GetStringOrDefault("BASE_URL", "http://localhost:8080")

	sender := email.NewSender(email.NewSenderOptions{
		BaseURL:                   baseURL,
		Log:                       log,
		MarketingEmailAddress:     env.GetStringOrDefault("MARKETING_EMAIL_ADDRESS", "marketing@example.com"),
		MarketingEmailName:        env.GetStringOrDefault("MARKETING_EMAIL_NAME", "Marketing"),
		ReplyToEmailAddress:       env.GetStringOrDefault("REPLY_TO_EMAIL_ADDRESS", "support@example.com"),
		ReplyToEmailName:          env.GetStringOrDefault("REPLY_TO_EMAIL_NAME", "Support"),
		Token:                     env.GetStringOrDefault("POSTMARK_TOKEN", ""),
		TransactionalEmailAddress: env.GetStringOrDefault("TRANSACTIONAL_EMAIL_ADDRESS", "transactional@example.com"),
		TransactionalEmailName:    env.GetStringOrDefault("TRANSACTIONAL_EMAIL_NAME", "Transactional"),
	})

	llmClient := llm.NewOpenAIClient(llm.NewOpenAIClientOptions{
		BaseURL: env.GetStringOrDefault("LLM_URL", "https://api.fireworks.ai/inference/v1"),
		Log:     log,
		Model:   llm.Model(env.GetStringOrDefault("LLM_MODEL", llm.ModelFireworksLlama3_1_8B.String())),
		Token:   env.GetStringOrDefault("LLM_TOKEN", ""),
	})

	if opts.LLMPrompterInjector != nil {
		opts.LLMPrompterInjector(llmClient)
	}

	jobs.Register(r, jobs.RegisterOpts{
		Log:    log,
		Sender: sender,
	})

	sqlHelper.JobsQ = q

	s := http.NewServer(http.NewServerOptions{
		Address:            env.GetStringOrDefault("SERVER_ADDRESS", ":8080"),
		AdminPassword:      env.GetStringOrDefault("ADMIN_PASSWORD", "correct horse battery staple"),
		BaseURL:            baseURL,
		HTTPRouterInjector: opts.HTTPRouterInjector,
		HTMLPage:           opts.HTMLPage,
		LLMClient:          llmClient,
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

	sqlHelper := sql.NewHelper(sql.NewHelperOptions{
		Log:  log,
		Path: env.GetStringOrDefault("DATABASE_PATH", "app.db"),
	})
	if err := sqlHelper.Connect(); err != nil {
		return err
	}

	if err := sqlHelper.MigrateUp(context.Background()); err != nil {
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
