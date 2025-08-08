package repository

func (r *Repository) SetupSchemas() (err error) {
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		profile_image_url TEXT
	);
	CREATE TABLE IF NOT EXISTS channels (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL		
	);
	`
	_, err = r.db.Exec(createTablesSQL)

	return err
}
