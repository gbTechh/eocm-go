package taxzone

import "time"

type TaxZone struct {
	ID          int64     `json:"id"`
	IDZone      int64     `json:"id_zone"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      bool      `json:"status"`
	Zone        *Zone     `json:"zone,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Zone struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
