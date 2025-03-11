package main

import "fmt"

func mainstrutc() {
	type course struct {
		Name string
		Profesor string
		Country string
	}

	db := course {
		Name: "Base de datos",
		Profesor: "Alexxys",
		Country: "Colombia",
	}

	git := course{"Git", "Beto", "Bolivia"}

	css := course{Profesor: "Alvaro"}

	fmt.Printf("%+v\n", db)
	fmt.Printf("%+v\n", git)
	fmt.Printf("%+v\n", css)
	fmt.Printf("%+v\n", db.Name)
	fmt.Printf("%+v\n", git.Country)
	fmt.Printf("%+v\n", css.Profesor)

	p := &db
	(*p).Country = "Peru"
	p.Country = "Peru" //esto es igual a lo de arriba pero sin la necesidad de desferencia (*p)
	fmt.Printf("%+v\n", db)
	fmt.Printf("%+v\n", p)
}