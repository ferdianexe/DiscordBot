package sqlite

import "database/sql"

const (
	createLoanTableQuery = `
		CREATE TABLE IF NOT EXISTS loan (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			amount FLOAT NOT NULL,
			term INTEGER NOT NULL,
			bill_amount FLOAT NOT NULL,
			outstanding FLOAT NOT NULL,
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	createUserTableQuery = `
			CREATE TABLE IF NOT EXISTS user (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				is_delinquent BOOLEAN NOT NULL
			);
	`
)

// SetupDatabase creates the database tables if they don't exist.
func (repo *Repository) SetupDatabase() error {
	_, err := repo.db.Exec(createLoanTableQuery)
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(createUserTableQuery)
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
