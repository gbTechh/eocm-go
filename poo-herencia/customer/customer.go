package customer


type Customer struct {
	name string
	address string
	phone string
}

func New(name, address, phone string) Customer {
	return Customer{name, address, phone}
}