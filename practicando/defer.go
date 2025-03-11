package main

import (
	"fmt"
	"os"
)

func mainDefer()  {
	
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)

	a:= 5
	defer fmt.Println("Defer: ",a)

	a=10
	fmt.Println(a)

	//Output:
	//10
	//Defer: 5

	/*
	Casos de uso:
	-limpiar recursos
	-cerrar archivos, conexiones de red, controladores de datos
	*/

	file, err := os.Create("hello.txt")
	if err != nil {
		fmt.Printf("ocurrio un error al crear el archivo %v", err)
		return
	}

	defer file.Close()
	_, err = file.Write([]byte("Hola a todos"))

	if err != nil {
		file.Close()
		fmt.Printf("Ocurrio un error al escribir el archivo: %v", err)
		return
	}

}