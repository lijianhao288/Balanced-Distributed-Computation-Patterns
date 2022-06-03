package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"strconv"
	"time"
)
type startWorker struct {
	Worker
}
func (w *startWorker) Work() {
	//block when waiting confirms
	nextWorkerQueue := <-w.confirmMsgs
	nextWorkerQueueName := string(nextWorkerQueue.Body)
	go func() {
		for x := 0; x < 10; x++ {
			msg := strconv.Itoa(x)
			err := w.localchos.Publish(w.exchangeName, 
				nextWorkerQueueName,false, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
				})
			failOnError(err, "Failed to publish a message")
			fmt.Println("Start worker published", msg)
			time.Sleep(1000 * time.Millisecond)
		}
		msg := "END"
		err := w.localchos.Publish(w.exchangeName, 
			nextWorkerQueueName,false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Start worker published", msg)
		fmt.Println("Start worker finished")
	}()
}