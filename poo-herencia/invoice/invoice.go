package invoice

import (
	"poo/customer"
	invoiceitem "poo/invoiceItem"
)

type Invoice struct {
	country string
	city string
	total float64
	cliente customer.Customer
	items invoiceitem.Items
}

func New(country, city string, cliente customer.Customer, items invoiceitem.Items) Invoice {
	return Invoice {
		country: country,
		city: city,
		cliente: cliente,
		items: items,
	}
}

func (i *Invoice) SetTotal()  {
	i.total = i.items.Total()
}