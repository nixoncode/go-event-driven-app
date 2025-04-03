package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

func sendWelcomeEmail(userEmail string) {
	// Simulate sending a welcome email
	fmt.Printf("Sending welcome email to %s\n", userEmail)
	// Simulate a delay
	time.Sleep(2 * time.Second)
	// Log the email sent
	fmt.Printf("Welcome email sent to %s\n", userEmail)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"registration_queue", // name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var event map[string]string
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				fmt.Printf("Error unmarshalling message: %s\n", err)
				continue
			}
			sendWelcomeEmail(event["email"])
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}
