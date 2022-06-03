package main
import "fmt"
func main() {
	fmt.Println(multiplicate(3, 4))
	fmt.Println(multipAndAdd(3, 4))
}
func multiplicate(x, y int) int {
	return x * y
}
func multipAndAdd(x, y int) (int, int) {
	return x * y, x + y
}