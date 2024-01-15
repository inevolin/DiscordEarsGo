# DiscordEarsGo

This works only on a Linux x86 platform

## Installation
- clone this repository and cd inside
- create `.env` file to provide your Discord bot token: `TOKEN=...`
- run `VOSK_PATH=$(pwd)/vosk-linux-x86_64-0.3.45 LD_LIBRARY_PATH=$VOSK_PATH CGO_CPPFLAGS="-I $VOSK_PATH" CGO_LDFLAGS="-L $VOSK_PATH" go run ears.go`
