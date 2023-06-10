package command_model

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Age       int       `json:"age"  db:"age"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Person) GetFullName() string {
	return p.FirstName + " " + p.LastName
}
