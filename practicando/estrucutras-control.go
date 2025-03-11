package main

import "fmt"

func mainswtich() {

	emoji := "gato"
	switch emoji {
	case "gato":
		fmt.Println("Es un gato")
		break;
	case "perro":
		fmt.Println("Es un perro")
		break;
	case "raton":
		fmt.Println("Es un raton")
		break;
	default:
		fmt.Println("No se lo que es")
		break
	}

	emoji2 := "manzana"
	switch emoji2 {
	case "gato", "perro":
		fmt.Println("Es un animal")
		break;
	case "manzana", "platano":
		fmt.Println("Es una fruta")
		break;
	default:
		fmt.Println("No se lo que es")
		break
	}
}