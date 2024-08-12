package sqlite

import (
	"database/sql"
	"ferdianexe/DiscordBot/entity"
)

const (
	getLoanQuery = `
		SELECT *
		FROM user
		WHERE id = ?;
	`
)

func (repo *Repository) GetUserByID(id int) (entity.User, error) {
	row := repo.db.QueryRow(getLoanQuery, id)

	var usr entity.User
	err := row.Scan(&usr.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}

	return usr, nil
}
