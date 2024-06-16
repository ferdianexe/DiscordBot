package entity

import "time"

// Loan is a loan entity.
type Loan struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float32   `json:"amount"`
	Term        int       `json:"term"`
	BillAmount  float32   `json:"bill_amount"`
	Outstanding float32   `json:"outstanding"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

// User is a user entity.
type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IsDelinquent bool   `json:"is_delinquent"`
}
