package service

import (
	"log"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

type BuyerService struct {
	repo pkg.BuyerRepository
}

func NewBuyer(repo pkg.BuyerRepository) *BuyerService {
	return &BuyerService{repo: repo}
}

func (service *BuyerService) GetAll() (buyer []pkg.Buyer, err error) {
	buyer, err = service.repo.GetAll()
	return buyer, err
}

func (service *BuyerService) GetOne(id int) (*pkg.Buyer, error) {
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

func (service *BuyerService) CreateBuyer(buyer pkg.BuyerAttributes) (*pkg.Buyer, error) {
	buyers, err := service.GetAll()
	if err != nil {
		log.Println("Error to load - ", err)
		return nil, err
	}

	id := getNextID(buyers)

	newBuyer := pkg.Buyer{
		ID: int64(id),
		BuyerAttributes: pkg.BuyerAttributes{
			CardNumberID: buyer.CardNumberID,
			FirstName:    buyer.FirstName,
			LastName:     buyer.LastName,
		},
	}

	err = service.validation(newBuyer)
	if err != nil {
		return nil, err
	}

	service.repo.CreateBuyer(newBuyer)

	return &newBuyer, nil
}

func (service *BuyerService) UpdateBuyer(updatedBuyer *pkg.Buyer) (*pkg.Buyer, error) {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error internal - ", err)
	}

	var buyerFound pkg.Buyer
	for _, buyer := range buyers {
		if buyer.ID == updatedBuyer.ID {
			buyerFound = buyer
		}
	}

	if buyerFound == (pkg.Buyer{}) {
		return nil, utils.ErrNotFound
	}

	for _, buyer := range buyers {
		if buyer.CardNumberID == updatedBuyer.CardNumberID {
			log.Println("Error Card number already exist in our system")
			return nil, utils.ErrConflict
		}
	}

	service.repo.UpdateBuyer(updatedBuyer)

	return updatedBuyer, nil
}

func (service *BuyerService) DeleteBuyer(id int) error {
	buyers, err := service.GetAll()

	if err != nil {
		log.Println("Error in GetAll - ", err)
		return err
	}

	for _, buyer := range buyers {
		if int(buyer.ID) == id {
			service.repo.DeleteBuyer(id)
			return nil
		}
	}

	return err
}

func (service *BuyerService) validation(newBuyer pkg.Buyer) error {
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

func getNextID(buyers []pkg.Buyer) int {
	maxID := 0
	for _, buyer := range buyers {
		if int(buyer.ID) > maxID {
			maxID = int(buyer.ID)
		}
	}
	return maxID + 1
}
