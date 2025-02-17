package main

import (
	"os"

	"github.com/np123-glitch/vice-bot/bot"
)

func main() {
	bot.Run(os.Getenv("DISCORD_TOKEN")) // Pass token here
}
