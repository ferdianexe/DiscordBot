package usecase

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// repoProvider is the interface for the repository.
type repoProvider interface {
	ChannelMessageSend(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	// ------------------------------------------------------------------------------------------------
	// Functions specific to Discord Guilds
	// ------------------------------------------------------------------------------------------------

	// Guild returns a Guild structure of a specific Guild.
	// guildID   : The ID of a Guild
	Guild(guildID string, options ...discordgo.RequestOption) (st *discordgo.Guild, err error)

	// ChannelVoiceJoin joins the session user to a voice channel.
	//
	//    gID     : Guild ID of the channel to join.
	//    cID     : Channel ID of the channel to join.
	//    mute    : If true, you will be set to muted upon joining.
	//    deaf    : If true, you will be set to deafened upon joining.
	ChannelVoiceJoin(gID, cID string, mute, deaf bool) (voice *discordgo.VoiceConnection, err error)
}

// UseCase is the usecase entity
type Usecase struct {
	repo repoProvider
}

// NewUseCase creates a new use case.
func NewUseCase(dc repoProvider) *Usecase {
	return &Usecase{
		repo: dc,
	}
}

func (usecase *Usecase) PlayMusic(message *discordgo.MessageCreate, voice *discordgo.VoiceState) error {
	if voice == nil {
		usecase.repo.ChannelMessageSend(message.ChannelID, "You need to join the voice channel first !")
		return nil
	}
	usecase.repo.ChannelMessageSend(message.ChannelID, "Music!!!")
	_, err := usecase.repo.ChannelVoiceJoin(message.GuildID, voice.ChannelID, false, false)
	if err != nil {
		log.Printf("usecase.repo.ChannelVoiceJoin return error (%v) - PlayMusic", err)
		return err
	}
	return nil
}
