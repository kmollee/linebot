package main

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	witai "github.com/wit-ai/wit-go"
)

// Outcome represents the outcome portion of a Wit message
type witEntity struct {
	Text       string  `json:"_text"`
	Intent     string  `json:"intent"`
	IntentID   string  `json:"intent_id"`
	Confidence float32 `json:"confidence"`
}

func witParse(msg string) (string, error) {
	resp, err := witClient.Parse(&witai.MessageRequest{
		Query: msg,
	})

	if err != nil {
		return "", errors.Wrap(err, "could not get wit.ai response")
	}

	var e witEntity

	var topCate string // most possible category
	var topCateConfident float32

	for cate, entitys := range resp.Entities {
		for _, entity := range entitys.([]interface{}) {
			err := mapstructure.Decode(entity, &e)
			if err != nil {
				return "", err
			}
			if e.Confidence > confidenceThreshold && e.Confidence > topCateConfident {
				topCate = cate
				topCateConfident = e.Confidence
			}
		}
	}

	return topCate, nil
}
