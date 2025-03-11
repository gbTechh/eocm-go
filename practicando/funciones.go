package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func mainFn() {
	hello()
	naming("Rodrigo", "silva")
	comida := "papa"
	change(comida)
	fmt.Println("\nComida",comida)

	changeRef(&comida)
	fmt.Println("\nComida", comida)

	text := "Hola mundo COMO estas"
	s1, s2 :=convert(text)
	fmt.Println(s1,s2)

	//lee archivo
	content, err := os.ReadFile("./things.txt")
	if err != nil {
		fmt.Printf("Ocurrio un error; %v", err)
		return;
	} 
	fmt.Println(string(content))
	

	//funcion con errores
	res, err := division(30,2)
	if err != nil {
		fmt.Printf("Ocurrio un error; %v", err)
		return;
	} 
	fmt.Println(res)

	//funcion que retorna funcion 
	funcRetfunc()

	//funcion variatica
	fmt.Println(sum3(1,2,6,3,4))
}

func hello() {
	fmt.Println("Hola")
}

func naming(firstName string, lastName string) {
	fmt.Printf("Hello %s %s", firstName, lastName)
}

//se pasa como copia y no cambia le valor
func change(value string)  {
	value = "arroz"
}
//se pasa como referencia y cambia el valor (se tiene que pasar la direccion)
func changeRef(value *string){
	*value = "arroz"
}

//funcion con retorno de valor
func sum(n1 int, n2 int) int  {
	return n1 + n2
}

func sum2(n1, n2 int) int  { // si dos parametros tienen el mis mo tipo de dato puedo especifciar el tiopo de dato en el ultimo paarametro
	return n1 + n2
}

//funcio para retornar multipels valores
func convert(text string)(string, string)  {
	min := strings.ToLower(text)
	may := strings.ToUpper(text)

	return min, may
}

//funciones con errores
func division(dividendo, divisor int) (int, error) {
	if divisor == 0 {
		return 0, errors.New("No peudes dividir por cero")
	}
	return (dividendo / divisor), nil
}
func division2(dividendo, divisor int) (result int, err error) { // return con parametros nombrados (no es muy usada ni reocmendada)
	if divisor == 0 {
		err = errors.New("No peudes dividir por cero")
		return result, err
	}
	result = dividendo / divisor
	return result, err
}

//funciones que reciben y retornan funciones
func funcRetfunc()  {
	nums := []int{1,4,6,78,56,2}
	result := filter(nums, func(num int) bool {
		return num <= 10
	})

	fmt.Println(result)
}

func filter(nums []int, callback func(int) bool) []int {
	result := []int{}
	for _,v := range nums {
		if callback(v) {
			result = append(result,v)
		}
	}
	return result
}

//funcion variatica
func sum3(nums ...int) int {
	total := 0
	for _,v := range nums {
		total += v
	}
	return total
	
}