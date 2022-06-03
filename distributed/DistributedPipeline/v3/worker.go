package main
import (
	"log"
	"os"
	"fmt"
	"github.com/streadway/amqp"
)
func main() {
	const exchangeName = "pipeExchangeV2"
	workerType := os.Args[1]
	inputConfirmQueue := os.Args[2]
	function := os.Args[3]
	conn1, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn1.Close()
	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn2.Close()
	conn3, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn3.Close()
	conn4, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn4.Close()
	if workerType == "startworker" {
		w := startWorker{Worker{conn1, conn2, conn3, conn4, 
			nil, nil, nil, exchangeName, inputConfirmQueue, function}}
		w.WaitNext()
		w.Work()
	} else if workerType == "midworker" {
		w := midWorker{Worker{conn1, conn2, conn3, conn4,
			nil, nil, nil, exchangeName, inputConfirmQueue, function}}
		w.ConnectPrevious()
		w.WaitNext()
		w.Work()
	} else {
		w := endWorker{Worker{conn1, conn2, conn3, conn4,
			nil, nil, nil, exchangeName, inputConfirmQueue, function}}
		w.ConnectPrevious()
		w.Work()
	}
	forever := make(chan bool)
	<-forever
}
type Worker struct{
	conn1             *amqp.Connection
	conn2             *amqp.Connection
	conn3             *amqp.Connection
	conn4             *amqp.Connection
	inputMsgs         <-chan amqp.Delivery
	confirmMsgs       <-chan amqp.Delivery
	localchos         *amqp.Channel
	exchangeName      string
	inputConfirmQueue string
	function          string
} 
func (w *Worker) ConnectPrevious() {
	//channel send confirm to inputExchange
	chis, err := w.conn1.Channel()
	failOnError(err, "Failed to open a channel")
	//channel receive from inputExchange
	chir, err := w.conn2.Channel()
	failOnError(err, "Failed to open a channel")
	//declare queue and bind to inputExchange
	queueForMsgs, err := chir.QueueDeclare("",
		false, false, true, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chir.QueueBind(queueForMsgs.Name, queueForMsgs.Name,
		w.exchangeName, false, nil)
	failOnError(err, "Failed to bind a queue")
	w.inputMsgs, err = chir.Consume(queueForMsgs.Name, "",
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	//send confirms with routekey as inputConfirmQueue
	confirmMessage := queueForMsgs.Name
	err = chis.Publish(w.exchangeName, w.inputConfirmQueue,
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(confirmMessage),
		})
	failOnError(err, "Failed to publish a message")
}
func (w *Worker) WaitNext() {
	//channel for receiving confirms from output channel
	chor, err := w.conn3.Channel()
	failOnError(err, "Failed to open a channel")
	//channel for sending to output channel
	chos, err := w.conn4.Channel()
	w.localchos = chos
	failOnError(err, "Failed to open a channel")
	err = chos.ExchangeDeclare(w.exchangeName, "direct",
		false, false, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	//declare confirm queue
	queueForConfirm, err := chor.QueueDeclare("",
		false, false, true, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chor.QueueBind(queueForConfirm.Name, queueForConfirm.Name,
		w.exchangeName, false, nil)
	failOnError(err, "Failed to bind a queue")
	w.confirmMsgs, err = chor.Consume(queueForConfirm.Name, "",
		false, true, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	//send and print out the confirm queue for next worker
	err = chos.Publish("dispatch", "resp",
		false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(queueForConfirm.Name),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Worker published input confirm queue",  
		queueForConfirm.Name)
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
