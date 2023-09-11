// package service сервис для синхонизации данных.
package service

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
)

// Register register new user.
func (s *Service) Register(ctx context.Context, u *models.User) (string, error) {
	s.memory = memory.NewDataStore()
	return s.client.Auth.Register(ctx, u.Login, u.Password)
}

// Login authorisation user.
func (s *Service) Login(ctx context.Context, u *models.User) (string, error) {
	s.memory = memory.NewDataStore()
	return s.client.Auth.Login(ctx, u.Login, u.Password)
}
