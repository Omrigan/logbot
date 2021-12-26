package main

import (
	"fmt"
	bot2 "github.com/omrigan/logbot/pkg/bot"
	"github.com/omrigan/logbot/pkg/storage"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"
)

type State uint8

const (
	IdleState State = iota
	OptionExpected
	CommentExpected
	TsExpected
)

func Markup(columns int, values []string) *tgbotapi.ReplyKeyboardMarkup {
	var buttons [][]tgbotapi.KeyboardButton
	var row []tgbotapi.KeyboardButton
	var i int
out:
	for {
		for j := 0; j < columns; j++ {
			row = append(row, tgbotapi.NewKeyboardButton(values[i]))
			i++
			if i == len(values) {
				break out
			}
		}
		buttons = append(buttons, row)
		row = []tgbotapi.KeyboardButton{}
	}
	buttons = append(buttons, row)
	keyboard := tgbotapi.NewReplyKeyboard(buttons...)
	return &keyboard
}

func ItemsMarkup(types map[string]*RecordType) *tgbotapi.ReplyKeyboardMarkup {
	var values []string
	for val := range types {
		values = append(values, val)
	}
	sort.Strings(values)
	return Markup(2, append(values, "cancel"))
}

func OptionsMarkup(typ *RecordType) *tgbotapi.ReplyKeyboardMarkup {
	return Markup(3, append(typ.Options, "cancel"))
}

func TsMarkup() *tgbotapi.ReplyKeyboardMarkup {
	return Markup(2, []string{"done", "cancel"})
}

func CommentMarkup() *tgbotapi.ReplyKeyboardMarkup {
	return Markup(2, []string{"skip", "cancel"})
}

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Panic(err)
	}
	cfgM, _ := yaml.Marshal(cfg)
	log.Println(string(cfgM))
	bot := bot2.NewBot(cfg.Bot)

	var storages []storage.Storage
	if cfg.Influx != nil {
		storages = append(storages, storage.NewInfluxStorage(cfg.Influx))
	}
	if cfg.File != nil {
		storages = append(storages, storage.NewFileStorage(cfg.File))
	}
	store := storage.NewMetaStorage(storages...)


	state := IdleState
	record := &storage.Record{}
	var recType *RecordType
	for  {
		update, err := bot.Next()
		if err != nil {
			log.Println(err)
			continue
		}
		reply := func(txt string, markup *tgbotapi.ReplyKeyboardMarkup) {
			err := bot.Send(update.Chat.ID, txt, markup)
			if err != nil {
				log.Println(err)
			}
		}
		toTSExpected := func() {
			state = TsExpected
			reply(record.Time(), TsMarkup())
		}
		toCommentExpected := func() {
			state = CommentExpected
			reply("comment", CommentMarkup())
		}
		toOptionsExpected := func() {
			state = OptionExpected
			reply("option", OptionsMarkup(recType))
		}
		txt := strings.TrimSpace(update.Text)
		if txt == "cancel" {
			reply("ok", ItemsMarkup(cfg.RecordTypes))
			state = IdleState
			continue
		}
		switch state {
		case IdleState:
			record = &storage.Record{
				Item: txt,
				TS:   time.Now(),
			}
			var ok bool
			recType, ok = cfg.RecordTypes[txt]

			if !ok {
				reply(fmt.Sprintf("Item %s received, but no such type exists", txt), nil)
				toCommentExpected()
				continue
			}

			reply(fmt.Sprintf("Item %s received", txt), nil)
			if recType == nil {
				toTSExpected()
				continue
			}
			if len(recType.Options) != 0 {
				toOptionsExpected()
				continue
			}
			if recType.Comment {
				toCommentExpected()
				continue
			}

			toTSExpected()
		case OptionExpected:
			record.Param = txt
			reply(fmt.Sprintf("Option %s received", txt), nil)
			if recType != nil && recType.Comment {
				toCommentExpected()
				continue
			}
			toTSExpected()
		case CommentExpected:
			if txt == "skip" {
				toTSExpected()
				continue
			}
			record.Comment = txt
			reply(fmt.Sprintf("Comment %s received", txt), nil)
			toTSExpected()
		case TsExpected:
			if txt != "done" {
				ts, err := time.Parse(time.RFC822, txt)
				if err != nil {
					reply(err.Error(), nil)
					continue
				}
				record.TS = ts
			}
			err := store.Write(record)

			if err != nil {
				reply(fmt.Sprintf(err.Error()), nil)
				continue
			}

			state = IdleState
			reply(fmt.Sprintf("Recorded\n%s", record.YAML()), ItemsMarkup(cfg.RecordTypes))
		}
	}
}
