package repository

import (
	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func NewSellerDbRepository(db map[int]pkg.Seller) *SellerDbRepository {
	defaultDb := make(map[int]pkg.Seller)
	if db != nil {
		defaultDb = db
	}
	return &SellerDbRepository{db: defaultDb}

}

type SellerDbRepository struct {
	db map[int]pkg.Seller
}

func (r *SellerDbRepository) GetAll() (map[int]pkg.Seller, error) {
	sellers := make(map[int]pkg.Seller)

	for key, value := range r.db {
		sellers[key] = value
	}

	return sellers, nil
}

func (r *SellerDbRepository) GetById(id int) (pkg.Seller, error) {
	seller, exists := r.db[id]
	if !exists {
		return pkg.Seller{}, utils.ErrNotFound
	}
	return seller, nil
}

func (r *SellerDbRepository) GetByCid(cid int) (pkg.Seller, error) {
	for _, v := range r.db {
		if v.Cid == cid {
			return r.db[v.ID], nil
		}
	}
	return pkg.Seller{}, nil
}

func (r *SellerDbRepository) Create(newSeller pkg.SellerRequest) (pkg.Seller, error) {
	newSellerId := utils.GetBiggestId(r.db) + 1

	createdSeller := pkg.Seller{
		ID:          newSellerId,
		Cid:         newSeller.Cid,
		CompanyName: newSeller.CompanyName,
		Address:     newSeller.Address,
		Telephone:   newSeller.Telephone,
	}

	r.db[newSellerId] = createdSeller
	return createdSeller, nil

}

func (r *SellerDbRepository) Update(seller pkg.Seller) (pkg.Seller, error) {
	r.db[seller.ID] = seller

	return r.db[seller.ID], nil
}

func (r *SellerDbRepository) Delete(id int) (bool, error) {
	delete(r.db, id)
	return true, nil
}
