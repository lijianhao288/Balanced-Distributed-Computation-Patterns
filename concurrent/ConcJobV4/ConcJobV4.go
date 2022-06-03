package main
import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"log"
	"time"
	"github.com/Jeffail/tunny"
)
var jobQueue = make(chan string, 100)
var maxGo uint64
var wg sync.WaitGroup
var pool *tunny.Pool
func main() {
	go goroutineCounter()
	pool = tunny.NewFunc(20, 
	func(payload interface{}) interface{} {
		var result string
		s, ok := payload.(string)
		if(!ok){
			log.Fatalln("type assertion fail")
		}
		result = linkTest(s)
		defer wg.Done()
		return result
	})
	defer pool.Close()
	start := time.Now()
	wg.Add(1)
	go linkSender()
	wg.Add(1)
	go workerCreator()
	wg.Wait()
	fmt.Println("Max goroutine number: ", 
	atomic.LoadUint64(&maxGo))
	duration := time.Since(start)
	fmt.Println("Time: ", duration)
}
func goroutineCounter() {
	for {
		n := runtime.NumGoroutine()
		u := uint64(n)
		if u > maxGo {
			atomic.StoreUint64(&maxGo, u)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
func linkSender() {
	defer wg.Done()
	links := []string{}
	var numOfLink = 100
	for i := 0; i < numOfLink; i++ {
		fakeLink := fmt.Sprintf("http://web%d.com", i)
		links = append(links, fakeLink)
	}
	for _, link := range links {
		jobQueue <- link
	}
	close(jobQueue)
}
func workerCreator() {
	defer wg.Done()
	for link := range jobQueue {
		wg.Add(1)
		result := pool.Process(link)
		fmt.Println(result)
	}
}
func linkTest(link string) string {
	time.Sleep(50 * time.Millisecond)
	if rand.Intn(2) == 1 {
		return link + ": Good"
	} else {
		return link + ": Bad"
	}
}
