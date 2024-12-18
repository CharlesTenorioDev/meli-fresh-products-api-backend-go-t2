package repository

import (
	"encoding/json"
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

var buyersFile = "buyers.json"

func (repo *BuyerDb) LoadBuyers() (map[int]pkg.Buyer, error) {
	file, err := os.ReadFile(buyersFile)
	if err != nil {
		return nil, err
	}

	var buyers map[int]pkg.Buyer

	if err := json.Unmarshal(file, &buyers); err != nil {
		return nil, err
	}

	return buyers, nil
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
