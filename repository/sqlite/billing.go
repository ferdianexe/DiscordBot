package sqlite

import (
	"database/sql"
	"time"

	"github.com/jery1402/billing-engine/entity"
)

const (
	getActiveLoanByUserIDQuery = `
		SELECT id, user_id, amount, term, bill_amount, outstanding, create_time, update_time
		FROM loan
		WHERE user_id = ? AND outstanding > 0
	`

	getLoanListQuery = `
		SELECT id, user_id, amount, term, bill_amount, outstanding, create_time, update_time
		FROM loan
	`

	getLoanQuery = `
		SELECT id, user_id, amount, term, bill_amount, outstanding, create_time, update_time
		FROM loan
		WHERE id = ?;
	`

	getUserQuery = `
		SELECT id, name, is_delinquent
		FROM user
		WHERE id = ?;
	`

	getUserListQuery = `
		SELECT id, name, is_delinquent
		FROM user
	`

	insertLoanQuery = `
		INSERT INTO loan (user_id, amount, term, bill_amount, outstanding, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	insertUserQuery = `
		INSERT INTO user (name, is_delinquent)
		VALUES (?, ?);
	`

	updateLoanQuery = `
		UPDATE loan
		SET outstanding = ?, update_time = ?
		WHERE id = ?;
	`

	updateUserQuery = `
		UPDATE user
		SET is_delinquent = ?
		WHERE id = ?;
	`
)

// GetActiveLoanByUserID retrieves a list of active loan information from the database based on the given user id.
func (repo *Repository) GetActiveLoanByUserID(userID int) ([]entity.Loan, error) {
	rows, err := repo.db.Query(getActiveLoanByUserIDQuery, userID)
	if err != nil {
		return []entity.Loan{}, err
	}
	defer rows.Close()

	var loans []entity.Loan
	for rows.Next() {
		var loan entity.Loan
		err = rows.Scan(&loan.ID, &loan.UserID, &loan.Amount, &loan.Term, &loan.BillAmount, &loan.Outstanding, &loan.CreateTime, &loan.UpdateTime)
		if err != nil {
			return []entity.Loan{}, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

// GetLoan retrieves a loan information from the database.
func (repo *Repository) GetLoan(id int) (entity.Loan, error) {
	row := repo.db.QueryRow(getLoanQuery, id)

	var loan entity.Loan
	err := row.Scan(&loan.ID, &loan.UserID, &loan.Amount, &loan.Term, &loan.BillAmount, &loan.Outstanding, &loan.CreateTime, &loan.UpdateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Loan{}, nil
		}
		return entity.Loan{}, err
	}

	return loan, nil
}

// GetLoanList retrieves a list of loan information from the database.
func (repo *Repository) GetLoanList() ([]entity.Loan, error) {
	rows, err := repo.db.Query(getLoanListQuery)
	if err != nil {
		return []entity.Loan{}, err
	}
	defer rows.Close()

	var loans []entity.Loan
	for rows.Next() {
		var loan entity.Loan
		err = rows.Scan(&loan.ID, &loan.UserID, &loan.Amount, &loan.Term, &loan.BillAmount, &loan.Outstanding, &loan.CreateTime, &loan.UpdateTime)
		if err != nil {
			return []entity.Loan{}, err
		}
		loans = append(loans, loan)
	}

	return loans, nil
}

// GetUser retrieves a user information from the database.
func (repo *Repository) GetUser(id int) (entity.User, error) {
	row := repo.db.QueryRow(getUserQuery, id)

	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.IsDelinquent)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}

	return user, nil
}

// GetUserList retrieves a list of user information from the database.
func (repo *Repository) GetUserList() ([]entity.User, error) {
	rows, err := repo.db.Query(getUserListQuery)
	if err != nil {
		return []entity.User{}, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.IsDelinquent)
		if err != nil {
			return []entity.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

// InsertLoan inserts a new loan in the database.
func (repo *Repository) InsertLoan(param entity.Loan) error {
	currentTime := time.Now()

	param.CreateTime = currentTime
	param.UpdateTime = currentTime

	_, err := repo.db.Exec(insertLoanQuery, param.UserID, param.Amount, param.Term, param.BillAmount, param.Outstanding, param.CreateTime, param.UpdateTime)
	return err
}

// InsertUser inserts a new user in the database.
func (repo *Repository) InsertUser(param entity.User) error {
	_, err := repo.db.Exec(insertUserQuery, param.Name, param.IsDelinquent)
	return err
}

// UpdateLoan updates a loan information in the database.
func (repo *Repository) UpdateLoan(param entity.Loan) error {
	_, err := repo.db.Exec(updateLoanQuery, param.Outstanding, param.UpdateTime, param.ID)
	return err
}

// UpdateUser updates a user information in the database.
func (repo *Repository) UpdateUser(param entity.User) error {
	_, err := repo.db.Exec(updateUserQuery, param.IsDelinquent, param.ID)
	return err
}
