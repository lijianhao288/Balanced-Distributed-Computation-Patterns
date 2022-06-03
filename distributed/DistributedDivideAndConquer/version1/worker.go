package main
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"github.com/streadway/amqp"
)
type workerArg struct {
	OutputKey  string
	ConfirmKey string
}
func main() {
	//command line args
	outputKey := os.Args[1]
	confirmKey := os.Args[2]
	conn1, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn1.Close()
	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn2.Close()
	conn3, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn3.Close()
	choc, err := conn1.Channel()
	failOnError(err, "Failed to open a channel")
	chor, err := conn2.Channel()
	failOnError(err, "Failed to open a channel")
	chi, err := conn3.Channel()
	failOnError(err, "Failed to open a channel")
	exchangeName := "conquer"
	err = chi.ExchangeDeclare(exchangeName,"direct",
		false,true,false,false,nil)
	failOnError(err, "Failed to declare an exchange")
	queueIn, err := chi.QueueDeclare("",false,true,true,false,nil)
	failOnError(err, "Failed to declare a queue")
	err = chi.QueueBind(queueIn.Name,queueIn.Name,
		exchangeName,false,nil)
	failOnError(err, "Failed to bind a queue")
	inputMsgs, err := chi.Consume(queueIn.Name,"",
		false,false,false,false,nil)
	failOnError(err, "Failed to register a consumer")
	//send confirm
	msg := queueIn.Name
	err = choc.Publish(exchangeName,confirmKey,
		false,false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
	fmt.Println("Published confirm:", msg)
	//receive tasks
	task := <-inputMsgs
	list := []int{}
	json.Unmarshal(task.Body, &list)
	task.Ack(false)
	fmt.Println("Received:", list)
	if length := len(list); length <= 4 {
		//do the sort
		sort.Ints(list)
		//publish with outputKey
		listB, err := json.Marshal(list)
		failOnError(err, "Failed to encode")
		err = chor.Publish(exchangeName,outputKey,
			false,false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        listB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Published:", list)
	} else {
		//devide
		first := list[0]
		res := list[1:]
		leftList := Filter(res, 
			func(i int) bool { return i < first })
		rightList := Filter(res, 
			func(i int) bool { return i >= first })
		//new connections and channels
		conn4, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn4.Close()
		conn5, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn5.Close()
		conn6, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn6.Close()
		conn7, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn7.Close()
		conn8, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn8.Close()
		conn9, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn9.Close()
		chilc, err := conn4.Channel()
		failOnError(err, "Failed to open a channel")
		chilr, err := conn5.Channel()
		failOnError(err, "Failed to open a channel")
		chol, err := conn6.Channel()
		failOnError(err, "Failed to open a channel")
		chirc, err := conn7.Channel()
		failOnError(err, "Failed to open a channel")
		chirr, err := conn8.Channel()
		failOnError(err, "Failed to open a channel")
		chor, err := conn9.Channel()
		failOnError(err, "Failed to open a channel")
		////bind queue and send args
		//left
		queueLeftConfirm, err := chilc.QueueDeclare("",
			false,true,true,false,nil)
		failOnError(err, "Failed to declare a queue")
		err = chilc.QueueBind(queueLeftConfirm.Name, 
			queueLeftConfirm.Name, exchangeName,false,nil)
		failOnError(err, "Failed to bind a queue")
		leftConfirmMsgs, err := chilc.Consume(queueLeftConfirm.Name,
			"",false,true,false,false,nil)
		failOnError(err, "Failed to register a consumer")
		queueLeftResult, err := chilr.QueueDeclare("",
			false,true,true,false,nil)
		failOnError(err, "Failed to declare a queue")
		err = chilr.QueueBind(queueLeftResult.Name,
			queueLeftResult.Name,exchangeName, false,nil)
		failOnError(err, "Failed to bind a queue")
		leftResultMsgs, err := chilr.Consume(queueLeftResult.Name,
			"",false,false,false,false,nil)
		failOnError(err, "Failed to register a consumer")
		//right
		queueRightConfirm, err := chirc.QueueDeclare("",
			false,true,true,false, nil)
		failOnError(err, "Failed to declare a queue")
		err = chirc.QueueBind(queueRightConfirm.Name,
			queueRightConfirm.Name,exchangeName,false,nil)
		failOnError(err, "Failed to bind a queue")
		rightConfirmMsgs, err := chirc.Consume(
			queueRightConfirm.Name,"",false,true,false,false,nil)
		failOnError(err, "Failed to register a consumer")
		queueRightResult, err := chirr.QueueDeclare("",
			false,true,true,false,nil)
		failOnError(err, "Failed to declare a queue")
		err = chirr.QueueBind(queueRightResult.Name,
			queueRightResult.Name,exchangeName,false,nil)
		failOnError(err, "Failed to bind a queue")
		rightResultMsgs, err := chirr.Consume(
			queueRightResult.Name,"",false,false,false,false,nil)
		failOnError(err, "Failed to register a consumer")
		//send args
		argsLeftP := &workerArg{
			OutputKey:  queueLeftResult.Name,
			ConfirmKey: queueLeftConfirm.Name,
		}
		argsLeftB, err := json.Marshal(argsLeftP)
		failOnError(err, "Failed to encode")
		err = chol.Publish(exchangeName,"generator",
			false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        argsLeftB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Published", string(argsLeftB))
		argsRightP := &workerArg{
			OutputKey:  queueRightResult.Name,
			ConfirmKey: queueRightConfirm.Name,
		}
		argsRightB, err := json.Marshal(argsRightP)
		failOnError(err, "Failed to encode")
		err = chor.Publish(exchangeName,"generator",
			false,false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        argsRightB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Published", string(argsRightB))
		//receive confirm and send task
		leftConfirm := <-leftConfirmMsgs
		leftTargetKey := string(leftConfirm.Body)
		leftConfirm.Ack(false)
		rightConfirm := <-rightConfirmMsgs
		rightTargetKey := string(rightConfirm.Body)
		rightConfirm.Ack(false)
		leftListB, err := json.Marshal(leftList)
		failOnError(err, "Failed to encode")
		err = chol.Publish(exchangeName,leftTargetKey,
			false,false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        leftListB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Published:", leftList)
		rightListB, err := json.Marshal(rightList)
		failOnError(err, "Failed to encode")
		err = chor.Publish(exchangeName,rightTargetKey,
			false,false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        rightListB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Published:", rightList)
		//receive left and right result and publish final result
		leftResultMsg := <-leftResultMsgs
		listLeftResult := []int{}
		json.Unmarshal(leftResultMsg.Body, &listLeftResult)
		leftResultMsg.Ack(false)
		fmt.Println("Left result:", listLeftResult)
		rightResultMsg := <-rightResultMsgs
		listRightResult := []int{}
		json.Unmarshal(rightResultMsg.Body, &listRightResult)
		rightResultMsg.Ack(false)
		fmt.Println("Right result:", listRightResult)
		//final result
		finalResult := append(listLeftResult, first)
		finalResult = append(finalResult, listRightResult...)
		finalResultB, err := json.Marshal(finalResult)
		failOnError(err, "Failed to encode")
		err = chor.Publish(exchangeName,outputKey,
			false,false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        finalResultB,
			})
		failOnError(err, "Failed to publish a message")
		fmt.Println("Result:",finalResult)
	}
	forever := make(chan bool)
	<-forever
}
func Filter(s []int, fn func(int) bool) []int {
	var p []int // == nil
	for _, i := range s {
		if fn(i) {
			p = append(p, i)
		}
	}
	return p
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
