package models

import (
	"errors"
	"fmt"

	"github.com/NevostruevK/GophKeeper/internal/utils/hash"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const (
	SaltSize = 16
	HashSize = 64
)

const Pepper = "secret pepper"

var ErrWrongHashForUser = errors.New("wrong hash for user: ")

type User struct {
	Login    string
	Password string
}

type UserDB struct {
	ID uuid.UUID
	User
	Hash string
	Salt string
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

func (u User) UserToDB() (*UserDB, error) {
	userDB := &UserDB{
		ID:   uuid.New(),
		User: u,
		Salt: lo.RandomString(SaltSize, lo.AllCharset),
	}
	err := userDB.CountHash()
	return userDB, err
}

func (u UserDB) CheckHash() error {
	h, err := CountHash(u)
	if err != nil {
		return err
	}
	if h != u.Hash {
		// TODO добавлю событие в лог
		fmt.Printf("TODO: %v %v != %s\n", ErrWrongHashForUser, u, h)
		return ErrWrongHashForUser
	}
	return nil
}
func CountHash(u UserDB) (string, error) {
	h, err := hash.Hash(fmt.Sprintf("%s:%s:%s", u.ID, u.Login, u.Password), u.Salt)
	if err != nil {
		return "", err
	}
	return hash.Hash(h, Pepper)
}

func (u *UserDB) CountHash() error {
	h, err := hash.Hash(fmt.Sprintf("%s:%s:%s", u.ID, u.Login, u.Password), u.Salt)
	if err != nil {
		return err
	}
	u.Hash, err = hash.Hash(h, Pepper)
	return err
}
