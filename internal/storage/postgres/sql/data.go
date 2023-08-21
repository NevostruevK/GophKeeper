package sql

const (
	InsertData = `
		INSERT INTO datas (id, data, key, nonce)
		VALUES ($1, $2, $3, $4)
	`

	SelectData = `
		SELECT data, key, nonce
		FROM datas
		WHERE id = $1
		LIMIT 1
	`
)
