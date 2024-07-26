package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(res http.ResponseWriter, req *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(res, req, &requestPayload)
	if err != nil {
		app.errorJSON(res, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(res, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(res, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// log successful authentication to logger-service
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(res, err)
		return
	}

	payLoad := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(res, http.StatusAccepted, payLoad)
}

func (app *Config) logRequest(name, data string) error {
	var logEntry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	logEntry.Name = name
	logEntry.Data = data

	// create json that we will send to the logger-service
	jsonData, _ := json.MarshalIndent(logEntry, "", "\t") // In production, use 'Marshal()' not 'MarshalIndent()'.

	// create a http request
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	// call the logger-service
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
