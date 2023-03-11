package server

import (
	"github.com/IamP5/wpp-assistant/internal/web/handler"
	"github.com/IamP5/wpp-assistant/pkg"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type Webserver struct{}

func MakeNewWebserver() *Webserver {
	return &Webserver{}
}

func (w *Webserver) Serve() {

	r := mux.NewRouter()
	n := negroni.New(negroni.NewLogger())
	t := pkg.MakeTwilio()
	o := pkg.MakeOpenAI()

	handler.MakeWppHandler(r, n, t, o)
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
