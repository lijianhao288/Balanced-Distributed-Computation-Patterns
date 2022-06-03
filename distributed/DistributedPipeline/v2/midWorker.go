package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)
type midWorker struct {
	Worker
}
func (w *midWorker) Work() {
	//block when waiting confirms
	nextWorkerQueue := <-w.confirmMsgs
	nextWorkerQueueName := string(nextWorkerQueue.Body)
	go func() {
		for d := range w.inputMsgs {
			var numberString = string(d.Body)
			if numberString == "END" {
				msg := "END"
				err := w.localchos.Publish(w.exchangeName, 
					nextWorkerQueueName,false, false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(msg),
					})
				failOnError(err, "Failed to publish a message")
				fmt.Println("Mid worker received", numberString,
					"Published", msg)
				fmt.Println("Mid worker finished")
				break
			}
			number, _ := strconv.Atoi(numberString)
			var changedNumber int
			switch w.function {
			case "+2":
				changedNumber = number + 2
			case "*2":
				changedNumber = number * 2
			}
			msg := strconv.Itoa(changedNumber)
			err := w.localchos.Publish(w.exchangeName, 
				nextWorkerQueueName,false, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
				})
			failOnError(err, "Failed to publish a message")
			fmt.Println("Mid worker received", numberString,
				"Published", msg)
			time.Sleep(1000 * time.Millisecond)
			d.Ack(false)
		}
	}()
}