package main
import (
	"fmt"
	"log"
	"time"
	"encoding/json"
	"github.com/streadway/amqp"
)
type workerArg struct {
	WorkerType        string
	InputConfirmQueue string
	Function          string
}
func main() {
	startTime := time.Now()
	conn1, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn1.Close()
	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn2.Close()
	conn3, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn3.Close()
	cho, err := conn1.Channel()
	failOnError(err, "Failed to open a channel")
	chi1, err := conn2.Channel()
	failOnError(err, "Failed to open a channel")
	chi2, err := conn3.Channel()
	failOnError(err, "Failed to open a channel")
	err = cho.ExchangeDeclare("dispatch", "direct", 
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	//resp
	queueResp, err := chi1.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi1.QueueBind(queueResp.Name, "resp", "dispatch", false, nil)
	failOnError(err, "Failed to bind a queue")
	respMsgs, err := chi1.Consume(queueResp.Name, "", 
		false, true, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	//end
	queueEnd, err := chi2.QueueDeclare("", false, true, true, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = chi2.QueueBind(queueEnd.Name, "end", "dispatch", false, nil)
	failOnError(err, "Failed to bind a queue")
	endMsgs, err := chi2.Consume(queueEnd.Name, "", 
		false, true, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	////send args
	//start worker index 1
	numWorkers := 10
	indexWorker := 1
	//start worker
	startArgP := &workerArg{
		WorkerType:        "startworker",
		InputConfirmQueue: "null",
		Function:          "null",
	}
	startArgB, err := json.Marshal(startArgP)
	failOnError(err, "Failed to encode")
	err = cho.Publish("dispatch", "generator", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        startArgB,
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Organizer Published:")
	fmt.Println(string(startArgB))
	//block consume from queueResp
	//mid and end workers
	for d := range respMsgs {
		indexWorker++
		//2~numWorkers
		workerType := "midworker"
		function := "+2"
		if indexWorker == numWorkers {
			workerType = "endworker"
			function = "null"
		}
		//publish args for mid worker
		midArgP := &workerArg{
			WorkerType:        workerType,
			InputConfirmQueue: string(d.Body),
			Function:          function,
		}
		midArgB, _ := json.Marshal(midArgP)
		err = cho.Publish("dispatch", "generator", false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        midArgB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Organizer Published:")
		fmt.Println(string(midArgB))
		d.Ack(false)
		//limit number
		if indexWorker == numWorkers {
			break
		}
	}
	//wait for the end message
	for d := range endMsgs {
		d.Ack(false)
		break
	}
	//print the time
	elapsed := time.Since(startTime)
	fmt.Println("Time:", elapsed)
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}