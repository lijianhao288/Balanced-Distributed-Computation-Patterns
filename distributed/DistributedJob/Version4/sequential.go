package main
import (
	"fmt"
	"net/http"
	"time"
)
func main() {
	startTime := time.Now()
	links := []string{}
	var numOfLink = 50
	for i := 0; i < numOfLink; i++ {
		fakeLink := fmt.Sprintf("http://web%d.com", i)
		links = append(links, fakeLink)
	}
	for _, link := range links {
		result := linkTest(link)
		fmt.Println(result)
	}
	elapsed := time.Since(startTime)
	fmt.Printf("Time: %s \n", elapsed)
}
func linkTest(link string) string {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	_, err := client.Get(link)
	if err != nil {
		resultString := link + " status: might down"
		return resultString
	}
	resultString := link + " status: up. "
	return resultString
}
