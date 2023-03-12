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
	log.Println("[Twilio GetMessageBySid] - started method")

	msg, err := t.Client.Api.FetchMessage(Sid, &openapi.FetchMessageParams{})
	if err != nil {
		log.Printf("[Twilio] Error getting message %s: %s", Sid, err.Error())
	}

	log.Println("[Twilio GetMessageBySid] - received message with success")
	return msg
}

func (t *Twilio) SendMessage(from string, to string, message string) error {
	log.Println("[Twilio SendMessage] - started method")

	params := &openapi.CreateMessageParams{
		From: &from,
		To:   &to,
		Body: &message,
	}

	log.Println("[Twilio SendMessage] - sending message")
	_, err := t.Client.Api.CreateMessage(params)

	if err != nil {
		log.Printf("[Twilio] Failed to sending message from %s to %s: %s", from, to, err)
	}

	log.Println("[Twilio SendMessage] - message sent with success!")
	return err
}
