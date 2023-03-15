package pkg

import (
	speech "cloud.google.com/go/speech/apiv1"
	"cloud.google.com/go/speech/apiv1/speechpb"
	"context"
	"log"
	"strings"
)

type SpeechToText struct {
	Client *speech.Client
}

func MakeSpeechToText() *SpeechToText {
	client, err := speech.NewClient(context.Background())
	if err != nil {
		log.Fatalf("[SpeechToText] Failed to create speech client: %v", err)
	}
	return &SpeechToText{
		Client: client,
	}
}

func (s *SpeechToText) RecognizeAudio(audio []byte) (string, error) {
	log.Println("[SpeechToText] - started RecognizeAudio method")

	resp, err := s.Client.Recognize(context.Background(), &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_FLAC,
			SampleRateHertz: 48000,
			LanguageCode:    "pt-BR",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audio},
		},
	})
	if err != nil {
		log.Printf("[SpeechToText] - RecognizeAudio error: %s", err.Error())
		return "", err
	}

	var transcriptions []string
	for _, result := range resp.Results {
		transcriptions = append(transcriptions, result.Alternatives[0].Transcript)
	}

	log.Println("[SpeechToText] - RecognizeAudio with success!")
	return strings.Join(transcriptions, " "), nil
}
