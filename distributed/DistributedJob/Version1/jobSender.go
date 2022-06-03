package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)
func main() {
	links := []string{
		"http://google.com",
		"http://facebook.com",
		"http://stackoverflow.com",
		"http://golang.org",
		"http://amazon.com",
	}
	var numOfLink = 10
	for i := 0; i < numOfLink; i++ {
		fakeLink := fmt.Sprintf("http://web%d.com", i)
		links = append(links, fakeLink)
	}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	cho, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer cho.Close()
	err = cho.ExchangeDeclare("jobExchange", "direct",
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	for _, link := range links {
		err := cho.Publish("jobExchange", "jobkey", false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(link),
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Published " + link)
	}
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}