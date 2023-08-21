package sql

const (
	InsertUser = `
		INSERT INTO users (id, login, hash, salt)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT ON CONSTRAINT users_login_unique
			DO NOTHING	
		RETURNING id
	`
	SelectUser = `
		SELECT id, hash, salt
		FROM users
		WHERE login = $1
		LIMIT 1
	`
)
