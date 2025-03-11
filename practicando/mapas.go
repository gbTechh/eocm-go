package main

import "fmt"

func main3() {
	mapa := map[string]string {
		"colobmia": "bogota",
		"peru": "lima",
		"argentina": "buenos aires",
	}

	fmt.Println("Capitales y paises: ", mapa)

  numeros := []int{1, 2, 3, 4}

	if _, ok := mapa["brasil"]; !ok {
		mapa["brasil"] = "brasilia"
	}
    
	mostrarArgumentos(append([]int{}, numeros...)...)
	fmt.Println("Original después:", numeros) // El original no cambia
	
	// vs usar el slice directamente
	mostrarArgumentos(numeros...)
}

func mostrarArgumentos(nums ...int) {
    fmt.Printf("Recibí %d argumentos: %v\n", len(nums), nums)
    // Modificar el slice local
    nums[0] = 999
}

