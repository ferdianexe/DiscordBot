package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jery1402/billing-engine/entity"
	"github.com/jery1402/billing-engine/usecase"
)

// CreateUser creates a new user.
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodPost {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	var user entity.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return

	}

	err = h.usecase.InsertUser(user)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = "success to create user"
	writeResponse(w, http.StatusOK, response)
}

// GetLoanList retrieves a list of loan information.
func (h *Handler) GetLoanList(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodGet {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	loans, err := h.usecase.GetLoanList()
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeResponse(w, http.StatusOK, loans)
}

// GetNextPayment returns the next payment date and the payment amount.
func (h *Handler) GetNextPayment(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodGet {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	loanID, err := strconv.Atoi(r.URL.Query().Get("loan_id"))
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	result, err := h.usecase.GetNextPayment(loanID)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = result
	writeResponse(w, http.StatusOK, response)
}

// GetOutstanding retrieves the outstanding amount of a loan.
func (h *Handler) GetOutstanding(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodGet {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	loanID, err := strconv.Atoi(r.URL.Query().Get("loan_id"))
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	loan, err := h.usecase.GetLoan(loanID)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = fmt.Sprintf("outstanding amount of loan ID %d is %0.2f", loanID, loan.Outstanding)
	writeResponse(w, http.StatusOK, response)
}

// GetUserDelinquentStatus retrieves the delinquent status of a user.
func (h *Handler) GetUserDelinquentStatus(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodGet {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	user, err := h.usecase.GetUser(userID)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = "user is not delinquent"
	if user.IsDelinquent {
		response.Message = "user is delinquent"
	}
	writeResponse(w, http.StatusOK, response)
}

// GetUserList retrieves a list of user information.
func (h *Handler) GetUserList(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodGet {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	users, err := h.usecase.GetUserList()
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	writeResponse(w, http.StatusOK, users)
}

// MakeLoan makes a loan.
func (h *Handler) MakeLoan(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodPost {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	var loan LoanRequest

	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return

	}

	err = validateLoanRequest(loan)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	err = h.usecase.CreateLoan(entity.Loan{
		UserID: loan.UserID,
		Amount: loan.Amount,
		Term:   loan.Term,
	})
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = "success to create a loan"
	writeResponse(w, http.StatusOK, response)
}

// MakePayment makes a payment.
func (h *Handler) MakePayment(w http.ResponseWriter, r *http.Request) {
	response := GeneralResponse{}

	if r.Method != http.MethodPost {
		response.Message = "method not allowed"
		writeResponse(w, http.StatusMethodNotAllowed, response)
		return
	}

	var payment PaymentRequest

	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return

	}

	err = validatePaymentRequest(payment)
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusBadRequest, response)
		return
	}

	err = h.usecase.MakePayment(usecase.MakePaymentParam{
		LoanID:      payment.LoanID,
		UserID:      payment.UserID,
		Amount:      payment.Amount,
		PaymentDate: payment.PaymentDate,
	})
	if err != nil {
		response.Message = err.Error()
		writeResponse(w, http.StatusInternalServerError, response)
		return
	}

	response.Message = "success to make a payment"
	writeResponse(w, http.StatusOK, response)
}

// validateLoanRequest validates the loan request.
func validateLoanRequest(loan LoanRequest) error {
	if loan.UserID == 0 {
		return errors.New("user id is required")
	}

	if loan.Amount == 0 {
		return errors.New("amount is required")
	}

	if loan.Term == 0 {
		return errors.New("term is required")
	}

	return nil
}

// validatePaymentRequest validates the payment request.
func validatePaymentRequest(payment PaymentRequest) error {
	if payment.LoanID == 0 {
		return errors.New("loan id is required")
	}

	if payment.UserID == 0 {
		return errors.New("user id is required")
	}

	if payment.Amount == 0 {
		return errors.New("amount is required")
	}

	if payment.PaymentDate != "" {
		_, err := time.Parse("2006-01-02", payment.PaymentDate)
		if err != nil {
			return errors.New("payment date is invalid")
		}
	}

	return nil
}
