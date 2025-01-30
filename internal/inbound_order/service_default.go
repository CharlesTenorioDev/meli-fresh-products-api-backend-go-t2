package inbound_order

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type InboundOrderService struct {
	repo internal.InboundOrderRepository
}

func NewInboundOrderService(repo internal.InboundOrderRepository) *InboundOrderService {
	return &InboundOrderService{repo: repo}
}

func (s *InboundOrderService) CreateInboundOrder(newOrder internal.InboundOrderAttributes) (internal.InboundOrder, error) {
	if newOrder.OrderDate == "" || newOrder.OrderNumber == "" || newOrder.EmployeeID == 0 || newOrder.ProductBatchID == 0 || newOrder.WarehouseID == 0 {
		return internal.InboundOrder{}, utils.ErrInvalidArguments
	}
	_, err := s.repo.FindByID(newOrder.EmployeeID)
	if err != nil && errors.Is(err, utils.ErrNotFound) {
		return internal.InboundOrder{}, utils.ErrNotFound
	}
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return internal.InboundOrder{}, err
	}

	//Check for duplicate
	existingOrder, _ := s.repo.FindByOrderNumber(newOrder.OrderNumber)
	if existingOrder.ID != 0 {
		return internal.InboundOrder{}, utils.ErrConflict
	}

	createdInbound, err := s.repo.CreateInboundOrder(newOrder)

	return createdInbound, err
}

func (s *InboundOrderService) GenerateInboundOrdersReport(ids []int) ([]internal.EmployeeInboundOrdersReport, error) {
	var reports []internal.EmployeeInboundOrdersReport

	if len(ids) == 0 {
		var err error

		reports, err = s.repo.GenerateInboundOrdersReport()
		if err != nil {
			return []internal.EmployeeInboundOrdersReport{}, err
		}
	} else {
		for _, id := range ids {
			_, err := s.repo.FindByID(id)
			if err != nil && errors.Is(err, utils.ErrNotFound) {
				return nil, utils.ErrNotFound
			}
			if !errors.Is(err, utils.ErrNotFound) {
				return nil, err
			}

			report, err := s.repo.GenerateByIDInboundOrdersReport(id)
			if err != nil {
				if err == utils.ErrNotFound {
					return nil, utils.ErrNotFound
				}

				return nil, err
			}

			reports = append(reports, report)
		}
	}

	return reports, nil
}
