package order

import "time"

type OrderDto struct {
	Id            string    `json:"id"`
	UserId        string    `json:"userId"`
	Amount        float32   `json:"amount"`
	Detail        string    `json:"detail"`
	Status        string    `json:"status"`
	SubStatus     string    `json:"subStatus"`
	FailureReason string    `json:"failureReason"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
