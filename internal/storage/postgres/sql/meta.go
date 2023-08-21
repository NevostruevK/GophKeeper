package sql

const (
	InsertMetaWithDescription = `
		INSERT INTO metas (id, user_id, type, title, has_description)
		VALUES ($1, $2, $3, $4, 'true')
		RETURNING id
	`

	InsertMeta = `
		INSERT INTO metas (id, user_id, type, title, data_size, description_size)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	SelectMetasForUser = `
		SELECT id, type, title, data_size, description_size
		FROM metas
		WHERE user_id = $1
	`

	SelectMetasForUserAndType = `
		SELECT id, type, title, data_size, description_size
		FROM metas
		WHERE user_id = $1 AND type = $2
	`
)
