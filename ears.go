package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"encoding/binary"

	vosk "discord/vosk"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"gopkg.in/hraban/opus.v2"
)

var model, _ = vosk.NewModel("./vosk_models/en")
var stt, _ = vosk.NewRecognizer(model, 48000)

func handleVoice(c chan *discordgo.Packet) {
	decoder, _ := opus.NewDecoder(48000, 2)

	buffer := new(bytes.Buffer)

	frame := make([]int16, 48000)
	for {
		select {
		case s, ok := <-c:
			if !ok {
				break
			}
			decoder.Decode(s.Opus, frame)
			binary.Write(buffer, binary.LittleEndian, frame)
			log.Println(frame)

			// if stt.AcceptWaveform(out) != 0 {
			// 	fmt.Println(rec.Result())
			// }
		}
	}
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, "*join") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				// join channel
				log.Println(g.ID, vs.ChannelID)
				v, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, true, false)
				if err != nil {
					fmt.Println("failed to join voice channel:", err)
					return
				}

				// TODO: when leave command is given, close objects:
				// go func() {
				// 	time.Sleep(10 * time.Second)
				// 	close(v.OpusRecv)
				// 	v.Close()
				// }()

				go handleVoice(v.OpusRecv)
			}
		}

	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "!airhorn")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	Token := os.Getenv("TOKEN")

	client, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	client.AddHandler(ready)
	defer client.Close()

	// We only really care about receiving voice and message updates.
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
