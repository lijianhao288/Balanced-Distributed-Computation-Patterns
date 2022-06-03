package main
import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)
const MaxWorker = 50
var jobQueue = make(chan string, 100)
var maxGo uint64
var wg sync.WaitGroup
func main() {
	go goroutineCounter()
	start := time.Now()
	wg.Add(1)
	go linkSender()
	wg.Add(1)
	go workerCreator()
	wg.Wait()
	fmt.Println("Max goroutine number: ", atomic.LoadUint64(&maxGo))
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
		time.Sleep(200 * time.Millisecond)
	}
}
func linkSender() {
	defer wg.Done()
	links := []string{}
	var numOfLink = 1000
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
	//index := 0
	workers := []worker{}
	for i := 0; i < MaxWorker; i++ {
		worker := worker{i, make(chan string), make(chan string)}
		worker.Start()
		workers = append(workers, worker)
	}
	for link := range jobQueue {
		//selectedWorker := roundRobin(workers, &index)
		selectedWorker := randomSelect(workers)
		selectedWorker.jobChannel <- link
	}
	for _, w := range workers {
		w.quitChannel <- "q"
	}

}
func roundRobin(l []worker, currentIndex *int) worker {
	selected := l[(*currentIndex)]
	if (*currentIndex) >= len(l)-1 {
		(*currentIndex) = 0
	} else {
		(*currentIndex)++
	}
	return selected
}
func randomSelect(l []worker) worker {
	lens := len(l)
	index := rand.Intn(lens)
	selected := l[index]
	return selected
}
type worker struct {
	id          int
	jobChannel  chan string
	quitChannel chan string
}
func (w worker) Start() {
	wg.Add(1)
	go func() {
		defer wg.Done()
	L:
		for {
			select {
			case link := <-w.jobChannel:
				fmt.Println("Worker ",
					strconv.Itoa(w.id), ": ",
					linkTest(link))
			case <-w.quitChannel:
				fmt.Println("Worker ",
					strconv.Itoa(w.id), "Quit")
				break L
			}
		}
	}()
}
func linkTest(link string) string {
	time.Sleep(500 * time.Millisecond)
	if rand.Intn(2) == 1 {
		return link + ": Good"
	} else {
		return link + ": Bad"
	}
}