package sql

const (
	InsertRecord = `
		INSERT INTO records (id, user_id, type, title, data, data_size, has_description, description)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'true', $6)
		RETURNING id
	`

	InsertRecordWithoutDescription = `
		INSERT INTO records (id, user_id, type, title, data, data_size, has_description)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'false')
		RETURNING id
	`
	SelectSpecs = `
		SELECT id, type, title, data_size, has_description
		FROM records
		WHERE user_id = $1
	`

	SelectSpecsOfType = `
		SELECT id, type, title, data_size, has_description
		FROM records
		WHERE user_id = $1 AND type = $2
	`
	SelectData = `
		SELECT data
		FROM records
		WHERE id = $1
		LIMIT 1
	`
)
