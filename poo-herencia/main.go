package main

import (
	"fmt"
	"poo/customer"
	"poo/invoice"
	"poo/invoiceItem"
)

func main()  {
	i := invoice.New(
		"Colombia", 
		"Popayan",
		customer.New("Alejandro", "direcicon", "98765432"),
		invoiceItem.NewItems(
			invoiceItem.New(1, "Curso de go", 12.34),
			invoiceItem.New(2, "Curso de go", 12.34),
			invoiceItem.New(3, "Curso de go", 12.34),
		),
	)
	i.SetTotal()

	fmt.Println(i)
}