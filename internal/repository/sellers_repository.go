package repository

import "github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"


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
	return r.db[id], nil

}

func (r *SellerDbRepository) GetByCid(cid int) (pkg.Seller, error) {
	return r.db[cid], nil
}

func (r *SellerDbRepository) Create(seller pkg.Seller) (pkg.Seller, error) {
	

}