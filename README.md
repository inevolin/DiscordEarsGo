# DiscordEarsGo

A speech-to-text bot for Discord written in GoLang. Can be useful for hearing impaired and deaf people. Try the bot for yourself on our Discord server: https://discord.gg/ApdTMG9

Note: This works only on a Linux x86 platform


## Setup

In your Discord Developers Bot settings, you need to enable these intents:

![image](https://github.com/inevolin/DiscordEarsBot/assets/53948000/6e926a75-a709-435a-b4f8-e9f8f0226856)

1. You may need to install opus `apt-get install opus-tools`
2. Clone this repository and cd inside
3. Create `.env` file to provide your Discord bot token: `TOKEN=...`

## Usage
- Run `VOSK_PATH=$(pwd)/vosk-linux-x86_64-0.3.45 LD_LIBRARY_PATH=$VOSK_PATH CGO_CPPFLAGS="-I $VOSK_PATH" CGO_LDFLAGS="-L $VOSK_PATH" go run ears.go`
- Or use `go build` to (with the above CLI prefix) and run the binary

## Usage
1. Enter one of your voice channels.
2. In one of your text channels type: `*join`, the bot will join the voice channel.
3. Everything said within that channel will be transcribed into text (as long as the bot is within the voice channel).
4. Type `*leave` to make the bot leave the voice channel.
