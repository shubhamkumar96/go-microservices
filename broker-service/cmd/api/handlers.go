package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shubhamkumar96/go-microservices/broker-service/event"
)

// Common request payload struct for all the microservices
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// payload for auth-service
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// payload for logger-service
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// payload for mail-service
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payLoad := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payLoad)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// Log by directly calling logger-service
		app.logData(w, requestPayload.Log)
	case "logViaRMQ":
		// Log by pushing event to RabbitMQ, from where the event will be consumed
		// by listener-service, and then listener-service calls logger-service
		app.logDataViaRabbitMQ(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json that we will send to the auth-service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// create a http request
	request, err := http.NewRequest("POST", "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	// call the auth-service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// get the correct status-code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth-service"))
		return
	}

	// read the response body
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payLoad jsonResponse
	payLoad.Error = false
	payLoad.Message = "Authenticated!"
	payLoad.Data = jsonFromService.Data

	// write to response
	app.writeJSON(w, http.StatusAccepted, payLoad)
}

func (app *Config) logData(w http.ResponseWriter, l LogPayload) {
	// create json that we will send to the logger-service
	jsonData, _ := json.MarshalIndent(l, "", "\t")

	// create a http request
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	// call the logger-service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// get the correct status-code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling auth-service"))
		return
	}

	var payLoad jsonResponse
	payLoad.Error = false
	payLoad.Message = "Logged!"

	// write to response
	app.writeJSON(w, http.StatusAccepted, payLoad)
}

func (app *Config) sendMail(w http.ResponseWriter, m MailPayload) {
	// create json that we will send to the mail-service
	jsonData, _ := json.MarshalIndent(m, "", "\t")

	// create a http request
	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	// call the mail-service
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// get the correct status-code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail-service"))
		return
	}

	var payLoad jsonResponse
	payLoad.Error = false
	payLoad.Message = "Message Sent to " + m.To

	// write to response
	app.writeJSON(w, http.StatusAccepted, payLoad)
}

func (app *Config) logDataViaRabbitMQ(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)
}

// Utility function used to push messages to Queue.
func (app *Config) pushToQueue(name, msg string) error {
	producer, err := event.NewProducer(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = producer.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}
