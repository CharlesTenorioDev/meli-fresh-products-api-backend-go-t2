package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

type BuyerDb struct {
	buyerTable map[int]pkg.Buyer
}

func NewBuyerDb(buyerTab map[int]pkg.Buyer) *BuyerDb {

	buyerDb := make(map[int]pkg.Buyer)
	if buyerTab != nil {
		buyerDb = buyerTab
	}
	return &BuyerDb{buyerTable: buyerDb}
}

var buyersFile = "/Users/peddpereira/Desktop/Project/meli-fresh-products-api-backend-go-t2/buyers.json"

func (repo *BuyerDb) LoadBuyers() (map[int]pkg.Buyer, error) {
	file, err := os.ReadFile(buyersFile)
	if err != nil {
		log.Println("Error to read file", err)
		return nil, err
	}

	var buyers []pkg.Buyer

	if err := json.Unmarshal(file, &buyers); err != nil {
		log.Println("Error to unmarshal - ")
		return nil, err
	}

	repo.buyerTable = make(map[int]pkg.Buyer)
	for _, buyer := range buyers {
		repo.buyerTable[int(buyer.ID)] = buyer
	}

	return repo.buyerTable, nil
}

func (repo *BuyerDb) GetAll() ([]pkg.Buyer, error) {
	buyersMap, err := repo.LoadBuyers()
	if err != nil {
		return nil, err
	}

	var buyers []pkg.Buyer
	for _, buyer := range buyersMap {
		buyers = append(buyers, buyer)
	}

	return buyers, nil
}
