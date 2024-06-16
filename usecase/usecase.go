package usecase

import "github.com/jery1402/billing-engine/entity"

// repoProvider is the interface for the repository.
type repoProvider interface {
	// InsertLoan inserts a new loan in the database.
	InsertLoan(param entity.Loan) error

	// InsertUser inserts a new user in the database.
	InsertUser(param entity.User) error

	// GetActiveLoanByUserID retrieves a list of active loan information from the database based on the given user id.
	GetActiveLoanByUserID(userID int) ([]entity.Loan, error)

	// GetLoan retrieves a loan information based on the given id.
	GetLoan(id int) (entity.Loan, error)

	// GetLoanList retrieves a list of loan information from the database.
	GetLoanList() ([]entity.Loan, error)

	// GetUser retrieves a user information based on the given id.
	GetUser(id int) (entity.User, error)

	// GetUserList retrieves a list of user information from the database.
	GetUserList() ([]entity.User, error)

	// UpdateLoan updates a loan information (ex: payment).
	UpdateLoan(param entity.Loan) error

	// UpdateUser updates a user information based on the given id.
	UpdateUser(param entity.User) error

	// SetupDatabase creates a new database in case it doesn't exist.
	SetupDatabase() error
}

// UseCase is the usecase entity
type UseCase struct {
	repo repoProvider
}

// NewUseCase creates a new use case.
func NewUseCase(repo repoProvider) *UseCase {
	return &UseCase{repo: repo}
}
