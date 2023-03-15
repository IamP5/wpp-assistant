package server

import (
	"github.com/IamP5/wpp-assistant/internal/usecase"
	"github.com/IamP5/wpp-assistant/internal/web/handler"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type Webserver struct {
	messageToChatUsecase *usecase.MessageToChat
}

func MakeNewWebserver(messageToChatUsecase *usecase.MessageToChat) *Webserver {
	return &Webserver{
		messageToChatUsecase: messageToChatUsecase,
	}
}

func (w *Webserver) Serve() {

	r := mux.NewRouter()
	n := negroni.New(negroni.NewLogger())

	handler.MakeWppHandler(r, n, w.messageToChatUsecase)
	http.Handle("/", r)

	server := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              ":8080",
		Handler:           http.DefaultServeMux,
		ErrorLog:          log.New(os.Stderr, "log: ", log.Lshortfile),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
