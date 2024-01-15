package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	vosk "github.com/alphacep/vosk-api/go"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"layeh.com/gopus"
)

type STTResult struct {
	Text string `json:"text"`
}

var model, _ = vosk.NewModel("./vosk_models/en") // path to Vosk model (default: english-very-small)
var stt, _ = vosk.NewRecognizer(model, 48000)    // 48kHz
var speakers, _ = gopus.NewDecoder(48000, 1)     // 48kHz and mono-channel
var voiceConnection *discordgo.VoiceConnection

func handleVoice(client *discordgo.Session, channelID string, user string, c chan *discordgo.Packet) {
	var buffer = new(bytes.Buffer)
	for {
		select {
		case s, ok := <-c:
			if !ok {
				break
			}
			if buffer == nil {
				buffer = new(bytes.Buffer)
			}
			packet, _ := speakers.Decode(s.Opus, 960, false) // frameSize is 960 (20ms)
			pcm := new(bytes.Buffer)
			binary.Write(pcm, binary.LittleEndian, packet)
			buffer.Write(pcm.Bytes())
			stt.AcceptWaveform(pcm.Bytes())

			var dur float32 = (float32(len(buffer.Bytes())) / 48000 / 2) // duration of audio

			// When silence packet detected, send result (skip audio shorter than 500ms)
			if dur > 0.5 && len(s.Opus) == 3 && s.Opus[0] == 248 && s.Opus[1] == 255 && s.Opus[2] == 254 {
				log.Println("dur", dur)
				var result STTResult
				json.Unmarshal([]byte(stt.FinalResult()), &result)
				if len(result.Text) > 0 {
					log.Println(fmt.Sprintf("%s: %s", user, result.Text))
					// process the transcription result:
					client.ChannelMessageSend(channelID, fmt.Sprintf("%s: %s", user, result.Text)) // send as text to channel
				}
				buffer.Reset()
			}
		}
	}
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// handle text commands:

	if strings.HasPrefix(m.Content, "*join") {
		// Find the guild for that channel.
		guild, err := s.State.Guild(m.GuildID)
		if err != nil {
			fmt.Println("Could not find guild", err)
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Author.ID {
				// join channel where user is at
				v, err := s.ChannelVoiceJoin(guild.ID, vs.ChannelID, true, false)
				if err != nil {
					fmt.Println("failed to join voice channel:", err)
					return
				}
				voiceConnection = v
				handleVoice(s, m.ChannelID, m.Author.Username, v.OpusRecv)
			}
		}

	} else if strings.HasPrefix(m.Content, "*leave") {
		voiceConnection.Disconnect()
		close(voiceConnection.OpusRecv)
		voiceConnection.Close()
	}
}

func main() {
	// get TOKEN from .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	Token := os.Getenv("TOKEN")

	// init discordgo client
	client, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	defer client.Close()

	// Intents for receiving voice, message updates, read members and send new messages
	client.Identify.Intents = discordgo.IntentsAll
	client.AddHandler(handleMessage)
	client.StateEnabled = true
	err = client.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	fmt.Println("Bot is running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
