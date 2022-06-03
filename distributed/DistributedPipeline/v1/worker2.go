package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)
func main() {
	conn1, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn1.Close()
	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn2.Close()
	cho, err := conn1.Channel()
	failOnError(err, "Failed to open a channel")
	defer cho.Close()
	chi, err := conn2.Channel()
	failOnError(err, "Failed to open a channel")
	defer chi.Close()
	err = chi.ExchangeDeclare("pipelineExchangeV1", "direct",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	q, err := chi.QueueDeclare("Worker2Queue",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi.QueueBind(q.Name, "KeyA", "pipelineExchangeV1",
		false, nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := chi.Consume(q.Name, "",
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			var numberString = string(d.Body)
			number, _ := strconv.Atoi(numberString)
			var doubleNumber = number * 2
			body := strconv.Itoa(doubleNumber)
			err = cho.Publish("pipelineExchangeV1", "KeyB",
				false, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			fmt.Println("Worker 2 received ", numberString,
				"sent", body)
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