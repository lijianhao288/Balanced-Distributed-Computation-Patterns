package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
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
	err = cho.ExchangeDeclare("jobExchange", "direct",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	err = chi.ExchangeDeclare("responseExchange", "direct",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	q, err := chi.QueueDeclare("jobQueue",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi.QueueBind(q.Name, "jobkey", "jobExchange", false, nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := chi.Consume(q.Name, "",
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range msgs {
			var result = linkTest(string(d.Body))
			fmt.Println(result)
			var err = cho.Publish("responseExchange", d.ReplyTo,
				false, false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(result),
				})
			failOnError(err, "Failed to publish a message")
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for jobs")
	forever := make(chan bool)
	<-forever
}
func linkTest(link string) string {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	_, err := client.Get(link)
	if err != nil {
		resultString := link + " status: might down."
		return resultString
	}
	resultString := link + " status: up."
	return resultString
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
