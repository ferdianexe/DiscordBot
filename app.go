package main

import (
	"ferdianexe/DiscordBot/handler"
	"ferdianexe/DiscordBot/infra/config"
	"ferdianexe/DiscordBot/service/music"
	"ferdianexe/DiscordBot/usecase"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type Application struct {
	handler *handler.Handler
}

func NewApplication(handler *handler.Handler) *Application {
	return &Application{handler: handler}
}

func (app *Application) routes() {
	http.HandleFunc("/ping", app.handler.Ping)
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal("Error message")
	}
}

func main() {
	// db, err := initDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// dbRepo := sqlite.NewRepository(db)

	fileReader := config.NewConfigFileReader("")
	configProvider := config.NewService(fileReader)

	// err := dbRepo.SetupDatabase()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Connected to database")

	// create a session
	discord, err := discordgo.New("Bot " + configProvider.GetConfig().BotID)
	checkNilErr(err)
	musicSvc := music.NewService()
	uc := usecase.NewUseCase(discord, musicSvc)
	h := handler.NewHandler(uc)
	app := NewApplication(h)
	app.routes()

	// add a event handler
	discord.AddHandler(h.IncomingMessageWrapper)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func initDB() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", "./db.db")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }
