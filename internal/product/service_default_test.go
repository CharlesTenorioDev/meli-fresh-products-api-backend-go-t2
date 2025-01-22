package product

import (
	"reflect"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockProductRepository struct {
	mock.Mock
}

type mockProductTypeValidation struct {
	mock.Mock
}

type mockSellerValidation struct {
	mock.Mock
}

func (m *mockProductRepository) GetAll() (listProducts []internal.Product, err error) {
	args := m.Called()
	return args.Get(0).([]internal.Product), args.Error(1)
}

func (m *mockProductRepository) GetByID(id int) (product internal.Product, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *mockProductRepository) Create(newproduct internal.ProductAttributes) (product internal.Product, err error) {
	args := m.Called(newproduct)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *mockProductRepository) Update(inputProduct internal.Product) (product internal.Product, err error) {
	args := m.Called(inputProduct)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *mockProductRepository) Delete(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockProductTypeValidation) GetProductTypeByID(id int) (productType internal.ProductType, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductType), args.Error(1)
}

func (m *mockSellerValidation) GetByID(id int) (internal.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func TestUnitProduct_GetProducts(t *testing.T) {
	type fields struct {
		repo                  internal.ProductRepository
		validationProductType internal.ProductTypeValidation
		validationSeller      internal.SellerValidation
	}
	tests := []struct {
		name             string
		fields           fields
		wantListProducts []internal.Product
		wantErr          bool
	}{
		{
			name: "GetProducts OK",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
			},
			wantListProducts: []internal.Product{
				{
					ID: 1,
					ProductAttributes: internal.ProductAttributes{
						ProductCode:                    "123",
						Description:                    "Product 1",
						Width:                          1.0,
						Height:                         1.0,
						Length:                         1.0,
						NetWeight:                      1.0,
						ExpirationRate:                 1.0,
						RecommendedFreezingTemperature: 1.0,
						FreezingRate:                   1.0,
						ProductType:                    1,
						SellerID:                       1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "GetProducts Error",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
			},
			wantListProducts: []internal.Product{},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicProductService{
				repo:                  tt.fields.repo,
				validationProductType: tt.fields.validationProductType,
				validationSeller:      tt.fields.validationSeller,
			}
			// Mock the GetAll method
			if tt.wantErr {
				(tt.fields.repo.(*mockProductRepository)).On("GetAll").Return(tt.wantListProducts, utils.ErrNotFound)
			} else {
				(tt.fields.repo.(*mockProductRepository)).On("GetAll").Return(tt.wantListProducts, nil)
			}

			gotListProducts, err := s.GetProducts()
			if (err != nil) != tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.wantListProducts, gotListProducts)
		})
	}
}

func TestUnitProduct_GetProductByID(t *testing.T) {
	type fields struct {
		repo                  internal.ProductRepository
		validationProductType internal.ProductTypeValidation
		validationSeller      internal.SellerValidation
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    internal.Product
		wantErr bool
	}{
		{
			name: "GetProductByID OK",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
			},
			args: args{
				id: 1,
			},
			want: internal.Product{
				ID: 1,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "Product 1",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    1,
					SellerID:                       1,
				},
			},
			wantErr: false,
		},
		{
			name: "GetProductByID Error",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
			},
			args: args{
				id: 1,
			},
			want:    internal.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicProductService{
				repo:                  tt.fields.repo,
				validationProductType: tt.fields.validationProductType,
				validationSeller:      tt.fields.validationSeller,
			}
			// Mock the GetByID method
			if tt.wantErr {
				(tt.fields.repo.(*mockProductRepository)).On("GetByID", tt.args.id).Return(tt.want, utils.ErrNotFound)
			} else {
				(tt.fields.repo.(*mockProductRepository)).On("GetByID", tt.args.id).Return(tt.want, nil)
			}

			gotProduct, err := s.GetProductByID(tt.args.id)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, gotProduct)

		})
	}
}

func TestUnitProduct_CreateProduct(t *testing.T) {
	type fields struct {
		repo                  internal.ProductRepository
		validationProductType internal.ProductTypeValidation
		validationSeller      internal.SellerValidation
	}
	type args struct {
		newProduct internal.ProductAttributes
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        internal.Product
		wantErr     bool
		expectedErr error
	}{
		{
			name: "CreateProduct OK",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "Product 2",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    1,
					SellerID:                       1,
				},
			},
			want: internal.Product{
				ID: 2,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "Product 2",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    1,
					SellerID:                       1,
				},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "CreateProduct Error Conflict",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "Product 1",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    1,
					SellerID:                       1,
				},
			},
			want:        internal.Product{},
			wantErr:     true,
			expectedErr: utils.ErrConflict,
		},
		{
			name: "CreateProduct Error - Empty Field : ProductCode",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "",
					Description:                    "",
					Width:                          0,
					Height:                         0,
					Length:                         0,
					NetWeight:                      0,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : Description",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "",
					Width:                          0,
					Height:                         0,
					Length:                         0,
					NetWeight:                      0,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : Width",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "description",
					Width:                          0,
					Height:                         0,
					Length:                         0,
					NetWeight:                      0,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : height",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "description",
					Width:                          1,
					Height:                         0,
					Length:                         0,
					NetWeight:                      0,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : Length",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "description",
					Width:                          10,
					Height:                         20,
					Length:                         0,
					NetWeight:                      0,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : NetWeight",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "description",
					Width:                          10,
					Height:                         20,
					Length:                         30,
					NetWeight:                      0,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : ExopirationRate",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "description",
					Width:                          10,
					Height:                         20,
					Length:                         30,
					NetWeight:                      40,
					ExpirationRate:                 0,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : RecommendedFreezingTemperature",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "description",
					Width:                          10,
					Height:                         20,
					Length:                         30,
					NetWeight:                      40,
					ExpirationRate:                 50,
					RecommendedFreezingTemperature: 0,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
		{
			name: "CreateProduct Error - Empty Field : FreezingRate",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "1234",
					Description:                    "description",
					Width:                          10,
					Height:                         20,
					Length:                         30,
					NetWeight:                      40,
					ExpirationRate:                 50,
					RecommendedFreezingTemperature: 60,
					FreezingRate:                   0,
					ProductType:                    0,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},

		{
			name: "CreateProduct Error - Invalid Product Type",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "Product 1",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    0,
					SellerID:                       1,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},

		{
			name: "CreateProduct Error - Invalid Seller",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				newProduct: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "Product 1",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    1,
					SellerID:                       0,
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicProductService{
				repo:                  tt.fields.repo,
				validationProductType: tt.fields.validationProductType,
				validationSeller:      tt.fields.validationSeller,
			}
			// Mock the Create method
			if tt.wantErr {
				(tt.fields.repo.(*mockProductRepository)).On("Create", tt.args.newProduct).Return(tt.want, utils.ErrConflict)
			} else {
				(tt.fields.repo.(*mockProductRepository)).On("Create", tt.args.newProduct).Return(tt.want, nil)
			}
			// Mock the GetProductTypeByID method
			(tt.fields.validationProductType.(*mockProductTypeValidation)).On("GetProductTypeByID", tt.args.newProduct.ProductType).Return(internal.ProductType{}, nil)
			// Mock the GetByID method
			(tt.fields.validationSeller.(*mockSellerValidation)).On("GetByID", tt.args.newProduct.SellerID).Return(internal.Seller{}, nil)
			// Mock the GetALL method
			(tt.fields.repo.(*mockProductRepository)).On("GetAll").Return([]internal.Product{
				{
					ID: 1,
					ProductAttributes: internal.ProductAttributes{
						ProductCode:                    "123",
						Description:                    "Product 1",
						Width:                          1.0,
						Height:                         1.0,
						Length:                         1.0,
						NetWeight:                      1.0,
						ExpirationRate:                 1.0,
						RecommendedFreezingTemperature: 1.0,
						FreezingRate:                   1.0,
						ProductType:                    1,
						SellerID:                       1,
					},
				},
			}, nil)

			got, err := s.CreateProduct(tt.args.newProduct)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicProductService.CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicProductService.CreateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitProduct_UpdateProduct(t *testing.T) {
	type fields struct {
		repo                  internal.ProductRepository
		validationProductType internal.ProductTypeValidation
		validationSeller      internal.SellerValidation
	}
	type args struct {
		inputProduct internal.Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    internal.Product
		wantErr bool
	}{
		{
			name: "UpdateProduct OK",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				inputProduct: internal.Product{
					ID: 1,
					ProductAttributes: internal.ProductAttributes{
						ProductCode:                    "1234",
						Description:                    "Product 2",
						Width:                          1.0,
						Height:                         1.0,
						Length:                         1.0,
						NetWeight:                      1.0,
						ExpirationRate:                 1.0,
						RecommendedFreezingTemperature: 1.0,
						FreezingRate:                   1.0,
						ProductType:                    1,
						SellerID:                       1,
					},
				},
			},
			want: internal.Product{
				ID: 1,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "package product",
					Description:                    "package product",
					Width:                          1.0,
					Height:                         1.0,
					Length:                         1.0,
					NetWeight:                      1.0,
					ExpirationRate:                 1.0,
					RecommendedFreezingTemperature: 1.0,
					FreezingRate:                   1.0,
					ProductType:                    1,
					SellerID:                       1,
				},
			},
			wantErr: false,
		},
		{
			name: "UpdateProduct Error - Product Not Found",
			fields: fields{
				repo: &mockProductRepository{
					mock.Mock{},
				},
				validationProductType: &mockProductTypeValidation{
					mock.Mock{},
				},
				validationSeller: &mockSellerValidation{
					mock.Mock{},
				},
			},
			args: args{
				inputProduct: internal.Product{
					ID: 99,
					ProductAttributes: internal.ProductAttributes{
						ProductCode: "1234",
						Description: "Product 2",
						Width:       1.0,
					},
				},
			},
			want:    internal.Product{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicProductService{
				repo:                  tt.fields.repo,
				validationProductType: tt.fields.validationProductType,
				validationSeller:      tt.fields.validationSeller,
			}
			// Mock the GetByID method
			if tt.wantErr {
				(tt.fields.repo.(*mockProductRepository)).On("GetByID", tt.args.inputProduct.ID).Return(internal.Product{}, utils.ErrNotFound)
			} else {
				(tt.fields.repo.(*mockProductRepository)).On("GetByID", tt.args.inputProduct.ID).Return(tt.args.inputProduct, nil)
			}
			// Mock the Update method
			if tt.wantErr {
				(tt.fields.repo.(*mockProductRepository)).On("Update", tt.args.inputProduct).Return(tt.want, utils.ErrConflict)
			} else {
				(tt.fields.repo.(*mockProductRepository)).On("Update", tt.args.inputProduct).Return(tt.want, nil)
			}
			// Mock the GetProductTypeByID method
			(tt.fields.validationProductType.(*mockProductTypeValidation)).On("GetProductTypeByID", tt.args.inputProduct.ProductType).Return(internal.ProductType{}, nil)
			// Mock the GetByID method
			(tt.fields.validationSeller.(*mockSellerValidation)).On("GetByID", tt.args.inputProduct.SellerID).Return(internal.Seller{}, nil)
			// Mock the GetALL method
			(tt.fields.repo.(*mockProductRepository)).On("GetAll").Return([]internal.Product{
				{
					ID: 1,
					ProductAttributes: internal.ProductAttributes{
						ProductCode:                    "123",
						Description:                    "Product 1",
						Width:                          1.0,
						Height:                         1.0,
						Length:                         1.0,
						NetWeight:                      1.0,
						ExpirationRate:                 1.0,
						RecommendedFreezingTemperature: 1.0,
						FreezingRate:                   1.0,
						ProductType:                    1,
						SellerID:                       1,
					},
				},
			}, nil)

			got, err := s.UpdateProduct(tt.args.inputProduct)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicProductServiceUpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BasicProductService.UpdateProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
