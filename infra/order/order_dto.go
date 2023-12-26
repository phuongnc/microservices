package order

import "time"

type OrderDto struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	Amount    float32   `json:"amount"`
	Detail    string    `json:"detail"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
