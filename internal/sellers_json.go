package internal

import (
	"encoding/json"
	"os"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func NewSellerJSONFile(path string) *SellerJSONFile {
	return &SellerJSONFile{
		path: path,
	}
}

type SellerJSONFile struct {
	path string
}

type SellerJSON struct {
	ID            int     `json:"id"`
	Cid           int     `json:"cid"`
	CompanyName   string  `json:"company_name"`
	Address       string  `json:"address"`
	Telephone     string  `json:"telephone"`
}

// Load is a method that loads the vehicles
func (l *SellerJSONFile) Load() (v map[int]pkg.Seller, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var sellersJSON []SellerJSON
	err = json.NewDecoder(file).Decode(&sellersJSON)
	if err != nil {
		return
	}

	// serialize vehicles
	v = make(map[int]pkg.Seller)
	for _, vh := range sellersJSON {
		v[vh.ID] = pkg.Seller{
			ID: 		  vh.ID,
			Cid:          vh.Cid,
			CompanyName:  vh.CompanyName,
			Address:      vh.Address,
			Telephone:    vh.Telephone,
		}
	}

	return
}
