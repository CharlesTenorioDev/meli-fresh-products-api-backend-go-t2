package service_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLocalityRepository struct {
	mock.Mock
}

func (m *MockLocalityRepository) Save(locality *internal.Locality) error {
	args := m.Called(locality)
	return args.Error(0)
}
func (m *MockLocalityRepository) GetByID(id int) (internal.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Locality), args.Error(1)
}
func (m *MockLocalityRepository) GetSellersByLocalityID(localityId int) ([]internal.SellersByLocality, error) {
	args := m.Called(localityId)
	return args.Get(0).([]internal.SellersByLocality), args.Error(1)
}

func (m *MockLocalityRepository) GetCarriesByLocalityID(localityId int) ([]internal.CarriesByLocality, error) {
	args := m.Called(localityId)
	return args.Get(0).([]internal.CarriesByLocality), args.Error(1)
}

type MockProvinceRepository struct {
	mock.Mock
}

func (m *MockProvinceRepository) Save(province *internal.Province) error {
	args := m.Called(province)
	return args.Error(0)
}
func (m *MockProvinceRepository) GetByName(name string) (internal.Province, error) {
	args := m.Called(name)
	return args.Get(0).(internal.Province), args.Error(1)
}

type MockCountryRepository struct {
	mock.Mock
}

func (m *MockCountryRepository) Save(country *internal.Country) error {
	args := m.Called(country)
	return args.Error(0)
}
func (m *MockCountryRepository) GetByName(name string) (internal.Country, error) {
	args := m.Called(name)
	return args.Get(0).(internal.Country), args.Error(1)
}

func TestUnitLocality_Save(t *testing.T) {

	cases := []struct {
		TestName      string
		Mock          func(*MockLocalityRepository, *MockProvinceRepository, *MockCountryRepository)
		ErrorToReturn error
		DataLocality  internal.Locality
		DataProvince  internal.Province
		DataCountry   internal.Country
	}{
		{
			TestName:      "given an invalid LocalityName, return utils.ErrInvalidArguments",
			Mock:          func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {},
			ErrorToReturn: utils.ErrInvalidArguments,
			DataLocality:  internal.Locality{},
		},
		{
			TestName:      "given an invalid Locality ID, return utils.ErrInvalidArguments",
			Mock:          func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {},
			ErrorToReturn: utils.ErrInvalidArguments,
			DataLocality: internal.Locality{
				LocalityName: "A random locality",
			},
		},
		{
			TestName:      "given an invalid ProvinceName, return utils.ErrInvalidArguments",
			Mock:          func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {},
			ErrorToReturn: utils.ErrInvalidArguments,
			DataLocality: internal.Locality{
				ID:           7000,
				LocalityName: "A random locality",
			},
			DataProvince: internal.Province{},
		},
		{
			TestName:      "given an invalid ProvinceName, return utils.ErrInvalidArguments",
			Mock:          func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {},
			ErrorToReturn: utils.ErrInvalidArguments,
			DataLocality: internal.Locality{
				ID:           7000,
				LocalityName: "A random locality",
			},
			DataProvince: internal.Province{
				ProvinceName: "Westails",
			},
		},
		{
			TestName:      "given an invalid CountryName, return utils.ErrInvalidArguments",
			Mock:          func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {},
			ErrorToReturn: utils.ErrInvalidArguments,
			DataLocality: internal.Locality{
				ID:           7000,
				LocalityName: "A random locality",
			},
			DataProvince: internal.Province{
				ProvinceName: "Westails",
			},
			DataCountry: internal.Country{},
		},
		{
			TestName: "given an existing locality ID, return utils.ErrConflict",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{ID: 2, LocalityName: "Westails"}, nil)
			},
			ErrorToReturn: utils.ErrConflict,
			DataLocality: internal.Locality{
				ID:           7000,
				LocalityName: "A random locality",
			},
			DataProvince: internal.Province{
				ProvinceName: "Westails",
			},
			DataCountry: internal.Country{
				CountryName: "Ostania",
			},
		},
		{
			TestName: "given an existing locality ID, return ",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{}, nil)
			},
			ErrorToReturn: utils.ErrConflict,
			DataLocality: internal.Locality{
				ID:           7000,
				LocalityName: "A random locality",
			},
			DataProvince: internal.Province{
				ProvinceName: "Westails",
			},
			DataCountry: internal.Country{
				CountryName: "Ostania",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			lr := new(MockLocalityRepository)
			pr := new(MockProvinceRepository)
			cr := new(MockCountryRepository)
			c.Mock(lr, pr, cr)
			service := service.NewBasicLocalityService(lr, pr, cr)
			err := service.Save(&c.DataLocality, &c.DataProvince, &c.DataCountry)
			require.ErrorIs(t, err, c.ErrorToReturn)
		})

	}

}

func TestUnitLocality_GetSellersByLocalityId(t *testing.T) {
	sampleSellerByLocality := internal.SellersByLocality{
		LocalityID:   1,
		LocalityName: "A random locality",
		SellersCount: 2,
	}
	t.Run("given an existing locality ID, return the count", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		lr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
		lr.On("GetSellersByLocalityID", mock.Anything).Return([]internal.SellersByLocality{sampleSellerByLocality}, nil)
		service := service.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetSellersByLocalityID(0)
		require.NoError(t, err)
		require.Len(t, report, 1)
	})
	t.Run("given a not existing locality ID, return an empty report and utils.", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		lr.On("GetByID", mock.Anything).Return(internal.Locality{}, utils.ErrNotFound)
		service := service.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetSellersByLocalityID(99)
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Len(t, report, 0)
	})
}
