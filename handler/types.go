package handler

// GeneralResponse is the general response entity.
type GeneralResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GetOutstandingResponse is the response entity for getting the outstanding amount of a loan.
type GetOutstandingResponse struct {
	LoanID      int     `json:"loan_id"`
	Outstanding float32 `json:"outstanding"`
}

// GetUserDelinquentStatusResponse is the response entity for getting the delinquent status of a user.
type GetUserDelinquentStatusResponse struct {
	UserID       int  `json:"user_id"`
	IsDelinquent bool `json:"is_delinquent"`
}

// LoanRequest is the loan request entity.
type LoanRequest struct {
	UserID int     `json:"user_id"`
	Amount float32 `json:"amount"`
	Term   int     `json:"term"`
}

// PaymentRequest is the payment request entity.
type PaymentRequest struct {
	LoanID int     `json:"loan_id"`
	UserID int     `json:"user_id"`
	Amount float32 `json:"amount"`

	// To override payment date.
	// If not set, the payment date will be the current time.
	// Format: 2006-01-02
	PaymentDate string `json:"payment_date"`
}
