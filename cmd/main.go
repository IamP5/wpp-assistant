package main

import (
	"github.com/IamP5/wpp-assistant/internal/web/server"
	"github.com/IamP5/wpp-assistant/pkg"
	"github.com/IamP5/wpp-assistant/usecase"
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
	msgToChatUsecase := usecase.MakeMessageToChat(openApi, twilio)

	webServer := server.MakeNewWebserver(msgToChatUsecase)
	webServer.Serve()
}
