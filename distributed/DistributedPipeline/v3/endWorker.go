package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)
type endWorker struct {
	Worker
}
func (w *endWorker) Work() {
	go func() {
		for d := range w.inputMsgs {
			var numberString = string(d.Body)
			fmt.Println("End worker received", numberString)
			if numberString == "END" {
				c, err := w.conn3.Channel()
				failOnError(err, "Failed to create a channel")
				msg := "END"
				err = c.Publish("dispatch", "end",
					false, false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(msg),
				})
				failOnError(err, "Failed to publish a message")
				fmt.Println("End worker Published to dispather:", msg)
				fmt.Println("End worker finished")
				break
			}
			time.Sleep(1000 * time.Millisecond)
			d.Ack(false)
		}
	}()
}