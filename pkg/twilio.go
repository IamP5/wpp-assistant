package pkg

import (
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
)

type Twilio struct {
	Client *twilio.RestClient
}

func MakeTwilio() *Twilio {
	return &Twilio{
		Client: twilio.NewRestClient(),
	}
}

func (t *Twilio) GetMessageBySid(Sid string) *openapi.ApiV2010Message {
	msg, err := t.Client.Api.FetchMessage(Sid, &openapi.FetchMessageParams{})
	if err != nil {
		log.Printf("[Twilio] Error getting message %s: %s", Sid, err.Error())
	}

	return msg
}

func (t *Twilio) SendMessage(from string, to string, message string) error {
	params := &openapi.CreateMessageParams{
		From: &from,
		To:   &to,
		Body: &message,
	}

	_, err := t.Client.Api.CreateMessage(params)

	if err != nil {
		log.Printf("[Twilio] Failed to sending message from %s to %s: %s", from, to, err)
	}

	return err
}
