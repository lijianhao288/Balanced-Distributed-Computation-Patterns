package main
import (
	"fmt"
	"github.com/rs/xid"
	"github.com/streadway/amqp"
	"log"
	"os"
)
func main() {
	links := []string{
		"http://google.com",
		"http://golang.org",
	}
	var numOfLink = 10
	for i := 0; i < numOfLink; i++ {
		fakeLink := fmt.Sprintf("http://web%d.com", i)
		links = append(links, fakeLink)
	}
	links2 := []string{
		"http://facebook.com",
		"http://amazon.com",
	}
	var numOfLink2 = 20
	for i := 11; i < numOfLink2; i++ {
		fakeLink := fmt.Sprintf("http://web%d.com", i)
		links2 = append(links2, fakeLink)
	}
	var linksToSend []string
	arg := os.Args[1]
	switch arg {
	case "1":
		linksToSend = links
	case "2":
		linksToSend = links2
	}
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
	err = chi.ExchangeDeclare("responseExchange", "direct", false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	q, err := chi.QueueDeclare("", false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi.QueueBind(q.Name, q.Name, "responseExchange", false, nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := chi.Consume(q.Name, "", false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	var jobCorr = make(map[string]string)
	for _, link := range linksToSend {
		var corrId = randomString()
		err := cho.Publish("jobExchange", "jobkey", false, false,
			amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: corrId,
				ReplyTo:       q.Name,
				Body:          []byte(link),
			})
		failOnError(err, "Failed to publish")
		fmt.Println("Published " + link)
		jobCorr[corrId] = link
	}
	go func() {
		for d := range msgs {
			if _, ok := jobCorr[d.CorrelationId]; ok {
				delete(jobCorr, d.CorrelationId)
				fmt.Println("Get result: " + string(d.Body))
			} else {
				fmt.Println("Got a not related msg")
			}
			d.Ack(false)
		}
	}()
	forever := make(chan bool)
	<-forever
}
func randomString() string {
	guid := xid.New()
	return guid.String()
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}