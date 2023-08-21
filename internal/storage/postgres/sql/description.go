package sql

const (
	InsertDescription = `
		INSERT INTO descriptions (id, description, is_compressed)
		VALUES ($1, $2, $3)
	`

	SelectDescription = `
		SELECT description, is_compressed
		FROM descriptions
		WHERE id = $1
	`
)
