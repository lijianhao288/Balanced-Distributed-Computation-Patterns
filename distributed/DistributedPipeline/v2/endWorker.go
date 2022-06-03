package main
import (
	"fmt"
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
				fmt.Println("End worker finished")
				break
			}
			time.Sleep(1000 * time.Millisecond)
			d.Ack(false)
		}
	}()
}