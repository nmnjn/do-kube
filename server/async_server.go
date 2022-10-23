package main

import (
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
	failOnError(err, "Failed to declare queue.")

	msgs, err := ch.Consume(
		queue.Name, // queue name
		"",     // consumer
		false,  // auto-act
		false,  // exclusif
		false,  // no-local
		false,  // no-wait
		nil,
	)
	if err != nil { log.Println("Consumer registration failed with error: ", err)}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("New Message Recieved: %s", msg.Body)
			msg.Ack(true)
		}
	}()

	log.Println(" [*] Waiting for new messages... Press CTRL+C to exit.")
	<- forever
}
