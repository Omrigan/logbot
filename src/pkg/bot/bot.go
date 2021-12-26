package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotConfig struct {
	Token     string   `yaml:"token"`
	Whitelist []string `yaml:"whitelist"`
}
type Bot struct {
	cfg *BotConfig
	bot *tgbotapi.BotAPI
	chn tgbotapi.UpdatesChannel
}

func NewBot(cfg *BotConfig) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	return &Bot{
		cfg: cfg,
		bot: bot,
		chn: updates,
	}
}
func (m *Bot) Send(chatID int64, text string, markup *tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	msg.ParseMode = "Markdown"
	_, err := m.bot.Send(msg)
	return err
}

func (m *Bot) Next() (*tgbotapi.Message, error) {
	update := <-m.chn
	if update.Message == nil || update.Message.Text == ""{
		err := m.Send(update.Message.Chat.ID, "Please use only text messages", nil)
		if err != nil {
			return nil, err
		}
		return m.Next()
	}
	if len(m.cfg.Whitelist) == 0 {
		return update.Message, nil
	}
	for _, user := range m.cfg.Whitelist {
		if user == update.Message.From.UserName {
			return update.Message, nil
		}
	}
	err := m.Send(update.Message.Chat.ID, "You are not in whitelist", nil)
	if err != nil {
		return nil, err
	}
	return m.Next()
}
