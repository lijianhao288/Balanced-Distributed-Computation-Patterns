package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
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
	go func() {
		for x := 0; x < 11; x++ {
			body := strconv.Itoa(x)
			err = ch.Publish("pipelineExchangeV1", "KeyA",
				false, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			failOnError(err, "Failed to publish a message")
			fmt.Println("Worker 1 sent", body)
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	forever := make(chan bool)
	<-forever
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}