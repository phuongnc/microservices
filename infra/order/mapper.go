package order

func MapOrderToModel(order *OrderDto) *OrderModel {
	if order == nil {
		return nil
	}
	return &OrderModel{
		Id:     order.Id,
		UserId: order.UserId,
		Amount: order.Amount,
		Detail: order.Detail,
	}
}

func mapOrderFromModel(order *OrderModel) *OrderDto {
	if order == nil {
		return nil
	}
	return &OrderDto{
		Id:        order.Id,
		UserId:    order.UserId,
		Amount:    order.Amount,
		Detail:    order.Detail,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}

}

func mapOrderToGorm(order *OrderModel) *OrderEntity {
	if order == nil {
		return nil
	}
	return &OrderEntity{
		Id:        order.Id,
		UserId:    order.UserId,
		Amount:    order.Amount,
		Detail:    order.Detail,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func mapOrderFromGorm(order *OrderEntity) *OrderModel {
	if order == nil {
		return nil
	}
	return &OrderModel{
		Id:        order.Id,
		UserId:    order.UserId,
		Amount:    order.Amount,
		Detail:    order.Detail,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
