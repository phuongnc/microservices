package order

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type OrderRepository interface {
	CreateOrder(ctx echo.Context, order *OrderModel) error
	UpdateOrder(ctx echo.Context, order *OrderModel) error
	Query(ctx echo.Context) OrderQuery
}

type OrderQuery interface {
	ById(id UUID) OrderQuery
	Result() (*OrderModel, error)
}

func NewOrderRepo() OrderRepository {
	return &orderRepo{}
}

type orderRepo struct{}

func (repo *orderRepo) CreateOrder(ctx echo.Context, o *OrderModel) error {
	db := ctx.Get("db").(*gorm.DB)
	orderGM := mapOrderToGorm(o)
	if err := db.Create(orderGM).Error; err != nil {
		return err
	}
	return nil
}

func (repo *orderRepo) UpdateOrder(ctx echo.Context, o *OrderModel) error {
	db := ctx.Get("db").(*gorm.DB)
	orderGM := mapOrderToGorm(o)
	return db.Model(orderGM).Update(orderGM).Error
}

func (repo *orderRepo) Query(ctx echo.Context) OrderQuery {
	return &orderQuery{db: ctx.Get("db").(*gorm.DB).Model(&OrderEntity{})}
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
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return mapOrderFromGorm(&result), nil
}
