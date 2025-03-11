package tag

import "time"

type Tag struct {
	ID          int64     `json:"id"`
	Name    		string    `json:"name"`
	Code				string    `json:"code"`
	Active    	bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}