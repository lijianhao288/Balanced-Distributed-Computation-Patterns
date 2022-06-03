package main
import "fmt"
func main() {
    lunch:= "apple"
    switch lunch {
    case "banana":
        fmt.Println("My lunch was a banana.")
    case "apple":
        fmt.Println("My lunch was an apple.")
    case "pear":
        fmt.Println("My lunch was a pear.")
    case "tomato":
        fmt.Println("My lunch was a tomato.")
    default:
        fmt.Printf("I did not have lunch.")
    }
}