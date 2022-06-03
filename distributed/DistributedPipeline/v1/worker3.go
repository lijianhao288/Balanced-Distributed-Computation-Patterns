package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare("pipelineExchangeV1", "direct",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	q, err := ch.QueueDeclare("Worker3Queue",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(q.Name, "KeyB", "pipelineExchangeV1",
		false, nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := ch.Consume(q.Name, "",
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			var numberString = string(d.Body)
			fmt.Println("Worker 3 received", numberString)
			time.Sleep(1000 * time.Millisecond)
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for jobs")
	forever := make(chan bool)
	<-forever
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}