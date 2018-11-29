package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/logger"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
)

const (
	logPath     = "run.log"
	defaultPort = "80"
)

func main() {

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer logger.Init("LoggerExample", true, true, lf).Close()

	// setup line bot
	channelSecret := os.Getenv("ChannelSecret")
	channelToken := os.Getenv("ChannelAccessToken")

	if channelSecret == "" || channelToken == "" {
		logger.Error("missing token")
		panic("missing token")
	}

	b, err := newBot(channelSecret, channelToken)
	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}

	// setup http server
	port := defaultPort
	if p := os.Getenv("PORT"); len(p) > 0 {
		port = p
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/callback", b.callback)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatalf("could not serve on port %s: %v", port, err)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Hello world")
}

type bot struct {
	*linebot.Client
}

func newBot(secretKey, token string) (*bot, error) {
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
				msg := linebot.NewTextMessage(message.ID + ":" + message.Text + " OK!")
				if _, err = b.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
					logger.Error(err)
				}
			}
		}
	}

}
