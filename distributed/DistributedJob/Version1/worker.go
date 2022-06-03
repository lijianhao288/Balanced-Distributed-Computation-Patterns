package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	chi, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer chi.Close()
	err = chi.ExchangeDeclare("jobExchange", "direct",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	queueIn, err := chi.QueueDeclare("jobQueue",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi.QueueBind(queueIn.Name, "jobkey", "jobExchange", false, nil)
	failOnError(err, "Failed to bind a queue")
	inputMsgs, err := chi.Consume(queueIn.Name, "",
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range inputMsgs {
			fmt.Println(linkTest(string(d.Body)))
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
		resultString := link + " status: might down"
		return resultString
	}
	resultString := link + " status: up. "
	return resultString
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}