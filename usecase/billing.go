package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/jery1402/billing-engine/entity"
)

const (
	interestRate = 0.1
)

// CreateDatabase creates a new database in case it doesn't exist.
func (uc *UseCase) CreateDatabase() error {
	return uc.repo.SetupDatabase()
}

// CreateLoan creates a new loan.
// If the user is delinquent, an error is returned.
func (uc *UseCase) CreateLoan(param entity.Loan) error {
	user, err := uc.repo.GetUser(param.UserID)
	if err != nil {
		return errors.New("uc.repo.GetUser() error - " + err.Error())
	}
	if user.IsDelinquent {
		return errors.New("user is delinquent")
	}

	param.BillAmount = param.Amount / float32(param.Term) * (1 + interestRate)
	param.Outstanding = param.Amount

	return uc.repo.InsertLoan(param)
}

// GetLoanByUserID retrieves a list of active loan information from the database based on the given user id.
func (uc *UseCase) GetActiveLoanByUserID(userID int) ([]entity.Loan, error) {
	return uc.repo.GetActiveLoanByUserID(userID)
}

// CreateUser creates a new user.
func (uc *UseCase) InsertUser(param entity.User) error {
	return uc.repo.InsertUser(param)
}

// GetLoan retrieves a loan information based on the given id.
func (uc *UseCase) GetLoan(id int) (entity.Loan, error) {
	return uc.repo.GetLoan(id)
}

func (uc *UseCase) GetLoanList() ([]entity.Loan, error) {
	return uc.repo.GetLoanList()
}

// GetNextPayment returns the next payment date and the payment amount.
func (uc *UseCase) GetNextPayment(loanID int) (GetNextPaymentResponse, error) {
	loan, err := uc.repo.GetLoan(loanID)
	if err != nil {
		return GetNextPaymentResponse{}, err
	}

	currentTime := time.Now()

	var response GetNextPaymentResponse

	// truncate the time so we only compare the day.
	// nextPaymentDate example: loan that was created on 2020-01-01 will have a next payment due on 2020-01-08.
	nextPaymentDate := loan.CreateTime.Truncate(24*time.Hour).AddDate(0, 0, int(loan.Amount-loan.Outstanding+1)*7)

	if currentTime.Truncate(24*time.Hour).Sub(nextPaymentDate).Hours()/24 > 14 {
		_ = uc.repo.UpdateUser(entity.User{ID: loan.UserID, IsDelinquent: true})
		response.IsLatePayment = true
	}

	amount := loan.BillAmount

	if response.IsLatePayment {
		multiplier := int(currentTime.Truncate(24*time.Hour).Sub(nextPaymentDate).Hours()/24) / 7
		amount = loan.BillAmount * float32(multiplier)
		if amount > loan.Outstanding {
			amount = loan.Outstanding * (1 + interestRate)
		}
	}

	response.PaymentAmount = amount
	response.NextPaymentDate = nextPaymentDate.Format("2006-01-02")

	return response, nil
}

// GetUser retrieves a user information based on the given id.
func (uc *UseCase) GetUser(id int) (entity.User, error) {
	user, err := uc.repo.GetUser(id)
	if err != nil {
		return entity.User{}, errors.New("uc.repo.GetUser() error - " + err.Error())
	}

	if user.ID == 0 {
		return entity.User{}, errors.New("user not found")
	}

	return user, nil
}

// GetUserList retrieves a list of user information from the database.
func (uc *UseCase) GetUserList() ([]entity.User, error) {
	return uc.repo.GetUserList()
}

// UpdateLoan updates a loan information (ex: payment).
func (uc *UseCase) MakePayment(param MakePaymentParam) (MakePaymentResponse, error) {
	var (
		isDelinquent bool
		response     MakePaymentResponse
	)

	user, err := uc.repo.GetUser(param.UserID)
	if err != nil {
		return response, err
	}

	loan, err := uc.repo.GetLoan(param.LoanID)
	if err != nil {
		return response, err
	}

	response.LoanID = loan.ID

	if user.ID != loan.UserID {
		response.Message = "user id does not match"
		return response, nil
	}

	if loan.Outstanding == 0 {
		response.Message = "loan has been paid"
		return response, nil
	}

	currentTime := time.Now()
	if param.PaymentDate != "" {
		currentTime, _ = time.Parse("2006-01-02", param.PaymentDate)
	}

	billAmount := loan.BillAmount

	// truncate the time so we only compare the day.
	if currentTime.Truncate(24*time.Hour).Sub(loan.UpdateTime.Truncate(24*time.Hour)).Hours()/24 > 14 {
		isDelinquent = true
		err = uc.repo.UpdateUser(entity.User{ID: loan.UserID, IsDelinquent: isDelinquent})
		if err != nil {
			return response, err
		}

		multiplier := int(currentTime.Truncate(24*time.Hour).Sub(loan.UpdateTime.Truncate(24*time.Hour)).Hours() / 24 / 7)
		billAmount = billAmount * float32(multiplier)
		if billAmount > loan.Outstanding {
			billAmount = loan.Outstanding * (1 + interestRate)
		}
	}

	if billAmount != param.Amount {
		response.Message = fmt.Sprintf("amount is not equal to the bill amount. Expected: %0.2f, Actual: %0.2f", billAmount, param.Amount)
		return response, nil
	}

	loan.Outstanding -= param.Amount / (1 + interestRate)
	loan.UpdateTime = currentTime

	err = uc.repo.UpdateLoan(loan)
	if err != nil {
		return response, err
	}

	if isDelinquent && loan.Outstanding == 0 {
		remainingLoans, err := uc.GetActiveLoanByUserID(param.UserID)
		if err != nil {
			response.Message = "failed to get remaining loans for delinquent validation - contact administrator"
			return response, nil
		}

		if len(remainingLoans) == 0 {
			err = uc.repo.UpdateUser(entity.User{ID: param.UserID, IsDelinquent: false})
			if err != nil {
				response.Message = "failed to remove delinquent status - contact administrator"
				return response, nil
			}
		}
	}

	response.IsSuccess = true
	response.Message = "success to make a payment"
	return response, nil
}

// UpdateUser updates a user information based on the given id.
func (uc *UseCase) UpdateUser(param entity.User) error {
	return uc.repo.UpdateUser(param)
}
