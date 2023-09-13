package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User структура для хранения данных о пользователе.
type User struct {
	Login    string // login
	Password string // password
}

// UserDB структура для хранения данных о пользователе для базы данных.
type UserDB struct {
	ID    uuid.UUID // user ID
	Login string    // login
	Hash  []byte    // password hash
}

// NewUser returns User.
func NewUser(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}

// NewUserDB returns UserDB.
func NewUserDB(login, password string) (*UserDB, error) {
	u := NewUser(login, password)
	return u.UserToDB()
}

// UserToDB converts User to UserDB.
func (u User) UserToDB() (userDB *UserDB, err error) {
	userDB = &UserDB{
		ID:    uuid.New(),
		Login: u.Login,
	}
	userDB.Hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	return userDB, err
}

// CheckHash checks user's password.
func (u UserDB) CheckHash(password string) error {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password))
}

// IsReadyForStorage check User's fields for ready to store.
func (u *User) IsReadyForStorage() (bool, string) {
	const (
		loginIsEmpty    = "login is empty"
		passwordIsShort = "password is short (at least 4 characters)"
		passwordMinLen  = 4
	)
	if u.Login == "" {
		return false, loginIsEmpty
	}
	if len(u.Password) < passwordMinLen {
		return false, passwordIsShort
	}
	return true, ""
}
