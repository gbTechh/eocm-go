package handler

import (
	"encoding/json"
	"net/http"
)


type Person struct {
	Name string
	Age int
}

func Get(w http.ResponseWriter, r *http.Request) {
	p := Person {
		"Jhoana",
		31,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}