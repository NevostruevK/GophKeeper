package models

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string
	Password string
}

type UserDB struct {
	ID    uuid.UUID
	Login string
	Hash  []byte
}

func NewUser(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}

func NewUserDB(login, password string) (*UserDB, error) {
	u := NewUser(login, password)
	return u.UserToDB()
}

func (u User) UserToDB() (userDB *UserDB, err error) {
	userDB = &UserDB{
		ID:    uuid.New(),
		Login: u.Login,
	}
	userDB.Hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	return userDB, err
}

func (u UserDB) CheckHash(password string) error {
	return bcrypt.CompareHashAndPassword(u.Hash, []byte(password))
}
