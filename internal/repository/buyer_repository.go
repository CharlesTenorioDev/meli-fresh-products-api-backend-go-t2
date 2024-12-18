package repository

import (
	"encoding/json"
	"os"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

var buyersFile = "buyers.json"

func LoadBuyers() (map[int]pkg.Buyer, error) {
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

func GetAllBuyers() ([]pkg.Buyer, error) {
	buyersMap, err := LoadBuyers()
	if err != nil {
		return nil, err
	}

	var buyers []pkg.Buyer
	for _, buyer := range buyersMap {
		buyers = append(buyers, buyer)
	}

	return buyers, nil
}
