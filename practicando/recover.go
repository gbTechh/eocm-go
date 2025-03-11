package main

import "fmt"

func main()  {
	div2(10,2)
	div2(10,1)
	div2(10,0)
}

func div2(dividendo, divisor int){
	defer func() {
		if r:= recover(); r != nil {
			fmt.Println("recuperandime del panic", r)
		}
	}()
	validarDivisor2(divisor)
	fmt.Println(dividendo/divisor)
}

func validarDivisor2(divisor int){
	if divisor == 0 {
		panic("Panico!!!!")
	}
}