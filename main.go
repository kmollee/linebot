package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/google/logger"

	// "github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	witai "github.com/wit-ai/wit-go"
)

type config struct {
	Port              string `env:"PORT" envDefault:"8000"`
	lineChannelSecret string `env:"ChannelSecret"`
	linehannelToken   string `env:"ChannelAccessToken"`
	witToekn          string `env:"WitToken"`
}

const (
	logFilename         = "run.log"
	confidenceThreshold = 0.5
)

var (
	witClient *witai.Client
)

func main() {
	lf, err := newLoggerFile()
	if err != nil {
		log.Fatal(err)
	}
	defer lf.Close()
	defer logger.Init("LoggerExample", true, true, lf).Close()

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	// setup line bot
	b, err := newBot(cfg.lineChannelSecret, cfg.linehannelToken)
	if err != nil {
		logger.Error(err.Error())
		panic(err.Error())
	}

	// setup wit bot
	witClient = witai.NewClient(cfg.witToekn)

	http.HandleFunc("/", b.index)
	http.HandleFunc("/callback", b.callback)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		logger.Fatalf("could not serve on port %s: %v", cfg.Port, err)
	}

}

func newLoggerFile() (*os.File, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, errors.Wrap(err, "could not execute dir path")
	}
	exPath := filepath.Dir(ex)
	logPath := exPath + logFilename

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open log file")
	}
	return lf, nil
}
