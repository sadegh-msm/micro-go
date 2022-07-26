package main

import (
	"broker/cmd/api/auth"
	"broker/cmd/api/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type BrokerRequest struct {
	Action   string          `json:"action"`
	Auth     AuthPayload     `json:"auth,omitempty"`
	Log      LogPayload      `json:"log,omitempty"`
	Shortner shortnerPayload `json:"shortner,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type shortnerPayload struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	ExpireTime  time.Duration `json:"expireTime"`
}

func (app *Config) isAlive(w http.ResponseWriter, r *http.Request) {
	res := response{
		Message: "im still alive",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

// Broker all services will connect to broker and communicate with each other with broker service
// this service will take BrokerRequest and decide foe all request where to go
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	res := response{
		Error:   false,
		Message: "Broker is active",
	}

	_ = app.writeJson(w, http.StatusOK, res)
}

func (app *Config) HandleAll(w http.ResponseWriter, r *http.Request) {
	var req BrokerRequest

	err := app.readJson(w, r, &req)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "auth":
		app.authenticateWithGrpc(w, r)

	case "log":
		app.logEventWithGrpc(w, r)

	case "shortner":
		app.shortnerURL(w, req.Shortner)

	default:
		app.errorJson(w, errors.New("invalid action"), http.StatusBadRequest)
	}
}

// Authenticate Authenticats by http/rest
func (app *Config) Authenticate(w http.ResponseWriter, data AuthPayload) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")

	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("user Unauthorized"))
		return
	} else if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error when authorizing"))
		return
	}

	var responseFromService response

	err = json.NewDecoder(res.Body).Decode(&responseFromService)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	if responseFromService.Error {
		app.errorJson(w, err, http.StatusUnauthorized)
		return
	}

	payload := response{
		Error:   false,
		Message: "auth completed",
		Data:    responseFromService,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}

// LogItem logs items by http/rest
func (app *Config) LogItem(w http.ResponseWriter, data LogPayload) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")

	req, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJson(w, err)
		return
	}

	payload := response{
		Error:   false,
		Message: "logged data",
	}
	app.writeJson(w, http.StatusAccepted, payload)
}

func (app *Config) shortnerURL(w http.ResponseWriter, data shortnerPayload) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")

	req, err := http.NewRequest("POST", "http://urlshortner-service/api/v1", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		app.errorJson(w, err)
		return
	}

	payload := response{
		Error:   false,
		Message: "url shortend",
		Data:    res,
	}
	app.writeJson(w, http.StatusAccepted, payload)
}

// logs items by Grpc and
func (app *Config) logEventWithGrpc(w http.ResponseWriter, r *http.Request) {
	var req BrokerRequest

	err := app.readJson(w, r, &req)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: req.Log.Name,
			Data: req.Log.Data,
		},
	})
	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := response{
		Error:   false,
		Message: "logged by grpc",
	}

	app.writeJson(w, http.StatusAccepted, res)
}

func (app *Config) authenticateWithGrpc(w http.ResponseWriter, r *http.Request) {
	var req BrokerRequest

	err := app.readJson(w, r, &req)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	conn, err := grpc.Dial("authentication-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer conn.Close()

	c := auth.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.Authenticate(ctx, &auth.AuthRequest{
		AuthEntry: &auth.Authenticate{
			Email:    req.Auth.Email,
			Password: req.Auth.Password,
		},
	})
	if err != nil {
		app.errorJson(w, err)
		return
	}

	res := response{
		Error:   false,
		Message: "authenticated by grpc",
	}

	app.writeJson(w, http.StatusAccepted, res)
}
