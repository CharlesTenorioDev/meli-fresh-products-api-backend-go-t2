package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func NewSellerDBRepository(db map[int]internal.Seller) *SellerDBRepository {
	defaultDB := make(map[int]internal.Seller)
	if db != nil {
		defaultDB = db
	}

	return &SellerDBRepository{db: defaultDB}
}

type SellerDBRepository struct {
	db map[int]internal.Seller
}

func (r *SellerDBRepository) GetAll() (map[int]internal.Seller, error) {
	sellers := make(map[int]internal.Seller)

	for key, value := range r.db {
		sellers[key] = value
	}

	return sellers, nil
}

func (r *SellerDBRepository) GetByID(id int) (internal.Seller, error) {
	seller, exists := r.db[id]
	if !exists {
		return internal.Seller{}, utils.ErrNotFound
	}

	return seller, nil
}

func (r *SellerDBRepository) GetByCid(cid int) (internal.Seller, error) {
	for _, v := range r.db {
		if v.Cid == cid {
			return r.db[v.ID], nil
		}
	}

	return internal.Seller{}, nil
}

func (r *SellerDBRepository) Create(newSeller internal.SellerRequest) (internal.Seller, error) {
	newSellerID := utils.GetBiggestID(r.db) + 1

	createdSeller := internal.Seller{
		ID:          newSellerID,
		Cid:         newSeller.Cid,
		CompanyName: newSeller.CompanyName,
		Address:     newSeller.Address,
		Telephone:   newSeller.Telephone,
	}

	r.db[newSellerID] = createdSeller

	return createdSeller, nil
}

func (r *SellerDBRepository) Update(seller internal.Seller) (internal.Seller, error) {
	r.db[seller.ID] = seller

	return r.db[seller.ID], nil
}

func (r *SellerDBRepository) Delete(id int) (bool, error) {
	delete(r.db, id)
	return true, nil
}
