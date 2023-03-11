package handler

import (
	"github.com/IamP5/wpp-assistant/internal/web/handler/dto"
	"github.com/IamP5/wpp-assistant/pkg"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func MakeWppHandler(r *mux.Router, n *negroni.Negroni, t *pkg.Twilio, o *pkg.OpenAI) {
	r.Handle("/wpp/receive", n.With(
		negroni.Wrap(receiveWppMessage(t, o)),
	)).Methods("POST", "OPTIONS")
}

func receiveWppMessage(t *pkg.Twilio, o *pkg.OpenAI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		request := &dto.TwilioWebhook{
			Id:                r.Form.Get("SmsSid"),
			Status:            r.Form.Get("SmsStatus"),
			MessageStatus:     r.Form.Get("MessageStatus"),
			ChannelToAddress:  r.Form.Get("ChannelToAddress"),
			To:                r.Form.Get("To"),
			ChannelPrefix:     r.Form.Get("ChannelPrefix"),
			MessageSid:        r.Form.Get("MessageSid"),
			AccountSid:        r.Form.Get("AccountSid"),
			From:              r.Form.Get("From"),
			ApiVersion:        r.Form.Get("ApiVersion"),
			ChannelInstallSid: r.Form.Get("ChannelInstallSid"),
		}

		message := t.GetMessageBySid(request.MessageSid)
		chatData, chatErr := o.CompleteChat(*message.Body)

		if chatErr != nil {
			log.Printf("Error in Open AI")
		}

		err = t.SendMessage(request.To, request.From, chatData)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}
