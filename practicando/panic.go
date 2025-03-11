package main

import "fmt"

func mainPanic()  {
	div(10,2)
	div(10,1)
	div(10,0)
}

func div(dividendo, divisor int){
	validarDivisor(divisor)
	fmt.Println(dividendo/divisor)
}

func validarDivisor(divisor int){
	if divisor == 0 {
		panic("Panico!!!!")
	}
}