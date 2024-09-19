package sql_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"maragu.dev/is"

	"maragu.dev/goo/model"
	"maragu.dev/goo/sqltest"
)

func TestHelper_Signup(t *testing.T) {
	t.Run("signs up an account and user, and creates a token", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "Me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		is.Equal(t, 34, len(u.ID))
		is.True(t, strings.HasPrefix(u.ID.String(), "u_"))
		is.True(t, time.Since(u.Created.T) < time.Second)
		is.True(t, time.Since(u.Updated.T) < time.Second)
		is.Equal(t, "Me", u.Name)
		is.Equal(t, "me@example.com", u.Email.String())
		is.True(t, !u.Confirmed)
		is.True(t, u.Active)

		var a model.Account
		err = h.Get(context.Background(), &a, `select * from accounts where id = ?`, u.AccountID)
		is.NotError(t, err)
		is.Equal(t, 34, len(a.ID))
		is.True(t, strings.HasPrefix(a.ID.String(), "a_"))
		is.True(t, time.Since(a.Created.T) < time.Second)
		is.True(t, time.Since(a.Updated.T) < time.Second)
		is.Equal(t, "Me", a.Name)

		var token string
		err = h.Get(context.Background(), &token, `select value from tokens where userID = ?`, u.ID)
		is.NotError(t, err)
		is.Equal(t, 34, len(token))
		is.True(t, strings.HasPrefix(token, "t_"))
	})

	t.Run("errors on duplicate email", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "Me@example.com",
		}
		_, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		_, err = h.Signup(context.Background(), u)
		is.Error(t, model.ErrorEmailConflict, err)
	})
}

func TestHelper_Login(t *testing.T) {
	t.Run("marks token used and user confirmed and returns user", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "Me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		userID := u.ID

		var token string
		err = h.Get(context.Background(), &token, `select value from tokens where userID = ?`, userID)
		is.NotError(t, err)

		u, err = h.Login(context.Background(), token)
		is.NotError(t, err)
		is.Equal(t, userID, u.ID)
		is.Equal(t, "Me", u.Name)

		var used bool
		err = h.Get(context.Background(), &used, `select used from tokens where value = ?`, token)
		is.NotError(t, err)
		is.True(t, used)

		is.True(t, u.Confirmed)
	})

	t.Run("can login twice", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "Me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		userID := u.ID

		var token string
		err = h.Get(context.Background(), &token, `select value from tokens where userID = ?`, userID)
		is.NotError(t, err)

		_, err = h.Login(context.Background(), token)
		is.NotError(t, err)

		_, err = h.Login(context.Background(), token)
		is.NotError(t, err)
	})

	t.Run("returns token expired error when token is expired", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "Me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		var token string
		err = h.Get(context.Background(), &token, `select value from tokens where userID = ?`, u.ID)
		is.NotError(t, err)

		err = h.Exec(context.Background(), `update tokens set expires = '2001-01-01T00:00:00.000Z' where value = ?`, token)
		is.NotError(t, err)

		_, err = h.Login(context.Background(), token)
		is.Error(t, model.ErrorTokenExpired, err)
	})

	t.Run("returns user inactive error when user not active", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "Me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		var token string
		err = h.Get(context.Background(), &token, `select value from tokens where userID = ?`, u.ID)
		is.NotError(t, err)

		err = h.Exec(context.Background(), `update users set active = false where id = ?`, u.ID)
		is.NotError(t, err)

		_, err = h.Login(context.Background(), token)
		is.Error(t, model.ErrorUserInactive, err)
	})

	t.Run("returns error if no such token", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		_, err := h.Login(context.Background(), "t_018743ccfbed090e8f5eebc810ff797d")
		is.Error(t, model.ErrorTokenNotFound, err)
	})
}

func TestHelper_TryLogin(t *testing.T) {
	t.Run("creates token", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		err = h.TryLogin(context.Background(), "me@example.com")
		is.NotError(t, err)

		var token string
		err = h.Get(context.Background(), &token, `select value from tokens where userID = ?`, u.ID)
		is.NotError(t, err)
	})

	t.Run("errors when user not found", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		err := h.TryLogin(context.Background(), "doesnotexist@example.com")
		is.Error(t, model.ErrorUserNotFound, err)
	})

	t.Run("errors when user inactive", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		err = h.Exec(context.Background(), `update users set active = false where id = ?`, u.ID)
		is.NotError(t, err)

		err = h.TryLogin(context.Background(), "me@example.com")
		is.Error(t, model.ErrorUserInactive, err)
	})
}

func TestHelper_GetUser(t *testing.T) {
	t.Run("returns user", func(t *testing.T) {
		h := sqltest.NewHelper(t)

		u := model.User{
			Name:  "Me",
			Email: "me@example.com",
		}
		u, err := h.Signup(context.Background(), u)
		is.NotError(t, err)

		u2, err := h.GetUser(context.Background(), u.ID)
		is.NotError(t, err)

		is.Equal(t, u.ID, u2.ID)
		is.Equal(t, u.Name, u2.Name)
	})
}
