package tui

import (
	"testing"

	"github.com/NevostruevK/GophKeeper/internal/models"
	"github.com/stretchr/testify/require"
)

func TestUserForm(t *testing.T) {
	service, server, client, err := startService()
	require.NoError(t, err)
	defer func() {
		err = stopService(server, client)
		require.NoError(t, err)
	}()
	ch := make(chan any)
	go func() {
		err := run(service, "", "", ch)
		require.NoError(t, err)
	}()
	defer app.Stop()
	<-ch
	userFormTest := newUserForm()
	tests := []struct {
		name    string
		user    *models.User
		isLogin bool
		want    string
	}{
		{
			name:    "register ok",
			user:    models.NewUser("register ok", "register ok"),
			isLogin: false,
			want:    "register ok register ok\n",
		},
		{
			name:    "login ok",
			user:    models.NewUser("register ok", "register ok"),
			isLogin: true,
			want:    "register ok login ok\n",
		},
		{
			name:    "register err (login empty)",
			user:    models.NewUser("", "register err"),
			isLogin: false,
			want:    "[Warning] login is empty\n",
		},
		{
			name:    "login err (not found)",
			user:    models.NewUser("login err (not found)", "register err"),
			isLogin: true,
			want:    errIncorrectLoginPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userFormTest.user = tt.user
			userFormTest.userRequst(tt.isLogin)
			got := messager.GetText(false)
			if got != tt.want {
				t.Errorf("userRequst got %s, want %s", got, tt.want)
			}
		})
	}
}
