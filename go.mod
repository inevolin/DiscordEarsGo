module discord/earsbot

go 1.21.4

require (
	discord/vosk v0.0.0
	github.com/bwmarrin/discordgo v0.27.1
	github.com/joho/godotenv v1.5.1
	gopkg.in/hraban/opus.v2 v2.0.0-20230925203106-0188a62cb302
)

replace discord/vosk => ./internal/vosk/

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
)
