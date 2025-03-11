package main

import "fmt"

func mainbucles() {
	nums := []uint8{2,4,6,8}

	for range nums {
		fmt.Println("EdTeam")
	}
	for i,v := range nums { //i indice, v valor
		fmt.Println("Indice:", i, "Valor: ",v)
	}
	for i := range nums { //i indice, v valor
		nums[i] *= 2
	}
	fmt.Println(nums)


	//iteracion sobre mapas

	sports := map[string]string{"basketball": "B", "soccer": "S"}
	for key, v := range sports {
		fmt.Println("Key", key, "Valor", v)
	}

	//iteracion sobre strings

	hello := "hello"

	for _,v := range hello {
		fmt.Println(string(v))
	}
}