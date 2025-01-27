package locality_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/locality"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
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
	var internalServerError = errors.New("internal server error")
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
			TestName: "given an existing locality ID, return internal server error",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, internalServerError)
			},
			ErrorToReturn: internalServerError,
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
			TestName: "given a valid locality, countryRepo.GetByName, return internal server error",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{}, internalServerError)
			},
			ErrorToReturn: internalServerError,
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
			TestName: "given a valid locality, when countryRepo.Save, return internal server error",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{}, utils.ErrNotFound)
				mcr.On("Save", mock.Anything).Return(internalServerError)
			},
			ErrorToReturn: internalServerError,
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
			TestName: "given a valid locality, when provinceRepo.GetByName, return internal server error",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{
					ID:          192,
					CountryName: "Brazil",
				}, nil)
				mcr.On("Save", mock.Anything).Return(nil)
				mpr.On("GetByName", mock.Anything).Return(internal.Province{}, internalServerError)
			},
			ErrorToReturn: internalServerError,
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
			TestName: "given a valid locality, when provinceRepo.Save, return internal server error",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{
					ID:          192,
					CountryName: "Brazil",
				}, nil)
				mcr.On("Save", mock.Anything).Return(nil)
				mpr.On("GetByName", mock.Anything).Return(internal.Province{}, utils.ErrNotFound)
				mpr.On("Save", mock.Anything).Return(internalServerError)
			},
			ErrorToReturn: internalServerError,
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
			TestName: "given a valid locality, when localityRepo.Save, return nil",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{
					ID:          192,
					CountryName: "Brazil",
				}, nil)
				mcr.On("Save", mock.Anything).Return(nil)
				mpr.On("GetByName", mock.Anything).Return(internal.Province{
					ID:           333,
					ProvinceName: "Westails",
					CountryID:    1,
				}, nil)
				mlr.On("Save", mock.Anything).Return(internalServerError)
			},
			ErrorToReturn: internalServerError,
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
			TestName: "given a valid locality, when localityRepo.Save, return nil",
			Mock: func(mlr *MockLocalityRepository, mpr *MockProvinceRepository, mcr *MockCountryRepository) {
				mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
				mcr.On("GetByName", mock.Anything).Return(internal.Country{
					ID:          192,
					CountryName: "Brazil",
				}, nil)
				mcr.On("Save", mock.Anything).Return(nil)
				mpr.On("GetByName", mock.Anything).Return(internal.Province{
					ID:           333,
					ProvinceName: "Westails",
					CountryID:    1,
				}, nil)
				mlr.On("Save", mock.Anything).Return(nil)
			},
			ErrorToReturn: nil,
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
			service := locality.NewBasicLocalityService(lr, pr, cr)
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

	t.Run("given an invalid locality ID, return an empty report and err of type utils.ErrInvalidArguments", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		service := locality.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetSellersByLocalityID(-1)
		require.ErrorIs(t, err, utils.ErrInvalidArguments)
		require.Len(t, report, 0)
	})

	t.Run("given a valid and existing locality ID, return the report row", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		lr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
		lr.On("GetSellersByLocalityID", mock.Anything).Return([]internal.SellersByLocality{sampleSellerByLocality}, nil)
		service := locality.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetSellersByLocalityID(1)
		require.NoError(t, err)
		require.Len(t, report, 1)
	})

	t.Run("given a valid and not existing locality ID, return an empty report and utils.ErrNotFound", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		lr.On("GetByID", mock.Anything).Return(internal.Locality{}, utils.ErrNotFound)
		service := locality.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetSellersByLocalityID(99)
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Len(t, report, 0)
	})
}

func TestUnitLocality_GetCarriesByLocalityID(t *testing.T) {
	sampleCarriesByLocality := internal.CarriesByLocality{
		LocalityID:   1,
		LocalityName: "A random locality",
		CarriesCount: 2,
	}

	t.Run("given an invalid locality ID, return an empty report and err of type utils.ErrInvalidArguments", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		service := locality.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetCarriesByLocalityID(-1)
		require.ErrorIs(t, err, utils.ErrInvalidArguments)
		require.Len(t, report, 0)
	})

	t.Run("given a valid and existing locality ID, return the report row", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		lr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
		lr.On("GetCarriesByLocalityID", mock.Anything).Return([]internal.CarriesByLocality{sampleCarriesByLocality}, nil)
		service := locality.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetCarriesByLocalityID(1)
		require.NoError(t, err)
		require.Len(t, report, 1)
	})

	t.Run("given a valid and not existing locality ID, return an empty report and utils.ErrNotFound", func(t *testing.T) {
		lr := new(MockLocalityRepository)
		pr := new(MockProvinceRepository)
		cr := new(MockCountryRepository)
		lr.On("GetByID", mock.Anything).Return(internal.Locality{}, utils.ErrNotFound)
		service := locality.NewBasicLocalityService(lr, pr, cr)

		report, err := service.GetCarriesByLocalityID(99)
		require.ErrorIs(t, err, utils.ErrNotFound)
		require.Len(t, report, 0)
	})
}
