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
	err = ch.ExchangeDeclare("sharedQExchange", "direct",
		false, true, false, false, nil)
	printErrorAndExit(err, "Failed to declare an exchange")
	publishMsg(ch, "sharedQExchange", "one", "msg1")
	publishMsg(ch, "sharedQExchange", "one", "msg2")
	publishMsg(ch, "sharedQExchange", "one", "msg3")
}
func printErrorAndExit(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, ":", err)
	}
}
func publishMsg(c *amqp.Channel, ex string, key string, msg string) {
	body := msg
	err := (*c).Publish(ex, key, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	printErrorAndExit(err, "Failed to publish a message")
	fmt.Println("Sent: ", body)
}