package model

// Error is for errors in the business domain. See the constants below.
type Error string

const (
	ErrorEmailConflict = Error("EMAIL_CONFLICT")
	ErrorTokenExpired  = Error("TOKEN_EXPIRED")
	ErrorUserInactive  = Error("USER_INACTIVE")
	ErrorTokenNotFound = Error("TOKEN_NOT_FOUND")
	ErrorUserNotFound  = Error("USER_NOT_FOUND")
)

func (e Error) Error() string {
	return string(e)
}
