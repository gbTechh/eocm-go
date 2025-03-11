package invoiceitem

import "time"

type Model struct {
	ID uint
	InvoiceHeaderID uint
	ProductID uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage interface{
	Migrate() error
}