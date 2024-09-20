package sql

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"

	"maragu.dev/errors"

	"maragu.dev/goo/model"
)

// Signup creates an account, an unconfirmed user, and a token.
// Returns the new user.
func (h *Helper) Signup(ctx context.Context, u model.User) (model.User, error) {
	err := h.InTransaction(ctx, func(tx *Tx) error {
		var exists bool
		query := `select exists (select * from users where email = ?)`
		if err := tx.Get(ctx, &exists, query, u.Email.ToLower()); err != nil {
			return errors.Wrap(err, "error getting user by email")
		}
		if exists {
			return model.ErrorEmailConflict
		}

		token, err := createToken()
		if err != nil {
			return err
		}

		u.Name = strings.TrimSpace(u.Name)

		var a model.Account
		query = `insert into accounts (name) values (?) returning *`
		if err := tx.Get(ctx, &a, query, u.Name); err != nil {
			return errors.Wrap(err, "error creating account")
		}

		query = `insert into users (accountID, name, email) values (?, ?, ?) returning *`
		if err := tx.Get(ctx, &u, query, a.ID, u.Name, u.Email.ToLower()); err != nil {
			return errors.Wrap(err, "error creating user")
		}

		query = `insert into tokens (value, userID) values (?, ?)`
		if err := tx.Exec(ctx, query, token, u.ID); err != nil {
			return errors.Wrap(err, "error creating token")
		}

		m := stringMap{
			"type":  "signup",
			"name":  u.Name,
			"email": u.Email.ToLower().String(),
			"token": token,
		}
		if err := h.CreateJobInTx(ctx, tx, "send-email", m); err != nil {
			return errors.Wrap(err, "error creating job to send signup email")
		}

		return nil
	})
	return u, err
}

// Login with the given token. It marks the token as used (but this isn't currently checked anywhere)
// if it's not expired and if the user is marked active.
// It also sets the user confirmed.
// Returns the user.
func (h *Helper) Login(ctx context.Context, token string) (model.User, error) {
	var u model.User
	err := h.InTransaction(ctx, func(tx *Tx) error {
		var expired bool
		query := `select exists (select 1 from tokens where value = ? and expires <= strftime('%Y-%m-%dT%H:%M:%fZ'))`
		if err := tx.Get(ctx, &expired, query, token); err != nil {
			return err
		}
		if expired {
			return model.ErrorTokenExpired
		}

		var inactive bool
		query = `select exists (select 1 from users where id = (select userID from tokens where value = ?) and not active)`
		if err := tx.Get(ctx, &inactive, query, token); err != nil {
			return err
		}
		if inactive {
			return model.ErrorUserInactive
		}

		var userID model.ID
		query = `update tokens set used = 1 where value = ? returning userID`
		if err := tx.Get(ctx, &userID, query, token); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return model.ErrorTokenNotFound
			}
			return err
		}

		query = `update users set confirmed = 1 where id = ? returning *`
		if err := tx.Get(ctx, &u, query, userID); err != nil {
			return err
		}

		return nil
	})
	return u, err
}

// TryLogin checks whether the user exists and is active, creates a login token, and creates a job to send
// an email with the token in it.
func (h *Helper) TryLogin(ctx context.Context, email model.Email) error {
	return h.InTransaction(ctx, func(tx *Tx) error {
		var exists bool
		query := `select exists (select 1 from users where email = ?)`
		if err := tx.Get(ctx, &exists, query, email); err != nil {
			return err
		}
		if !exists {
			return model.ErrorUserNotFound
		}

		var inactive bool
		query = `select exists (select 1 from users where email = ? and not active)`
		if err := tx.Get(ctx, &inactive, query, email); err != nil {
			return err
		}
		if inactive {
			return model.ErrorUserInactive
		}

		token, err := createToken()
		if err != nil {
			return err
		}
		query = `insert into tokens (value, userID) values (?, (select id from users where email = ?))`
		if err := tx.Exec(ctx, query, token, email); err != nil {
			return errors.Wrap(err, "error creating token")
		}

		var u model.User
		query = `select * from users where email = ?`
		if err := tx.Get(ctx, &u, query, email); err != nil {
			return errors.Wrap(err, "error getting user")
		}

		m := stringMap{
			"type":  "login",
			"name":  u.Name,
			"email": u.Email.ToLower().String(),
			"token": token,
		}
		if err := h.CreateJobInTx(ctx, tx, "send-email", m); err != nil {
			return errors.Wrap(err, "error creating job to send login email")
		}

		return nil
	})
}

func (h *Helper) GetUser(ctx context.Context, id model.ID) (model.User, error) {
	var u model.User
	if err := h.Get(ctx, &u, `select * from users where id = ?`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, model.ErrorUserNotFound
		}
		return u, err
	}

	return u, nil
}

func createToken() (string, error) {
	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	return fmt.Sprintf("t_%x", secret), nil
}
