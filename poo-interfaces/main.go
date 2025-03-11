package main

import "fmt"

type Greeter interface {
	Greet()
}

type Byer interface {
	Bye()
}

//embeber interfaces
type GreeterByer interface {
	Greeter
	Byer
}

type Person struct {
	Name string
}

func (p Person) Greet()  {
	fmt.Printf("Hola soy %s\n", p.Name)
}
func (p Person) Bye()  {
	fmt.Printf("Adios soy %s\n", p.Name)
}

type Text string

func (t Text) Greet()  {
	fmt.Printf("Hola soy %s\n", t)
}
func (t Text) Bye()  {
	fmt.Printf("Adios soy %s\n", t)
}

func All(gs ...GreeterByer) {
	for _, gb := range gs {
		gb.Greet()
		gb.Bye()
	}
}



func main()  {
	//var g Greeter = Person{Name: "alejandro"}
	//var g2 Greeter = Text("Daisy")
	//g.Greet()
	//g2.Greet()

	p := Person{Name: "Francisco"}
	var t Text = "Poncio"
	fmt.Println()
	All(p, t)
}