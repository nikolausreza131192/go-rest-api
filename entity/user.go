package entity

import "time"

// User ...
type User struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"nama"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	Status    int64     `json:"status" db:"status"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	UpdatedBy string    `json:"updated_by" db:"updated_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
