package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/shubhamkumar96/go-microservices/broker-service/event"
	"github.com/shubhamkumar96/go-microservices/broker-service/logs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	case "logViaREST":
		// Log by directly calling logger-service via REST API
		app.logDataViaREST(w, requestPayload.Log)
	case "logViaRPC":
		// Log by directly calling logger-service via RPC
		app.logDataViaRPC(w, requestPayload.Log)
	case "logViaGRPC":
		// Log by directly calling logger-service via gRPC
		app.logDataViaGRPC(w, requestPayload.Log)
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

func (app *Config) logDataViaREST(w http.ResponseWriter, l LogPayload) {
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

// Define the type that remote RPC server expects
type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logDataViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	// pass on the exact RPC method name you want to call
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = result

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logDataViaGRPC(w http.ResponseWriter, l LogPayload) {
	// Connect to the gRPC server
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	// Create Client
	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the gRPC method
	respone, err := c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: l.Name,
			Data: l.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = respone.Result

	app.writeJSON(w, http.StatusAccepted, payload)
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
