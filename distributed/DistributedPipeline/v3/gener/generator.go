package main
import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"github.com/streadway/amqp"
)
type workerArg struct {
	WorkerType        string
	InputConfirmQueue string
	Function          string
}
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	err = ch.ExchangeDeclare("dispatch", "direct", 
		false, true, false, false, nil)
	failOnError(err, "Failed to declare an exchange")
	queueArg, err := ch.QueueDeclare("generatorQueue", 
		false, true, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(queueArg.Name, "generator", "dispatch", 
		false, nil)
	failOnError(err, "Failed to bind a queue")
	argMsgs, err := ch.Consume(queueArg.Name, "", 
		false, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")
	go func() {
		for d := range argMsgs {
			arg := workerArg{}
			json.Unmarshal(d.Body, &arg)
			fmt.Println(arg)
			cmd := exec.Command("cmd", "/C", "start", "../worker.exe",
				arg.WorkerType, arg.InputConfirmQueue, arg.Function)
			err = cmd.Run()
			failOnError(err, "Failed to generate worker")
			fmt.Println("generated one worker")
			d.Ack(false)
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