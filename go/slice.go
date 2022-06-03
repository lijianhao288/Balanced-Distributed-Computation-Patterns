package main
import "fmt"
func main() {
	animals := []string{
		"dog",
		"cat",
		"bird",
		"lion",
	}
	animals = append(animals, "panda")
	animals = append(animals, "tiger", "wolf")
	for index, animal := range animals {
		fmt.Println(index, animal)
	}
	for _, animal := range animals {
		fmt.Println(animal)
	}
	fmt.Println(animals[3])
	fmt.Println(animals[:3])
	fmt.Println(animals[3:])
	fmt.Println(animals[2:4])
	animals[3] = "SSS"
	fmt.Println(animals)
}
