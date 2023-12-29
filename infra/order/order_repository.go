package order

import (
	"context"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *OrderModel) error
	UpdateOrder(ctx context.Context, order *OrderModel) error
	Query(ctx context.Context) OrderQuery
}

type OrderQuery interface {
	ById(id UUID) OrderQuery
	Result() (*OrderModel, error)
}

func NewOrderRepo() OrderRepository {
	return &orderRepo{}
}

type orderRepo struct{}

func (repo *orderRepo) CreateOrder(ctx context.Context, o *OrderModel) error {
	db := ctx.Value("db").(*gorm.DB)
	orderGM := mapOrderToGorm(o)
	if err := db.Create(orderGM).Error; err != nil {
		return err
	}
	return nil
}

func (repo *orderRepo) UpdateOrder(ctx context.Context, o *OrderModel) error {
	db := ctx.Value("db").(*gorm.DB)
	orderGM := mapOrderToGorm(o)
	return db.Model(orderGM).Updates(orderGM).Error
	return nil
}

func (repo *orderRepo) Query(ctx context.Context) OrderQuery {
	return &orderQuery{db: ctx.Value("db").(*gorm.DB).Model(&OrderEntity{})}
}

type orderQuery struct {
	db *gorm.DB
}

func (query orderQuery) ById(id UUID) OrderQuery {
	return &orderQuery{db: query.db.Where("id = ?", id)}
}

func (query orderQuery) Result() (*OrderModel, error) {
	result := OrderEntity{}

	err := query.db.First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapOrderFromGorm(&result), nil
}
