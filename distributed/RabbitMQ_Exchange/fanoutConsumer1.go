package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	printErrorAndExit(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	printErrorAndExit(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare("fanoutExchange", "fanout",
		false, true, false, false, nil)
	printErrorAndExit(err, "Failed to declare an exchange")
	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	printErrorAndExit(err, "Failed to declare a queue")
	err = ch.QueueBind(q.Name, "anykey2", "fanoutExchange", false, nil)
	printErrorAndExit(err, "Failed to bind a queue")
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	printErrorAndExit(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			bodyString := string(d.Body)
			fmt.Println("Received:", bodyString)
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for msgs")
	forever := make(chan bool)
	<-forever
}
func printErrorAndExit(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ":", err)
	}
}