package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"tg-tw/telegram"
)

func main() {
	token := flag.String("token", "", "telegram api token")
	url := flag.String("api", "", "telegram api URL")
	flag.Parse()

	bot, err := telegram.New(
		*token,
		*url,
		slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})),
	)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
