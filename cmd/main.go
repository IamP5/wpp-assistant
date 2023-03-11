package main

import (
	"github.com/IamP5/wpp-assistant/internal/web/server"
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
	webServer := server.MakeNewWebserver()
	webServer.Serve()
}
