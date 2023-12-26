package order

import "time"

type (
	UUID = string
)

type OrderModel struct {
	Id        UUID
	UserId    UUID
	Amount    float32
	Detail    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
