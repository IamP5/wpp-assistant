package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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

func (t *Twilio) TranscriptAudioMessage(audioUrl string) (string, error) {
	log.Println("[Twilio TranscriptAudioMessage] - started method")

	transcriptionParams := url.Values{}
	transcriptionParams.Set("AudioUrl", audioUrl)
	transcriptionParams.Set("SpeechModel", "phone_call")

	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/"+os.Getenv("TWILIO_ACCOUNT_SID")+"/Transcriptions.json", strings.NewReader(transcriptionParams.Encode()))
	if err != nil {
		log.Println(err.Error())
	}

	req.SetBasicAuth(os.Getenv("TWILIO_ACCOUNT_SID"), os.Getenv("TWILIO_AUTH_TOKEN"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	// Parse the JSON response body into a map[string]interface{}
	var transcriptionResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&transcriptionResp)
	if err != nil {
		log.Println(err.Error())
	}

	// Do something with the transcription response, e.g. print the transcription SID:
	fmt.Println(transcriptionResp)
	return transcriptionResp["sid"].(string), nil
}
