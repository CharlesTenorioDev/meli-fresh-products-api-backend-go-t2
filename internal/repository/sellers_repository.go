package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func NewSellerDbRepository(db map[int]internal.Seller) *SellerDbRepository {
	defaultDb := make(map[int]internal.Seller)
	if db != nil {
		defaultDb = db
	}
	return &SellerDbRepository{db: defaultDb}

}

type SellerDbRepository struct {
	db map[int]internal.Seller
}

func (r *SellerDbRepository) GetAll() (map[int]internal.Seller, error) {
	sellers := make(map[int]internal.Seller)

	for key, value := range r.db {
		sellers[key] = value
	}

	return sellers, nil
}

func (r *SellerDbRepository) GetById(id int) (internal.Seller, error) {
	seller, exists := r.db[id]
	if !exists {
		return internal.Seller{}, utils.ErrNotFound
	}
	return seller, nil
}

func (r *SellerDbRepository) GetByCid(cid int) (internal.Seller, error) {
	for _, v := range r.db {
		if v.Cid == cid {
			return r.db[v.ID], nil
		}
	}
	return internal.Seller{}, nil
}

func (r *SellerDbRepository) Create(newSeller internal.SellerRequest) (internal.Seller, error) {
	newSellerId := utils.GetBiggestId(r.db) + 1

	createdSeller := internal.Seller{
		ID:          newSellerId,
		Cid:         newSeller.Cid,
		CompanyName: newSeller.CompanyName,
		Address:     newSeller.Address,
		Telephone:   newSeller.Telephone,
	}

	r.db[newSellerId] = createdSeller
	return createdSeller, nil

}

func (r *SellerDbRepository) Update(seller internal.Seller) (internal.Seller, error) {
	r.db[seller.ID] = seller

	return r.db[seller.ID], nil
}

func (r *SellerDbRepository) Delete(id int) (bool, error) {
	delete(r.db, id)
	return true, nil
}
