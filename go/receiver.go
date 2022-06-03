package main
import "fmt"
type student struct {
	name   string
	id string
}
func main() {
	adam := student{"Adam", "abcdef"}
	fmt.Println(adam.getNameR())
	fmt.Println(getNameA(adam))
}
func (s student) getNameR() string {
	return s.name
}
func getNameA(s student) string {
	return s.name
}
