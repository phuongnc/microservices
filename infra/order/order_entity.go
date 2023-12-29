package order

import "time"

type OrderEntity struct {
	Id            string `gorm:"primaryKey;type:VARCHAR(255)"`
	UserId        string `gorm:"type:VARCHAR(255)"`
	Amount        float32
	Detail        string `gorm:"type:VARCHAR(1000)"`
	Status        string `gorm:"type:VARCHAR(255)"`
	SubStatus     string `gorm:"type:VARCHAR(255)"`
	FailureReason string `gorm:"type:VARCHAR(1000)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (o *OrderEntity) TableName() string {
	return "order"
}
