package usecase

// GetNextPaymentResponse is the response entity for getting the next payment date and amount.
type GetNextPaymentResponse struct {
	NextPaymentDate string  `json:"next_payment_date"`
	PaymentAmount   float32 `json:"payment_amount"`
	IsLatePayment   bool    `json:"is_late_payment"`
}

// MakePaymentParam is the parameter for making a payment.
type MakePaymentParam struct {
	LoanID int
	UserID int
	Amount float32

	// To override payment date.
	// If not set, the payment date will be the current time.
	// Format: 2006-01-02
	PaymentDate string
}

// MakePaymentResponse is the response entity for making a payment.
type MakePaymentResponse struct {
	LoanID    int    `json:"loan_id"`
	IsSuccess bool   `json:"is_success"`
	Message   string `json:"message"`
}
