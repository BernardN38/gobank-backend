package main

import (
	"encoding/json"
	"errors"
	"github.com/bernardn38/gobank/broker-service/event"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type RequestPayload struct {
	Action   string          `json:"action"`
	Auth     AuthMessage     `json:"auth,omitempty"`
	Transfer TransferMessage `json:"transfer,omitempty"`
}

type AuthMessage struct {
	Name string      `json:"name"`
	Data AuthPayload `json:"data"`
}
type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type TransferMessage struct {
	Name string          `json:"name"`
	Data TransferPayload `json:"data"`
}

type TransferPayload struct {
	ToAccount   uuid.UUID `json:"to_account"`
	FromAccount uuid.UUID `json:"from_account"`
	Amount      int       `json:"amount"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authEvent(w, requestPayload.Auth)
	case "transfer":
		app.transferEvent(w, r, requestPayload.Transfer)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) transferEvent(w http.ResponseWriter, r *http.Request, l TransferMessage) {
	//token, err := r.Cookie("jwtToken")
	rawHeaderToken := r.Header.Get("Authorization")
	headerToken := strings.TrimPrefix(rawHeaderToken, "Bearer ")
	//if err != nil {
	//	app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
	//	return
	//}
	err := app.pushToQueue(l.Name, map[string]interface{}{"to_account": l.Data.ToAccount, "from_account": l.Data.FromAccount, "amount": l.Data.Amount, "token": headerToken})

	if err != nil {
		log.Println(err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via rabbit mq"
	app.writeJSON(w, http.StatusOK, payload)

}
func (app *Config) authEvent(w http.ResponseWriter, l AuthMessage) {
	err := app.pushToQueue(l.Name, map[string]interface{}{"username": l.Data.Username, "password": l.Data.Password})
	if err != nil {
		log.Println(err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via rabbit mq"
	app.writeJSON(w, http.StatusOK, payload)
}

type RabbitMessage struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

func (app *Config) pushToQueue(name string, data map[string]interface{}) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	payload := RabbitMessage{
		Name: name,
		Data: data,
	}
	j, _ := json.Marshal(&payload)

	err = emitter.Push(string(j), "auth")
	if err != nil {
		return err
	}
	return nil
}
