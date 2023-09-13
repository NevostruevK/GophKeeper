package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AddUser adds user to database.
func (s *Storage) AddUser(ctx context.Context, user models.User) (id uuid.UUID, err error) {
	u, err := user.UserToDB()
	if err != nil {
		return uuid.Nil, err
	}
	if err = s.QueryRow(ctx, sql.InsertUser, u.ID, u.Login, u.Hash).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, storage.ErrDuplicateLogin
		}
	}
	return id, err
}

// GetUser get user from database.
func (s *Storage) GetUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	u := models.UserDB{Login: user.Login}
	err := s.QueryRow(ctx, sql.SelectUser, u.Login).Scan(&u.ID, &u.Hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, storage.ErrNotFound
		}
		return uuid.Nil, err
	}
	err = u.CheckHash(user.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%v : %w", storage.ErrWrongPassword, err)
	}
	return u.ID, nil
}
