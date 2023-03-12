package usecase

import (
	"github.com/IamP5/wpp-assistant/internal/web/handler/dto"
	"github.com/IamP5/wpp-assistant/pkg"
	"log"
)

type MessageToChat struct {
	OpenAI *pkg.OpenAI
	Twilio *pkg.Twilio
}

type MessageToChatInput struct {
	To      string
	From    string
	Message *dto.TwilioWebhook
}

func MakeMessageToChat(openAi *pkg.OpenAI, twilio *pkg.Twilio) *MessageToChat {
	return &MessageToChat{
		OpenAI: openAi,
		Twilio: twilio,
	}
}

func (u *MessageToChat) Execute(input *MessageToChatInput) error {
	log.Println("[MsgToChatUsecase] - usecase executed")

	message := u.Twilio.GetMessageBySid(input.Message.MessageSid)

	chat, err := u.OpenAI.CompleteChat(*message.Body)
	if err != nil {
		log.Printf("[MsgToChatUsecase] - usecase failed)")
		return err
	}

	err = u.Twilio.SendMessage(input.From, input.To, chat)
	if err != nil {
		log.Printf("[MsgToChatUsecase] - usecase failed")
		return err
	}

	log.Println("[MsgToChatUsecase] - usecase executed with success!")
	return nil
}
