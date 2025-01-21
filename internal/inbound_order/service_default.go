package inbound_order

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type InboundOrderService struct {
	repo internal.InboundOrderRepository
}

func NewInboundOrderService(repo internal.InboundOrderRepository) *InboundOrderService {
	return &InboundOrderService{repo: repo}
}

func (s *InboundOrderService) CreateOrder(newOrder internal.InboundOrderAttributes) (internal.InboundOrder, error) {
	if newOrder.OrderDate == "" || newOrder.OrderNumber == "" || newOrder.EmployeeID == 0 || newOrder.ProductBatchID == 0 || newOrder.WarehouseID == 0 {
		return internal.InboundOrder{}, utils.ErrInvalidArguments
	}

	//Check for duplicate order number
	existingOrder, _ := s.repo.FindByOrderNumber(newOrder.OrderNumber)
	if existingOrder.ID != 0 {
		return internal.InboundOrder{}, utils.ErrConflict
	}

	return s.repo.Create(newOrder)
}

func (s *InboundOrderService) GenerateInboundOrdersReport(ids []int) ([]internal.EmployeeInboundOrdersReport, error) {
	var reports []internal.EmployeeInboundOrdersReport

	if len(ids) == 0 {
		var err error

		reports, err = s.repo.GenerateReport()
		if err != nil {
			return []internal.EmployeeInboundOrdersReport{}, err
		}
	} else {
		for _, id := range ids {
			report, err := s.repo.GenerateReportForEmployee(id)
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
