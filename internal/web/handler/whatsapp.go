package handler

import (
	"fmt"
	"github.com/IamP5/wpp-assistant/internal/web/handler/dto"
	"github.com/IamP5/wpp-assistant/usecase"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

func MakeWppHandler(r *mux.Router, n *negroni.Negroni, messageToChatUsecase *usecase.MessageToChat) {
	r.Handle("/wpp/receive", n.With(
		negroni.Wrap(receiveWppMessage(messageToChatUsecase)),
	)).Methods("POST", "OPTIONS")
}

func receiveWppMessage(msgToChatUsecase *usecase.MessageToChat) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		fmt.Println(r.PostForm)

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
			MediaContentType:  r.Form.Get("MediaContentType0"),
			MediaUrl:          r.Form.Get("MediaUrl0"),
		}

		input := &usecase.MessageToChatInput{
			To:      request.From,
			From:    request.To,
			Message: request,
		}

		err = msgToChatUsecase.Execute(input)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}
