package main
import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"github.com/streadway/amqp"
)
type workerArg struct {
	OutputKey  string
	ConfirmKey string
}
func main() {
	startTime := time.Now()
	//list := []int{3, 4, 7, 2, 5, 7, 8, 4, 6, 8, 6, 3, 66, 432, 
		//63, 6, 7, 8, 4, 65, 34, 4, 36}
	list := []int{5,7,3,4,1,9,6,2,8}
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
	chic, err := conn2.Channel()
	failOnError(err, "Failed to open a channel")
	chir, err := conn3.Channel()
	failOnError(err, "Failed to open a channel")
	exchangeName := "conquer"
	err = cho.ExchangeDeclare(exchangeName,"direct",
		false,true,false,false,nil)
	failOnError(err, "Failed to declare an exchange")
	queueForConfirm, err := chic.QueueDeclare("",
		false,true,true,false,nil)
	failOnError(err, "Failed to declare a queue")
	err = chic.QueueBind(queueForConfirm.Name,
	queueForConfirm.Name,exchangeName,false,nil)
	failOnError(err, "Failed to bind a queue")
	confirmMsgs, err := chic.Consume(queueForConfirm.Name,
		"",false,true,false,false,nil)
	failOnError(err, "Failed to register a consumer")
	queueForResult, err := chir.QueueDeclare("",
		false,true,true,false,nil)
	failOnError(err, "Failed to declare a queue")
	err = chir.QueueBind(queueForResult.Name,
		queueForResult.Name,exchangeName,false,nil)
	failOnError(err, "Failed to bind a queue")
	resultMsgs, err := chir.Consume(queueForResult.Name,"",
		false,false,false,false,nil)
	failOnError(err, "Failed to register a consumer")
	//send args to generate worker
	argsP := &workerArg{
		OutputKey:  queueForResult.Name,
		ConfirmKey: queueForConfirm.Name,
	}
	argsB, err := json.Marshal(argsP)
	failOnError(err, "Failed to encode")
	err = cho.Publish(exchangeName,"generator",
		false,false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        argsB,
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("starter Published", string(argsB))
	confirm := <-confirmMsgs
	targetKey := string(confirm.Body)
	confirm.Ack(false)
	listB, err := json.Marshal(list)
	failOnError(err, "Failed to encode")
	msg := listB
	err = cho.Publish(exchangeName,targetKey,
		false,false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Published:", list)
	result := []int{}
	resultMsg := <-resultMsgs
	json.Unmarshal(resultMsg.Body, &result)
	fmt.Println("result:", result)
	resultMsg.Ack(false)
	elapsed := time.Since(startTime)
	fmt.Println("Time: ", elapsed)
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}