package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maragudk/goqite/jobs"
	"maragu.dev/snorkel"

	"maragu.dev/goo/model"
)

type emailSender interface {
	SendTransactionalEmail(ctx context.Context, name string, email model.Email, subject, preheader, template string, kw map[string]string) error
}

func sendEmail(log *snorkel.Logger, sender emailSender) jobs.Func {
	return func(ctx context.Context, m2 []byte) error {
		m := mustUnmarshalJSON(m2)

		typ := m["type"]
		email := m["email"]

		log.Event("Sending email", 1, "type", typ, "email", email)

		var err error
		switch typ {
		case "signup":
			err = sendSignupEmail(sender, ctx, m)
		case "login":
			err = sendLoginEmail(sender, ctx, m)
		default:
			panic("unknown email type " + typ)
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func sendSignupEmail(sender emailSender, ctx context.Context, m map[string]string) error {
	subject := fmt.Sprintf("Welcome, %v!", m["name"])
	return sender.SendTransactionalEmail(ctx, m["name"], model.Email(m["email"]),
		subject, "Click the link to complete sign up.", "signup", m)
}

func sendLoginEmail(sender emailSender, ctx context.Context, m map[string]string) error {
	subject := fmt.Sprintf("Welcome back, %v!", m["name"])
	return sender.SendTransactionalEmail(ctx, m["name"], model.Email(m["email"]),
		subject, "Click the link to log in.", "login", m)
}

func mustUnmarshalJSON(m []byte) map[string]string {
	result := map[string]string{}
	err := json.Unmarshal(m, &result)
	if err != nil {
		panic(err)
	}
	return result
}
