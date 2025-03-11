package invoiceheader

import "time"

type Model struct {
	ID uint
	Client string
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Storage interface {
	Migrate() error
}