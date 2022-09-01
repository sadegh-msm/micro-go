package main

import (
	"context"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
)

type response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error while connecting to mongo", err)
		return nil, err
	}
	return conn, nil
}

// creating function for ease of working with json:
// reading a json, writing a json and pass an error with json

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxSize := 1048576 //one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must have only one json value")
	}
	return nil
}

func (app *Config) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJson(w http.ResponseWriter, err error, status ...int) error {
	var payload response
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(w, statusCode, payload)
}
