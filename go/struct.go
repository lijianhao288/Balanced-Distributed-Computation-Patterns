package main
import "fmt"
type student struct {
	name   string
	id id
	c contact
}
type id string
type contact struct {
	email   string
	mobile string
}
func main() {
	cont:= contact{"adam@gmail.com","0000"}
	adam := student{"Adam",  "abcdef",cont}
	fmt.Println(adam)
	adam.id = "nnnnn"
	fmt.Println(adam)
	fmt.Println(adam.name)
	type teacherId string
	type teacher struct{
		n string
		id teacherId
		c contact
	}
	v:= teacher {"v","a",contact{"v@gmail.com", "1111"}}
	fmt.Println(v)
}

