package main

import (
	"fmt"
	"github.com/bernardn38/gobank/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	//try to connect to rabitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()
	//start listening for messages
	log.Println("Listening for and consuming rabbitMq messages")
	//create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	//watch the queue and consume events
	err = consumer.Listen([]string{"auth"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
		} else {
			log.Println("Connected to RabbbitMQ")
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backoff = backoff + (time.Second * 5)
		log.Println("backing off")
		time.Sleep(backoff)
	}
	return connection, nil
}
