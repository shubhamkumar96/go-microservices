package main

import (
	"net/http"
)

func (app *Config) SendMail(res http.ResponseWriter, req *http.Request) {
	var requestPayload struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	err := app.readJSON(res, req, &requestPayload)
	if err != nil {
		app.errorJSON(res, err, http.StatusBadRequest)
		return
	}

	// Create Message object
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(res, err, http.StatusBadRequest)
		return
	}

	// Create and Send Json Response
	payload := jsonResponse{
		Error:   false,
		Message: "Message Sent to " + requestPayload.To,
	}

	app.writeJSON(res, http.StatusAccepted, payload)
}
