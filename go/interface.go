package main
import "fmt"
type animal interface{
	move() error
	breath() error
}
type readThinker interface{
	read() error
	think() error
}
type people interface{
	animal
	readThinker
}
type student struct{
	name string
}
func (s student) move() error{
	fmt.Println(s.name+" move")
	return nil
}
func (s student) breath() error{
	fmt.Println(s.name+" breath")
	return nil
}
func (s student) read() error{
	fmt.Println(s.name+" read")
	return nil
}
func (s student) think() error{
	fmt.Println(s.name+" think")
	return nil
}
func (s student) study() error{
	fmt.Println(s.name+" study")
	return nil
}
type dog struct{
	name string
}
func (d dog) move() error{
	fmt.Println(d.name+" move")
	return nil
}
func (d dog) breath() error{
	fmt.Println(d.name+" breath")
	return nil
}
func moveAndBreath(a animal) {
	_ = a.move()
	_ = a.breath()
}
func readAndThink(r readThinker){
	_ = r.read()
	_ = r.think()
}
func moveAndThink(p people){
	_ = p.move()
	_ = p.think()
}
func main(){
	ol := student{"Olivia"}
	pug := dog{"Bella"}
	moveAndBreath(ol)
	moveAndBreath(pug)
	readAndThink(ol)
	moveAndThink(ol)
}
