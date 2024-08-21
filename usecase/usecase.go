package usecase

import (
	"ferdianexe/DiscordBot/service/music"
	"fmt"
	"log"
	"strings"
	"sync"

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

// musicServiceProvider is the interface for the repository.
type musicServiceProvider interface {
	PlayMusicLocally(voice *discordgo.VoiceConnection) error

	PlayMusicYoutube(voice *discordgo.VoiceConnection, url string, playlistStatus *music.PlaylistStatus) error

	GetGuildIDPlaylistStatus(guildID string) *music.PlaylistStatus
}

// UseCase is the usecase entity
type Usecase struct {
	repo  repoProvider
	music musicServiceProvider
	lock  sync.Mutex
}

// NewUseCase creates a new use case.
func NewUseCase(dc repoProvider, music musicServiceProvider) *Usecase {
	return &Usecase{
		repo:  dc,
		lock:  sync.Mutex{},
		music: music,
	}
}

func (usecase *Usecase) AddQueueSong(message *discordgo.MessageCreate, url string) error {
	guildPlaylist := usecase.music.GetGuildIDPlaylistStatus(message.GuildID)

	if !guildPlaylist.IsPlaying {
		usecase.repo.ChannelMessageSend(message.ChannelID, "There is no song playing right now, use playt command instead !")
		return nil
	}
	if len(url) == 0 {
		url = strings.Trim(message.Content, "!queue")
	}

	guildPlaylist.PlayListMutex.Lock()
	guildPlaylist.PlaylistURL = append(guildPlaylist.PlaylistURL, url)
	guildPlaylist.PlayListMutex.Unlock()

	guildPlaylist.PlayListMutex.RLock()
	totalLength := len(guildPlaylist.PlaylistURL)
	guildPlaylist.PlayListMutex.RUnlock()

	usecase.repo.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Adding to queue play %s %d song(s) remaining", url, totalLength))
	return nil
}

func (usecase *Usecase) PlayMusic(message *discordgo.MessageCreate, voice *discordgo.VoiceState) error {
	if voice == nil {
		usecase.repo.ChannelMessageSend(message.ChannelID, "You need to join the voice channel first !")
		return nil
	}
	usecase.repo.ChannelMessageSend(message.ChannelID, "Music!!!")
	voiceConnection, err := usecase.repo.ChannelVoiceJoin(message.GuildID, voice.ChannelID, false, false)
	if err != nil {
		log.Printf("usecase.repo.ChannelVoiceJoin return error (%v) - PlayMusic", err)
		return err
	}

	voiceConnection.Speaking(true)
	defer voiceConnection.Speaking(false)

	usecase.music.PlayMusicLocally(voiceConnection)
	return nil
}

func (usecase *Usecase) PlayMusicYoutube(message *discordgo.MessageCreate, voice *discordgo.VoiceState) error {
	if voice == nil {
		usecase.repo.ChannelMessageSend(message.ChannelID, "You need to join the voice channel first !")
		return nil
	}
	voiceConnection, err := usecase.repo.ChannelVoiceJoin(message.GuildID, voice.ChannelID, false, false)
	if err != nil {
		log.Printf("usecase.repo.ChannelVoiceJoin return error (%v) - PlayMusicYoutube", err)
		return err
	}

	guildPlaylist := usecase.music.GetGuildIDPlaylistStatus(message.GuildID)

	url := strings.Trim(message.Content, "!playt")

	guildPlaylist.PlayMutex.RLock()
	if guildPlaylist.IsPlaying {
		guildPlaylist.PlayMutex.RUnlock()
		usecase.AddQueueSong(message, url)
		return nil
	}
	guildPlaylist.PlayMutex.RUnlock()

	voiceConnection.Speaking(true)
	defer voiceConnection.Speaking(false)

	guildPlaylist.PlaylistURL = append(guildPlaylist.PlaylistURL, url)
	for {
		guildPlaylist.PlayListMutex.RLock()
		if len(guildPlaylist.PlaylistURL) == 0 {
			guildPlaylist.PlayListMutex.RUnlock()
			break
		}
		playedUrl := guildPlaylist.PlaylistURL[0]
		guildPlaylist.PlayListMutex.RUnlock()
		guildPlaylist.PlayListMutex.Lock()
		guildPlaylist.PlaylistURL = guildPlaylist.PlaylistURL[1:]
		guildPlaylist.PlayListMutex.Unlock()
		usecase.repo.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Playing music from YouTube... %s", playedUrl))
		usecase.music.PlayMusicYoutube(voiceConnection, playedUrl, guildPlaylist)
	}
	return nil
}
