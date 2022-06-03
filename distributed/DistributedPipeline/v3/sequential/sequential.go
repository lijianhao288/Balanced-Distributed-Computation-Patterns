package main
import (
	"fmt"
	"strconv"
	"time"
)
func main() {
	startTime := time.Now()
	for x := 0; x < 10; x++ {
		num := x
		time.Sleep(1000 * time.Millisecond)
		for i := 0; i < 8; i++ {
			num = num + 2
			time.Sleep(1000 * time.Millisecond)
		}
		fmt.Println("number is:", strconv.Itoa(num))
		time.Sleep(1000 * time.Millisecond)
	}
	elapsed := time.Since(startTime)
	fmt.Println("Time:", elapsed)
}