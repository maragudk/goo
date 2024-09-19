package email

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"maragu.dev/errors"
	"maragu.dev/snorkel"

	"maragu.dev/goo/model"
)

const (
	marketingMessageStream     = "broadcast"
	transactionalMessageStream = "outbound"
)

type emailType int

const (
	marketing emailType = iota
	transactional
)

// nameAndEmail combo, of the form "Name <email@example.com>"
type nameAndEmail = string

type Keywords = map[string]string

// Sender can send transactional and marketing emails through Postmark.
// See https://postmarkapp.com/developer
type Sender struct {
	baseURL           string
	client            *http.Client
	endpointURL       string
	log               *snorkel.Logger
	marketingFrom     nameAndEmail
	replyTo           nameAndEmail
	token             string
	transactionalFrom nameAndEmail
}

type NewSenderOptions struct {
	BaseURL                   string
	EndpointURL               string
	Log                       *snorkel.Logger
	MarketingEmailAddress     string
	MarketingEmailName        string
	ReplyToEmailName          string
	ReplyToEmailAddress       string
	Token                     string
	TransactionalEmailAddress string
	TransactionalEmailName    string
}

func NewSender(opts NewSenderOptions) *Sender {
	if opts.Log == nil {
		opts.Log = snorkel.NewDiscard()
	}

	if opts.EndpointURL == "" {
		opts.EndpointURL = "https://api.postmarkapp.com/email"
	}

	return &Sender{
		baseURL:           strings.TrimSuffix(opts.BaseURL, "/"),
		client:            &http.Client{Timeout: 3 * time.Second},
		endpointURL:       strings.TrimSuffix(opts.EndpointURL, "/"),
		log:               opts.Log,
		marketingFrom:     createNameAndEmail(opts.MarketingEmailName, opts.MarketingEmailAddress),
		replyTo:           createNameAndEmail(opts.ReplyToEmailName, opts.ReplyToEmailAddress),
		token:             opts.Token,
		transactionalFrom: createNameAndEmail(opts.TransactionalEmailName, opts.TransactionalEmailAddress),
	}
}

func (s *Sender) SendTransactionalEmail(ctx context.Context, name string, email model.Email, subject, preheader, template string, kw Keywords) error {
	return s.send(ctx, transactional, createNameAndEmail(name, email.String()), subject, preheader, template, kw)
}

// requestBody used in Sender.send.
// See https://postmarkapp.com/developer/user-guide/send-email-with-api
type requestBody struct {
	MessageStream string
	From          nameAndEmail
	To            nameAndEmail
	ReplyTo       nameAndEmail
	Subject       string
	TextBody      string
	HtmlBody      string
}

func (s *Sender) send(ctx context.Context, typ emailType, to nameAndEmail, subject, preheader, template string, keywords Keywords) error {
	var messageStream string
	var from nameAndEmail
	switch typ {
	case marketing:
		from = s.marketingFrom
		messageStream = marketingMessageStream
	case transactional:
		from = s.transactionalFrom
		messageStream = transactionalMessageStream
	}

	// Keywords that are always included
	keywords["baseURL"] = s.baseURL

	err := s.sendRequest(ctx, requestBody{
		MessageStream: messageStream,
		From:          from,
		ReplyTo:       s.replyTo,
		To:            to,
		Subject:       subject,
		HtmlBody:      getEmail(template+".html", preheader, keywords),
	})

	return err
}

type postmarkResponse struct {
	ErrorCode int
	Message   string
}

// send using the Postmark API.
func (s *Sender) sendRequest(ctx context.Context, body requestBody) error {
	bodyAsBytes, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "error marshalling request body to json")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.endpointURL, bytes.NewReader(bodyAsBytes))
	if err != nil {
		return errors.Wrap(err, "error creating request")
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Postmark-Server-Token", s.token)

	response, err := s.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "error making request")
	}
	defer func() {
		_ = response.Body.Close()
	}()
	bodyAsBytes, err = io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}

	// https://postmarkapp.com/developer/api/overview#response-codes
	if response.StatusCode == http.StatusUnprocessableEntity {
		var r postmarkResponse
		if err := json.Unmarshal(bodyAsBytes, &r); err != nil {
			return errors.Wrap(err, "error unwrapping postmark error response body")
		}

		// https://postmarkapp.com/developer/api/overview#error-codes
		switch r.ErrorCode {
		case 406:
			s.log.Event("Not sending email, recipient is inactive", 1, "recipient", body.To)
			return nil
		default:
			s.log.Event("Error sending email, got error code", 1, "message", r.Message, "error code", r.ErrorCode)
			return errors.Newf("error sending email, got error code %v", r.ErrorCode)
		}
	}

	if response.StatusCode > 299 {
		s.log.Event("Error sending email, got http status code", 1, "status code", response.StatusCode, "body", string(bodyAsBytes))
		return errors.Newf("error sending email, got status %v", response.StatusCode)
	}

	return nil
}

// createNameAndEmail returns a name and email string ready for inserting into From and To fields.
func createNameAndEmail(name, email string) nameAndEmail {
	return fmt.Sprintf("%v <%v>", name, email)
}

//go:embed emails
var emails embed.FS

// getEmail from the given path, panicking on errors.
// It also replaces keywords given in the map.
// Email preheader text should be between 40-130 characters long.
func getEmail(path, preheader string, keywords Keywords) string {
	emailBody, err := emails.ReadFile("emails/" + path)
	if err != nil {
		panic(err)
	}

	layout, err := emails.ReadFile("emails/layout.html")
	if err != nil {
		panic(err)
	}

	email := string(layout)
	email = strings.ReplaceAll(email, "{{preheader}}", preheader)
	email = strings.ReplaceAll(email, "{{body}}", string(emailBody))

	if _, ok := keywords["unsubscribe"]; ok {
		email = strings.ReplaceAll(email, "{{unsubscribe}}", "{{{ pm:unsubscribe }}}")
	} else {
		email = strings.ReplaceAll(email, "{{unsubscribe}}", "")
	}

	for keyword, replacement := range keywords {
		email = strings.ReplaceAll(email, "{{"+keyword+"}}", replacement)
	}

	return email
}
