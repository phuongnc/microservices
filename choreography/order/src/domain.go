package src

import (
	"infra/log"

	"github.com/labstack/echo/v4"
)

type OrderDomain interface {
	CreateOrder(ctx echo.Context) error
}

type orderDomain struct {
	logger *log.Logger
	//orderRepo repo.SaleSummaryRepo
}

func NewOrderDomain(
	logger *log.Logger,
	//orderRepo repo.SaleSummaryRepo,
) OrderDomain {
	return &orderDomain{
		logger,
		//orderRepo,
	}
}

func (s *orderDomain) CreateOrder(ctx echo.Context) error {
	s.logger.Info("Getting sale summary data")
	// result, err := s.saleSummaryRepo.Query(ctx).ByProcessingDateRange(fromDate, toDate).ResultList()
	// if err != nil {
	// 	s.logger.Error("Can not get sale summary data", err)
	// 	return nil, err
	// }
	return nil
}
