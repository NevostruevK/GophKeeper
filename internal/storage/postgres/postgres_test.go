package postgres_test

import (
	"context"
	"fmt"

	"github.com/NevostruevK/GophKeeper/internal/models"
	storage "github.com/NevostruevK/GophKeeper/internal/storage/postgres"
	"github.com/NevostruevK/GophKeeper/internal/tools/crypto"
	"github.com/google/uuid"
	"github.com/samber/lo"
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

type idsDB struct {
	ids []uuid.UUID
}

func addUser(ctx context.Context, st *storage.Storage, ids *idsDB) (*models.UserDB, error) {
	user := models.NewUser(newLogPass())
	id, err := st.AddUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	ids.ids = append(ids.ids, id)
	return &models.UserDB{Login: user.Login, ID: id}, err
}

func deleteFromDB(ctx context.Context, st *storage.Storage, ids []uuid.UUID) error {
	for _, id := range ids {
		_, err := st.Exec(ctx, "DELETE FROM records where user_id = $1", id)
		if err != nil {
			return err
		}
		_, err = st.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
		if err != nil {
			return err
		}
	}
	return nil
}

func newStorage(ctx context.Context) (*storage.Storage, error) {
	key := lo.RandomString(32, lo.AllCharset)
	nonce := lo.RandomString(12, lo.AllCharset)
	c, err := crypto.NewCrypto([]byte(key), []byte(nonce))
	if err != nil {
		return nil, err
	}
	return storage.NewStorage(ctx, "user=postgres sslmode=disable", c)
}
