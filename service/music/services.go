package music

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	opus "layeh.com/gopus"
)

// Service is the Service layer of music
type Service struct {
}

// NewService creates a new use case.
func NewService() *Service {
	return &Service{}
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
