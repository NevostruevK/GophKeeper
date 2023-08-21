package postgres

const (
	createUsersSQL = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			login varchar(64) NOT NULL CHECK (char_length(login)>8),
			hash varchar(64) NOT NULL,
			salt varchar(16) NOT NULL,
			CONSTRAINT users_login_unique UNIQUE(login)
		)
	`
	createMetasSQL = `
		DO $$ BEGIN
			CREATE TYPE type_enum AS ENUM ('PAIR', 'TEXT', 'FILE', 'CARD');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		CREATE TABLE IF NOT EXISTS metas (
			id UUID CONSTRAINT metas_pkey PRIMARY KEY,
			user_id UUID REFERENCES users(id) NOT NULL,
			type type_enum NOT NULL, 
			title varchar(256) NOT NULL, 
			data_size integer NOT NULL,
			description_size integer DEFAULT 0
		)
	`

	createDescriptionsSQL = `
		CREATE TABLE IF NOT EXISTS descriptions (
			id UUID REFERENCES metas(id) PRIMARY KEY,
			description text NOT NULL,
			is_compressed BOOL DEFAULT FALSE			
		)
	`
	createDatasSQL = `
		CREATE TABLE IF NOT EXISTS datas (
			id UUID REFERENCES metas(id) PRIMARY KEY,
			data bytea NOT NULL,
			key varchar(32) NOT NULL,
			nonce varchar(12) NOT NULL
		)
	`
)
