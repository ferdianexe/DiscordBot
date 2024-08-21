package music

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"

	"github.com/bwmarrin/discordgo"
	opus "layeh.com/gopus"
)

// Service is the Service layer of music
type Service struct {
	playMutex   sync.Mutex
	mapInstance map[string]*PlaylistStatus
}

// NewService creates a new use case.
func NewService() *Service {
	return &Service{
		playMutex:   sync.Mutex{},
		mapInstance: make(map[string]*PlaylistStatus),
	}
}

func (service *Service) PlayMusicLocally(voice *discordgo.VoiceConnection) error {
	ffmpegCmd := exec.Command("ffmpeg", "-i", "Test2.webm", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		fmt.Println("ffmpegCmd.StdoutPipe(): ", err)
		return err
	}

	buffer := bufio.NewReaderSize(ffmpegOut, 16384)

	if err := ffmpegCmd.Start(); err != nil {
		log.Printf("Error starting ffmpeg: %v", err)
		return err
	}
	defer ffmpegCmd.Wait()
	fmt.Println("ffmpeg started")
	send := make(chan []int16, 2)

	go service.sendPCM(voice, send)
	for {
		audioBuffer := make([]int16, 960*2)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}
		send <- audioBuffer
	}

	fmt.Println("ffmpeg stopped")
	return nil
}

func (service *Service) GetGuildIDPlaylistStatus(guildID string) *PlaylistStatus {
	playListStatus, isExists := service.mapInstance[guildID]
	if !isExists {
		playList := PlaylistStatus{
			PlayMutex:     sync.RWMutex{},
			PlayListMutex: sync.RWMutex{},
			IsPlaying:     false,
			PlaylistURL:   make([]string, 0),
		}
		service.mapInstance[guildID] = &playList
		return &playList
	}
	return playListStatus
}

func (service *Service) sendPCM(voice *discordgo.VoiceConnection, pcm <-chan []int16) {
	encoder, err := opus.NewEncoder(48000, 2, opus.Audio)
	if err != nil {
		fmt.Println("NewEncoder error,", err)
		return
	}
	for {
		receive, ok := <-pcm
		if !ok {
			fmt.Println("PCM channel closed")
			return
		}
		opus, err := encoder.Encode(receive, 960, 960*2*2)
		if err != nil {
			fmt.Println("Encoding error,", err)
			return
		}
		if !voice.Ready || voice.OpusSend == nil {
			fmt.Printf("Discordgo not ready for opus packets. %+v : %+v", voice.Ready, voice.OpusSend)
			return
		}
		voice.OpusSend <- opus
	}
}

func (service *Service) PlayMusicYoutube(voice *discordgo.VoiceConnection, url string, playlistStatus *PlaylistStatus) error {
	var err error

	playlistStatus.PlayMutex.Lock()
	playlistStatus.IsPlaying = true
	playlistStatus.PlayMutex.Unlock()

	defer func() {
		playlistStatus.PlayMutex.Lock()
		playlistStatus.IsPlaying = false
		playlistStatus.PlayMutex.Unlock()
	}()

	cmd := exec.Command("yt-dlp", "-f", "bestaudio[ext=webm]", "-o", "-", url)
	ffmpegCmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegCmd.Stdin, err = cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error PlayMusicYoutube: cmd.StdoutPipe(): ", err)
		return err
	}
	ffmpegOut, err := ffmpegCmd.StdoutPipe()
	if err != nil {
		fmt.Println("ffmpegCmd.StdoutPipe(): ", err)
		return err
	}

	buffer := bufio.NewReaderSize(ffmpegOut, 16384)

	if err := cmd.Start(); err != nil {
		log.Printf("PlayMusicYoutube Error starting cmd.Start(): %v", err)
		return err
	}

	if err := ffmpegCmd.Start(); err != nil {
		log.Printf("PlayMusicYoutube Error starting ffmpeg: %v", err)
		return err
	}

	defer func() {
		cmd.Wait()
		ffmpegCmd.Wait()
	}()
	fmt.Println("ffmpeg started")
	send := make(chan []int16, 2)

	go service.sendPCM(voice, send)
	for {
		audioBuffer := make([]int16, 960*2)
		err = binary.Read(buffer, binary.LittleEndian, &audioBuffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}
		send <- audioBuffer
	}

	fmt.Println("ffmpeg stopped")

	return nil
}
