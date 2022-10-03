package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}
type AuthPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	for _, s := range topics {
		err := ch.QueueBind(q.Name, s, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload, &d)
		}
	}()
	fmt.Printf("Waiting fo message on [Exchange, Queue] [logs_topic, %s]", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload, d *amqp.Delivery) {
	switch payload.Name {
	case "auth":
		err := authEvent(payload)
		if err != nil {
			log.Println(err)
		}
		err = d.Ack(false)
		if err != nil {
			log.Println(err)
		}
	case "transfer":
		err := transferEvent(payload)
		if err != nil {
			log.Println(err)
		}
		err = d.Ack(false)
		if err != nil {
			log.Println(err)
		}
	default:
		log.Println("No name to process")
	}
}

func authEvent(entry Payload) error {
	jsonData, _ := json.Marshal(entry.Data)
	authServiceURL := "http://auth-service/login"

	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return err
	}

	log.Println(response.Cookies())
	return nil
}

type TransferPayload struct {
	ToAccount   uuid.UUID `json:"to_account"`
	FromAccount uuid.UUID `json:"from_account"`
	Amount      int       `json:"amount"`
}

func transferEvent(entry Payload) error {
	payload := TransferPayload{
		ToAccount:   uuid.Must(uuid.Parse(entry.Data["to_account"].(string))),
		FromAccount: uuid.Must(uuid.Parse(entry.Data["from_account"].(string))),
		Amount:      int(entry.Data["amount"].(float64)),
	}
	jsonData, _ := json.Marshal(payload)
	transferServiceURL := "http://transaction-service/transfers"

	request, err := http.NewRequest("POST", transferServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", entry.Data["token"]))
	client := &http.Client{}
	response, err := client.Do(request)
	log.Println("posting transfer")
	if err != nil {
		return err
	}
	defer response.Body.Close()
	log.Println(response.StatusCode)
	if response.StatusCode != http.StatusOK {
		log.Println(err)
		return err
	}

	log.Println(response.Cookies())
	return nil
}
