package sqlite

import "database/sql"

const (
	createUserTableQuery = `
			CREATE TABLE IF NOT EXISTS user (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL
			);
	`
)

// SetupDatabase creates the database tables if they don't exist.
func (repo *Repository) SetupDatabase() error {
	_, err := repo.db.Exec(createUserTableQuery)
	if err != nil {
		return err
	}

	return nil
}

// sqliteProvider is the interface for the sqlite database.
type sqliteProvider interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// Repository is the repository entity.
type Repository struct {
	db sqliteProvider
}

// NewRepository creates a new repository.
func NewRepository(db sqliteProvider) *Repository {
	return &Repository{db: db}
}
