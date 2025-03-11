package main

import (
	"log"
	"net/http"
)

func main2() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	log.Println("Servidor iniciando en: http://127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}