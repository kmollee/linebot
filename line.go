package main

import (
	"fmt"
	"net/http"

	"github.com/google/logger"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
)

type bot struct {
	*linebot.Client
}

func (b *bot) index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Hello world")
}

func newBot(secretKey, token string) (*bot, error) {
	if secretKey == "" || token == "" {
		return nil, fmt.Errorf("secretKey: %s or token:%s are empty", secretKey, token)
	}
	b, err := linebot.New(secretKey, token)
	if err != nil {
		return nil, errors.Wrap(err, "could not create bot")
	}
	return &bot{b}, nil
}

func (b *bot) callback(w http.ResponseWriter, r *http.Request) {
	events, err := b.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var msg *linebot.TextMessage

				msgType, err := witParse(message.Text)
				if err != nil || msgType == "" {
					// wit not able recognize message
					msg = linebot.NewTextMessage(message.ID + ":" + message.Text + "I don't know what you mean?")

				} else {

					switch msgType {
					case "greetings":
						msg = linebot.NewTextMessage(message.ID + ":" + "Hi!")
					case "wikipedia_search_query":
						res, err := wolframQuery(message.Text)
						if err != nil {
							msg = linebot.NewTextMessage(message.ID + ":" + message.Text + "I don't know what you mean?")
						} else {
							msg = linebot.NewTextMessage(res)
						}

					default:
						msg = linebot.NewTextMessage(message.ID + ":" + message.Text + "I don't know what you mean?")

					}
				}

				if _, err = b.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
					logger.Error(err)
				}
			}
		}
	}

}
