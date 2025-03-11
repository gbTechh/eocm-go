package main

import (
	"fmt"
	"strings"
)

type PayMethod interface {
	Pay()
}

type Paypal struct {

}

type Cash struct {
	
}

type CreditCart struct {

}

func (p Paypal) Pay()  {
	fmt.Println("pagado por paypal")
}	
	
func (p Cash) Pay()  {
	fmt.Println("pagado por efectivo")
}	
	

func (c CreditCart) Pay()  {
	fmt.Println("Pagado con tarjeta de credito")
}

func Factory(method uint) PayMethod {
	switch method {
	case 1:
		return Paypal{}
	case 2:
		return Cash{}
	
	case 3:
		return CreditCart{}
	default:
		return nil
	}
}

//inerfaces vacias
func wrapper(i interface{})  {
	fmt.Printf("valor: %v, Tipo: %T\n", i, i)

	//saber el tipo de la interface (type assertions)
	v,ok := i.(string)
	if ok {
		fmt.Println(strings.ToUpper(v),"Es de tipo string")
	}

	//swtich type para saber el tipo de interface
	switch v := i.(type) {
	case string:
		fmt.Println(strings.ToUpper(v),"ES DE TIPO string")
	case int:
		fmt.Println(v * 2)
	default:
		fmt.Printf("valor: %v, Tipo %T\n",v, v)
	}
}

func main()  {
	payMethod := Factory(2)
	payMethod.Pay()

	wrapper(12)
	wrapper(12.32)
	wrapper(false)
	wrapper("Alejandro")
}