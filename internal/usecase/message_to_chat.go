package usecase

import (
	"fmt"
	"github.com/IamP5/wpp-assistant/internal/web/handler/dto"
	"github.com/IamP5/wpp-assistant/pkg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

type MessageToChat struct {
	OpenAI *pkg.OpenAI
	Twilio *pkg.Twilio
	Speech *pkg.SpeechToText
}

type MessageToChatInput struct {
	To      string
	From    string
	Message *dto.TwilioWebhook
}

func MakeMessageToChat(openAi *pkg.OpenAI, twilio *pkg.Twilio, speech *pkg.SpeechToText) *MessageToChat {
	return &MessageToChat{
		OpenAI: openAi,
		Twilio: twilio,
		Speech: speech,
	}
}

func (u *MessageToChat) Execute(input *MessageToChatInput) error {
	log.Println("[MsgToChatUsecase] - usecase executed")

	message := u.Twilio.GetMessageBySid(input.Message.MessageSid)
	numMedia, err := strconv.Atoi(*message.NumMedia)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	if numMedia > 0 {
		err := u.proceedToAudio(input)
		if err != nil {
			log.Printf("[MsgToChatUsecase] - usecase failed")
			return err
		}
		return nil
	}

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

func (u *MessageToChat) proceedToAudio(input *MessageToChatInput) error {
	resp, err := http.Get(input.Message.MediaUrl)
	if err != nil {
		log.Printf("[MsgToChatUsecase - ProceedToAudio] - Failed to download audio file: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	flac, err := decodeAudio(resp.Body)
	if err != nil {
		log.Printf("[MsgToChatUsecase - ProceedToAudio] - Failed to decode audio to flac: %s", err.Error())
		return err
	}

	text, err := u.Speech.RecognizeAudio(flac)
	if err != nil {
		log.Printf("[MsgToChatUsecase - ProceedToAudio] - Failed to transcribe audio: %s", err.Error())
		return err
	}

	chat, err := u.OpenAI.CompleteChat(text)
	if err != nil {
		log.Printf("[MsgToChatUsecase - ProceedToAudio] - Failed to complete chat: %s", err.Error())
		return err
	}

	err = u.Twilio.SendMessage(input.From, input.To, chat)
	if err != nil {
		log.Printf("[MsgToChatUsecase - ProceedToAudio] - Failed to send message: %s", err.Error())
		return err
	}

	return nil
}

func decodeAudio(audio io.Reader) ([]byte, error) {
	// Convert OGG to FLAC
	// Save the OGG audio to a temporary file
	tmpOggFile, err := ioutil.TempFile("", "audio_*.ogg")
	if err != nil {
		return nil, err
	}
	//defer os.Remove(tmpOggFile.Name())

	_, err = io.Copy(tmpOggFile, audio)
	if err != nil {
		log.Printf("[MsgToChatUsecase - transcribeAudio] - Failed to copy audio to temporary file: %s", err.Error())
		return nil, err
	}

	// Convert OGG to FLAC using ffmpeg
	tmpFlacFile, err := ioutil.TempFile("", "audio_*.flac")
	if err != nil {
		log.Printf("[MsgToChatUsecase - transcribeAudio] - Failed to create temporary file: %s", err.Error())
		return nil, err
	}
	//defer os.Remove(tmpFlacFile.Name())

	cmd := exec.Command("ffmpeg", "-i", tmpOggFile.Name(), "-y", "-f", "flac", tmpFlacFile.Name())
	err = cmd.Run()
	if err != nil {
		log.Printf("[MsgToChatUsecase - transcribeAudio] - Failed to convert OGG to FLAC: %s", err.Error())
		return nil, fmt.Errorf("failed to convert OGG to FLAC: %v", err)
	}

	flacData, err := ioutil.ReadAll(tmpFlacFile)
	if err != nil {
		log.Printf("[MsgToChatUsecase - transcribeAudio] - Failed to read temporary file: %s", err.Error())
		return nil, err
	}

	return flacData, nil
}
