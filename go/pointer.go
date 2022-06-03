package main
import "fmt"
type student struct {
	name   string
	id string
}
func main() {
	animals := []string{
		"dog",
		"lion",
		"panda",
	}
	adam := student{"Adam", "abcdef"}
	fmt.Println(adam)
	modifyStudentName(&adam, "Levi")
	fmt.Println(adam)
	fmt.Println(animals)
	modifyFirstItem(animals, "cat")
	fmt.Println(animals)
}
func modifyStudentName(pointerToStudent *student, newName string) {
	(*pointerToStudent).name = newName
}
func modifyFirstItem(animals []string, newFirstItem string) {
	animals[0] = newFirstItem
}