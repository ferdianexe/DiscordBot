package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// usecaseProvider is the interface for the usecase.
type usecaseProvider interface {
	PlayMusic(message *discordgo.MessageCreate, voice *discordgo.VoiceState) error
}

// Handler is the handler enti.
type Handler struct {
	usecase usecaseProvider
}

// NewHandler creates a new handler.
func NewHandler(usecase usecaseProvider) *Handler {
	return &Handler{usecase: usecase}
}

// Ping is Ping
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, http.StatusCreated, GeneralResponse{Message: "success"})
}

// IncomingMessageWrapper all message discord should enter this function
func (h *Handler) IncomingMessageWrapper(discord *discordgo.Session, message *discordgo.MessageCreate) {
	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	case strings.Contains(message.Content, "!bye"):
		discord.ChannelMessageSend(message.ChannelID, "Good Bye👋")
	case strings.Contains(message.Content, "!play"):
		// get current user voice state
		voiceState, _ := discord.State.VoiceState(message.GuildID, message.Author.ID)
		h.usecase.PlayMusic(message, voiceState)
	}
}

// writeResponse writes a HTTP response.
func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonResponse, _ := json.Marshal(data)
	w.Write(jsonResponse)
}

// GeneralResponse is the general response entity.
type GeneralResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
