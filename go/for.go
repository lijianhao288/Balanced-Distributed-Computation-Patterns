package main
import "fmt"
func main() {
    product := 1
    for i := 1; i < 5; i++ {
        product *= i
    }
    fmt.Println(product)
}