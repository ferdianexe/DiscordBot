package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jery1402/billing-engine/handler"
	"github.com/jery1402/billing-engine/repository/sqlite"
	"github.com/jery1402/billing-engine/usecase"
)

type Application struct {
	handler *handler.Handler
}

func NewApplication(handler *handler.Handler) *Application {
	return &Application{handler: handler}
}

func (app *Application) routes() {
	http.HandleFunc("/init_db", app.handler.CreateDatabase)
	http.HandleFunc("/create_user", app.handler.CreateUser)
	http.HandleFunc("/get_payment_schedule", app.handler.GetNextPayment)
	http.HandleFunc("/get_outstanding", app.handler.GetOutstanding)
	http.HandleFunc("/loan_list", app.handler.GetLoanList)
	http.HandleFunc("/user_list", app.handler.GetUserList)
	http.HandleFunc("/is_delinquent", app.handler.GetUserDelinquentStatus)
	http.HandleFunc("/make_loan", app.handler.MakeLoan)
	http.HandleFunc("/make_payment", app.handler.MakePayment)
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	dbRepo := sqlite.NewRepository(db)

	uc := usecase.NewUseCase(dbRepo)
	h := handler.NewHandler(uc)

	app := NewApplication(h)

	err = dbRepo.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	app.routes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./billing-engine.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
