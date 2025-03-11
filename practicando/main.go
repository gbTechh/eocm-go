package main

import "fmt"

func main2() {
	var nombre string = "Leonardo"
	var apellido = "Caceres"
	edad := 32
	fmt.Println("Hola mundo con go " + nombre + " " + apellido + " tengo ", edad)	
	//arreglos
	var listaFrutas = [4]string{"Pera", "manzana"}
	lent := len(listaFrutas)
	fmt.Println(listaFrutas, lent)

	listaPaises := []string{"Peru", "Venezuela", "Chile"}
	listaPaises = append(listaPaises, "Brazil")
	listaPaises2 := listaPaises[1:5]
	fmt.Println(listaPaises)
	fmt.Println(listaPaises2)

	original := []int{1, 2, 3, 4, 5}
	length := len(original)
	fmt.Printf("Original: len=%d cap=%d\n", length, cap(original))
	// Slice desde el índice 2
	slice1 := original[2:]
	fmt.Printf("Slice1: len=%d cap=%d\n", len(slice1), cap(slice1))
	// La capacidad será 3 porque cuenta desde el índice 2 hasta el final

	// Slice con límite de capacidad
	slice2 := original[2:4:4]  // [bajo:alto:máximo]
	fmt.Printf("Slice2: len=%d cap=%d\n", len(slice2), cap(slice2))
	// La capacidad será 2 porque la limitamos explícitamente
}

