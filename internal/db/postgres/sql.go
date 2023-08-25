package postgres

const (
	createUsersSQL = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			login varchar(64) NOT NULL CHECK (char_length(login)>8),
			hash varchar(64) NOT NULL,
			CONSTRAINT users_login_unique UNIQUE(login)
		)
	`
	createRecordsSQL = `
		DO $$ BEGIN
			CREATE TYPE type_enum AS ENUM ('PAIR', 'TEXT', 'FILE', 'CARD');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		CREATE TABLE IF NOT EXISTS records (
			id UUID CONSTRAINT records_pkey PRIMARY KEY,
			user_id UUID REFERENCES users(id) NOT NULL,
			type type_enum NOT NULL, 
			title varchar(256) NOT NULL,
			data bytea NOT NULL,
			data_size integer NOT NULL,
			has_description bool NOT NULL,
			description text
		)
	`
)
