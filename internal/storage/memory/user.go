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
	// ErrTypeAssert error: can't provide type assertion.
	ErrTypeAssert = errors.New("can't provide type assertion")
)

type memUser struct {
	password string
	id       uuid.UUID
}

func newMemUser(password string) memUser {
	return memUser{password, uuid.New()}
}

// UserStore store for user in memory.
type UserStore struct {
	User sync.Map
}

// NewUserStore returns UserStore.
func NewUserStore() *UserStore {
	return &UserStore{sync.Map{}}
}

// AddUser save user in memory.
func (u *UserStore) AddUser(_ context.Context, user models.User) (uuid.UUID, error) {
	memUser := newMemUser(user.Password)
	_, ok := u.User.LoadOrStore(user.Login, memUser)
	if ok {
		return uuid.Nil, storage.ErrDuplicateLogin
	}
	return memUser.id, nil
}

// GetUser load user from memory.
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
	return mUser.id, nil
}
