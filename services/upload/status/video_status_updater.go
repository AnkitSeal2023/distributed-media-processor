package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func main() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"media_events", // name
		"fanout",       // type
		true,           // durability
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	fmt.Println("Exchange declared successfully")

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durability
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	fmt.Println("Queue declared successfully")

	// for _, s := range os.Args[1:] {
	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "media_events", "")
	err = ch.QueueBind(
		q.Name,         // queue name
		"",             // routing key
		"media_events", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	fmt.Printf("Queue %s bound to exchange %s with routing key %s\n", q.Name, "media_events", "")

	// }

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Println(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
