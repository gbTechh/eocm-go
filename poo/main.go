package main

import "fmt"


func main()  {
	Go := Course {
		Name: "Go desde cero",
		Price: 13.44,
		IsFree: false,
		UserIDs: []uint{12,23,45},
		Classes: map[uint]string {
			1: "Introduccion",
			2: "Estructuras",
			3: "Maps",
		},
	}

	css := Course{Name: "CSS dede cero", IsFree: true}
	js := Course{}
	js.Name = "Js desde cero"
	js.UserIDs = []uint{32,32}
	fmt.Println(Go.Name)
	fmt.Printf("%+v\n",css)
	fmt.Printf("%+v",js)


	//Go.PrintClasses()
	Go.ChangePrice(67.12)
	fmt.Println(Go.Price)

}