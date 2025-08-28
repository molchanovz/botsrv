package botsrv

import (
	"context"
	"fmt"

	"botsrv/pkg/db"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/vmkteam/embedlog"
)

const (
	startCommand = "/start"
)

type Config struct {
	Token string
}

type BotManager struct {
	embedlog.Logger
	dbo db.DB
	cr  db.CommonRepo
}

func NewBotManager(logger embedlog.Logger, dbo db.DB) *BotManager {
	return &BotManager{
		Logger: logger,
		dbo:    dbo,
		cr:     db.NewCommonRepo(dbo),
	}
}

func (bm *BotManager) RegisterBotHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, startCommand, bot.MatchTypePrefix, bm.startHandler)
}

func (bm *BotManager) DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Нажми %v", startCommand),
	})
	if err != nil {
		bm.Logger.Errorf("send message failed:%v", err)
		return
	}
}

func (bm *BotManager) startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Привет, %v!", update.Message.From.Username),
	})
	if err != nil {
		bm.Logger.Errorf("send message failed:%v", err)
		return
	}
}

func Ptr[T any](in T) *T {
	return &in
}

func Deref[T any](in *T) T {
	if in != nil {
		return *in
	}

	return *new(T) // return default value for type
}
