package jobs

import (
	"github.com/maragudk/goqite/jobs"
	"maragu.dev/snorkel"

	"maragu.dev/goo/email"
)

type RegisterOpts struct {
	Log    *snorkel.Logger
	Sender *email.Sender
}

func Register(r *jobs.Runner, opts RegisterOpts) {
	r.Register("send-email", sendEmail(opts.Log, opts.Sender))
}
