package service

import (
	"context"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/NevostruevK/GophKeeper/internal/storage/memory"
)

func (s *Service) Register(ctx context.Context, u *models.User) (string, error) {
	s.memory = memory.NewDataStore()
	return s.client.Auth.Register(ctx, u.Login, u.Password)
}

func (s *Service) Login(ctx context.Context, u *models.User) (string, error) {
	s.memory = memory.NewDataStore()
	return s.client.Auth.Login(ctx, u.Login, u.Password)
}
