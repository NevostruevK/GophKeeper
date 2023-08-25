package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/google/uuid"
)

var (
	ErrTypeAssert = errors.New("can't provide type assertion")
)

type memUser struct {
	password string
	ID       uuid.UUID
}

func newMemUser(password string) memUser {
	return memUser{password, uuid.New()}
}

type UserStore struct {
	User sync.Map
}

func NewUserStore() *UserStore {
	return &UserStore{sync.Map{}}
}

func (u *UserStore) AddUser(_ context.Context, user models.User) (uuid.UUID, error) {
	memUser := newMemUser(user.Password)
	_, ok := u.User.LoadOrStore(user.Login, memUser)
	if ok {
		return uuid.Nil, storage.ErrDuplicateLogin
	}
	return memUser.ID, nil
}

func (u *UserStore) GetUser(_ context.Context, user models.User) (uuid.UUID, error) {
	v, ok := u.User.Load(user.Login)
	if !ok {
		return uuid.Nil, storage.ErrNotFound
	}
	mUser, ok := v.(memUser)
	if !ok {
		return uuid.Nil, ErrTypeAssert
	}
	if mUser.password != user.Password {
		return uuid.Nil, storage.ErrWrongPassword
	}
	return mUser.ID, nil
}
