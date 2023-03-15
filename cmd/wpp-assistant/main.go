package main

import (
	"github.com/IamP5/wpp-assistant/internal/usecase"
	"github.com/IamP5/wpp-assistant/internal/web/server"
	"github.com/IamP5/wpp-assistant/pkg"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading main .env file: %s", err)
	}
}

func main() {
	openApi := pkg.MakeOpenAI()
	twilio := pkg.MakeTwilio()
	speechToText := pkg.MakeSpeechToText()

	msgToChatUsecase := usecase.MakeMessageToChat(openApi, twilio, speechToText)

	webServer := server.MakeNewWebserver(msgToChatUsecase)
	webServer.Serve()
}
