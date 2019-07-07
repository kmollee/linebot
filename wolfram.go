package main

import (
	"fmt"

	wolfram "github.com/Krognol/go-wolfram"
)

const (
	units   = wolfram.Metric
	timeout = 1000
)

func newWolframClient(ID string) (*wolfram.Client, error) {
	if ID == "" {
		return nil, fmt.Errorf("Could not create wolfram client: due ID is empty")
	}
	return &wolfram.Client{AppID: ID}, nil
}

func wolframQuery(text string) (string, error) {

	res, err := wolframClient.GetSpokentAnswerQuery(text, units, timeout)
	if err != nil {
		return "", err
	}
	return res, nil
}
