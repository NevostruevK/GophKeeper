package postgres_test

import (
	"context"
	"fmt"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/google/uuid"
)

func genLoginPassword() func() (string, string) {
	var num int
	return func() (string, string) {
		num++
		login := fmt.Sprintf("test_login_%d", num)
		password := fmt.Sprintf("test_password_%d", num)
		return login, password
	}
}

var newLogPass = genLoginPassword()

func addUser(ctx context.Context) (uuid.UUID, error) {
	user := models.NewUser(newLogPass())
	return testStorage.AddUser(ctx, *user)
}
