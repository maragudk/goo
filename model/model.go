package model

type ID string

func (i ID) String() string {
	return string(i)
}

type Account struct {
	ID      ID
	Created Time
	Updated Time
	Name    string
}

type User struct {
	ID        ID
	Created   Time
	Updated   Time
	AccountID ID `db:"accountID"`
	Name      string
	Email     Email
	Confirmed bool
	Active    bool
}
