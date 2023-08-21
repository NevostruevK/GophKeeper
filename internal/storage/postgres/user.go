package postgres

import (
	"context"
	"errors"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/postgres/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) AddUser(ctx context.Context, user models.User) (id uuid.UUID, err error) {
	u, err := user.UserToDB()
	if err != nil {
		return uuid.Nil, err
	}
	if err = s.QueryRow(ctx, sql.InsertUser, u.ID, u.Login, u.Hash, u.Salt).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, ErrDuplicateLogin
		}
	}
	return id, err
}

func (s *Storage) GetUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	u := models.UserDB{User: user}
	err := s.QueryRow(ctx, sql.SelectUser, u.Login).Scan(&u.ID, &u.Hash, &u.Salt)
	if err != nil {
		return uuid.Nil, err
	}
	return u.ID, u.CheckHash()
}
