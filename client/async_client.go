package main

import (
	"os"
	"strings"
	"log"
	"github.com/streadway/amqp"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@167.172.7.118:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a communication channel")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"DEMO-QUEUE-NAME", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusif
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil { log.Println("Queue declaration failed with error: ", err)}

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",     // exchange
		queue.Name, // routing key - Queue Name
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish message.")
	log.Printf("[*] Message Sent: %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "HELLO"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
