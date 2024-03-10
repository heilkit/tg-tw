package telegram

import (
	"fmt"
	"github.com/heilkit/tg"
	"github.com/heilkit/tg/scheduler"
	"log/slog"
	"os"
	"path"
	"tg-tw/twitter"
)

type Bot struct {
	tg  *tg.Bot
	tw  *twitter.API
	log *slog.Logger
}

func New(token string, api string, log *slog.Logger) (*Bot, error) {
	tgBot, err := tg.NewBot(tg.Settings{
		URL:       api,
		Token:     token,
		Local:     tg.LocalMoving(),
		Scheduler: scheduler.ExtraConservative(),
		Retries:   4,
		Logger:    tg.LoggerSlog(log),
	})
	if err != nil {
		log.Error("while creating bot", "token", token, "api", api)
		return nil, err
	}

	return &Bot{
		tg:  tgBot,
		tw:  twitter.New(),
		log: log,
	}, nil
}

func (bot *Bot) Start() {
	bot.tg.Handle(tg.OnText, func(ctx tg.Context) error {
		url, err := twitter.Vx(ctx.Text())
		if err != nil {
			return ctx.Reply("url parsing failed")
		}
		files, dir, post, err := bot.tw.DownloadTempVx(url)
		defer func(path string) {
			if err := os.RemoveAll(path); err != nil {
				bot.log.Error("while removing dir", "err", err, "path", path)
			}
		}(dir)
		for _, file := range files {
			filename := fmt.Sprintf("@%s_%s", post.UserScreenName, path.Base(file))
			if err := ctx.Reply(&tg.Document{FileName: filename, File: tg.FromDisk(file)}); err != nil {
				bot.log.Error("while uploading file", "err", err, "file", file)
				return err
			}
		}
		return nil
	})

	bot.tg.Handle("/start", func(ctx tg.Context) error {
		return bot.tg.React(ctx.Message(), tg.ReactionPray)
	})

	bot.tg.Start()
}
