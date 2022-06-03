package main
import (
	"fmt"
	"github.com/rs/xid"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)
func main() {
	startTime := time.Now()
	links := []string{}
	var numOfLink = 50
	for i := 0; i < numOfLink; i++ {
		fakeLink := fmt.Sprintf("http://web%d.com", i)
		links = append(links, fakeLink)
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
	err = cho.ExchangeDeclare("jobExchange", "direct", false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	err = cho.ExchangeDeclare("responseExchange", "direct", false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	q, err := chi.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi.QueueBind(q.Name, q.Name, "responseExchange", false, nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := chi.Consume(q.Name, "responseConsumer",
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	var jobCorr = make(map[string]string)
	for _, link := range links {
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for d := range msgs {
			if _, ok := jobCorr[d.CorrelationId]; ok {
				delete(jobCorr, d.CorrelationId)
				fmt.Println("Get result: " + string(d.Body))
			} else {
				fmt.Println("Got a not related msg")
			}
			d.Ack(false)
			if len(jobCorr) == 0 {
				err = chi.Cancel("responseConsumer", false)
				failOnError(err, "Failed to cancel a consumer")
			}
		}
		wg.Done()
	}()
	wg.Wait()
	elapsed := time.Since(startTime)
	fmt.Println("Time: " + elapsed.String())

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