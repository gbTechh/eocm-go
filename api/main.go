package main

import (
	"fmt"
	"net/http"
)

func saludar(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hola mundo")
}

func main()  {
	//registrando la ruta y el handler
	http.HandleFunc("/saludar", saludar)
	//sube usn ervidor ne el puerto 8080
	http.ListenAndServe(":8000", nil)
}