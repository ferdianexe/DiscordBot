package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jery1402/billing-engine/entity"
	"github.com/jery1402/billing-engine/usecase"
)

// usecaseProvider is the interface for the usecase.
type usecaseProvider interface {
	// CreateDatabase creates a new database in case it doesn't exist.
	CreateDatabase() error

	// CreateLoan creates a new loan.
	CreateLoan(param entity.Loan) error

	// CreateUser creates a new user.
	InsertUser(param entity.User) error

	// GetActiveLoanByUserID retrieves a list of active loan information from the database based on the given user id.
	GetActiveLoanByUserID(userID int) ([]entity.Loan, error)

	// GetLoan retrieves a loan information based on the given id.
	GetLoan(id int) (entity.Loan, error)

	// GetLoanList retrieves a list of loan information from the database.
	GetLoanList() ([]entity.Loan, error)

	// GetNextPayment returns the next payment date and the payment amount.
	GetNextPayment(loanID int) (string, error)

	// GetUser retrieves a user information based on the given id.
	GetUser(id int) (entity.User, error)

	// GetUserList retrieves a list of user information from the database.
	GetUserList() ([]entity.User, error)

	// UpdateLoan updates a loan information (ex: payment).
	MakePayment(param usecase.MakePaymentParam) error

	// UpdateUser updates a user information based on the given id.
	UpdateUser(param entity.User) error
}

// Handler is the handler entity.
type Handler struct {
	usecase usecaseProvider
}

// NewHandler creates a new handler.
func NewHandler(usecase usecaseProvider) *Handler {
	return &Handler{usecase: usecase}
}

// CreateDatabase creates a new database in case it doesn't exist.
func (h *Handler) CreateDatabase(w http.ResponseWriter, r *http.Request) {
	err := h.usecase.CreateDatabase()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, GeneralResponse{Message: err.Error()})
		return
	}

	writeResponse(w, http.StatusCreated, GeneralResponse{Message: "success"})
}

// writeResponse writes a HTTP response.
func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonResponse, _ := json.Marshal(data)
	w.Write(jsonResponse)
}
