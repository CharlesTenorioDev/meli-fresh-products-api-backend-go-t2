package service

import (
	"log"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BuyerService struct {
	repo internal.BuyerRepository
}

func NewBuyer(repo internal.BuyerRepository) *BuyerService {
	return &BuyerService{repo: repo}
}

func (service *BuyerService) GetAll() (buyer []internal.Buyer, err error) {
	buyer, err = service.repo.GetAll()
	return buyer, err
}

func (service *BuyerService) GetOne(id int) (*internal.Buyer, error) {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error in GetAll - ", err)
		return nil, err
	}

	for _, buyer := range buyers {
		if int(buyer.ID) == id {
			return &buyer, err
		}
	}

	return nil, err
}

func (service *BuyerService) CreateBuyer(buyer internal.BuyerAttributes) (*internal.Buyer, error) {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error to load - ", err)
		return nil, err
	}

	id := getNextID(buyers)

	newBuyer := internal.Buyer{
		ID: int64(id),
		BuyerAttributes: internal.BuyerAttributes{
			CardNumberID: buyer.CardNumberID,
			FirstName:    buyer.FirstName,
			LastName:     buyer.LastName,
		},
	}

	err = service.validation(newBuyer)
	if err != nil {
		return nil, err
	}

	return service.repo.CreateBuyer(newBuyer)
}

func (service *BuyerService) UpdateBuyer(updatedBuyer *internal.Buyer) (*internal.Buyer, error) {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error internal - ", err)
	}

	var buyerFound internal.Buyer

	for _, buyer := range buyers {
		if buyer.ID == updatedBuyer.ID {
			buyerFound = buyer
		}
	}

	if buyerFound == (internal.Buyer{}) {
		return nil, utils.ErrNotFound
	}

	for _, buyer := range buyers {
		if buyer.CardNumberID == updatedBuyer.CardNumberID {
			log.Println("Error Card number already exist in our system")
			return nil, utils.ErrConflict
		}
	}

	return service.repo.UpdateBuyer(updatedBuyer)
}

func (service *BuyerService) DeleteBuyer(id int) error {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error in GetAll - ", err)
		return err
	}

	for _, buyer := range buyers {
		if int(buyer.ID) == id {
			if err := service.repo.DeleteBuyer(id); err != nil {
				return err
			}

			return nil
		}
	}

	return err
}

func (service *BuyerService) validation(newBuyer internal.Buyer) error {
	buyers, err := service.repo.GetAll()

	if err != nil {
		log.Println("Error in load -", err)
		return err
	}

	for _, buyer := range buyers {
		if buyer.ID == newBuyer.ID {
			log.Println("There is a user with this ID - ", err)
			return utils.ErrConflict
		}

		if buyer.CardNumberID == newBuyer.CardNumberID {
			log.Println("There is a user with this card number - ", err)
			return utils.ErrConflict
		}
	}

	return nil
}

func getNextID(buyers []internal.Buyer) int {
	maxID := 0
	for _, buyer := range buyers {
		if int(buyer.ID) > maxID {
			maxID = int(buyer.ID)
		}
	}

	return maxID + 1
}
