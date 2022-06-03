package main
import "fmt"
func f(x int) (a int, b int, s string) {
    a = 4
    b = x
    s = "apple"
    return 
}
func main() {
	fmt.Println(f(5))
}
