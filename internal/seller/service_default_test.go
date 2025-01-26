package seller

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type MockSellerRepository struct {
	mock.Mock
}

func (ms *MockSellerRepository) Create(seller *internal.Seller) error {
	args := ms.Called(seller)
	return args.Error(0)
}

func (ms *MockSellerRepository) Update(seller *internal.Seller) error {
	args := ms.Called(seller)
	return args.Error(0)
}

func (ms *MockSellerRepository) GetAll() (sellers []internal.Seller, err error) {
	args := ms.Called()
	return args.Get(0).([]internal.Seller), args.Error(1)
}

func (ms *MockSellerRepository) GetByID(id int) (seller internal.Seller, err error) {
	args := ms.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (ms *MockSellerRepository) GetByCid(cid int) (seller internal.Seller, err error) {
	args := ms.Called(cid)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (ms *MockSellerRepository) Delete(id int) (err error) {
	args := ms.Called(id)
	return args.Error(0)
}

// Mock Locality Repository
type MockLocalityRepository struct {
	mock.Mock
}

func (ml *MockLocalityRepository) Save(locality internal.Locality) (id int, err error) {
	args := ml.Called(locality)
	return args.Int(0), args.Error(1)
}

func (ml *MockLocalityRepository) GetByID(id int) (locality internal.Locality, err error) {
	args := ml.Called(id)
	return args.Get(0).(internal.Locality), args.Error(1)
}

func (ml *MockLocalityRepository) GetSellersByLocalityID(localityID int) (sellers []internal.Seller, err error) {
	args := ml.Called(localityID)
	return args.Get(0).([]internal.Seller), args.Error(1)
}

func (ml *MockLocalityRepository) GetCarriesByLocalityID(localityID int) (carries []internal.Carry, err error) {
	args := ml.Called(localityID)
	return args.Get(0).([]internal.Carry), args.Error(1)
}

func TestUnitSeller_GetAll_Success(t *testing.T) {
	sellers := []internal.Seller{
		{ID: 1, Cid: 55, CompanyName: "Company", Address: "Address", Telephone: "1199999999", LocalityID: 1},
		{ID: 2, Cid: 56, CompanyName: "Company2", Address: "Address2", Telephone: "1199999992", LocalityID: 2},
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetAll").Return(sellers, nil)

	service := NewSellerService(msr, mlr)

	result, err := service.GetAll()

	require.NoError(t, err)
	require.Len(t, result, len(sellers))
	require.Equal(t, result, sellers)
}

func TestUnitSeller_GetAll_InternalServerError(t *testing.T) {

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	internalError := errors.New("internal error")
	msr.On("GetAll").Return([]internal.Seller(nil), internalError)

	service := NewSellerService(msr, mlr)
	result, err := service.GetAll()

	require.ErrorIs(t, err, internalError)
	require.Len(t, result, 0)
	require.Equal(t, result, []internal.Seller(nil))
}

func TestUnitSeller_GetByID_Success(t *testing.T) {
	seller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(seller, nil)

	service := NewSellerService(msr, mlr)

	result, err := service.GetByID(seller.ID)

	require.NoError(t, err)
	require.Equal(t, result, seller)
}

func TestUnitSeller_GetByID_NotFound(t *testing.T) {

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(internal.Seller{}, nil)

	service := NewSellerService(msr, mlr)

	result, err := service.GetByID(1)

	require.Equal(t, utils.ENotFound("Seller"), err)
	require.Equal(t, result, internal.Seller{})
}

func TestUnitSeller_GetByID_InternalServerError(t *testing.T) {

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	internalError := errors.New("internal error")
	msr.On("GetByID", mock.Anything).Return(internal.Seller{}, internalError)

	service := NewSellerService(msr, mlr)

	result, err := service.GetByID(1)

	require.ErrorIs(t, err, internalError)
	require.Equal(t, result, internal.Seller{})
}

func TestUnitSeller_Create_Success(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByCid", mock.Anything).Return(internal.Seller{}, nil)
	mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
	msr.On("Create", mock.Anything).Return(nil)

	service := NewSellerService(msr, mlr)

	result := service.Create(&newSeller)

	require.NoError(t, result)
}

func TestUnitSeller_Create_CidAlreadyExists(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByCid", mock.Anything).Return(internal.Seller{Cid: 55}, nil)
	mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
	msr.On("Create", mock.Anything).Return(nil)

	service := NewSellerService(msr, mlr)

	result := service.Create(&newSeller)

	require.Equal(t, utils.EConflict("Cid", "Seller"), result)
}

func TestUnitSeller_Create_EmptyOrInvalidArguments(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByCid", mock.Anything).Return(internal.Seller{Cid: 0}, nil)

	service := NewSellerService(msr, mlr)

	result := service.Create(&newSeller)

	require.ErrorIs(t, result, utils.ErrInvalidArguments)
}

func TestUnitSeller_Create_LocalityDoesNotExist(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	errGetLocality := utils.EDependencyNotFound("Seller", "locality ID")

	msr.On("GetByCid", mock.Anything).Return(internal.Seller{Cid: 0}, nil)
	mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, errGetLocality)

	service := NewSellerService(msr, mlr)

	result := service.Create(&newSeller)

	require.ErrorIs(t, result, errGetLocality)
}

func TestUnitSeller_Create_InternalServerError(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	internalErr := errors.New("internal server error")

	msr.On("GetByCid", mock.Anything).Return(internal.Seller{Cid: 0}, nil)
	mlr.On("GetByID", mock.Anything).Return(internal.Locality{}, nil)
	msr.On("Create", mock.Anything).Return(internalErr)

	service := NewSellerService(msr, mlr)

	result := service.Create(&newSeller)

	require.ErrorIs(t, result, internalErr)
}

func TestUnitSeller_Update_Success(t *testing.T) {
	updatedSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(internal.Seller{ID: 1, LocalityID: 1}, nil)
	msr.On("GetByCid", mock.Anything).Return(internal.Seller{ID: 1, Cid: 55}, nil)
	msr.On("Update", mock.Anything).Return(nil)

	service := NewSellerService(msr, mlr)

	result, err := service.Update(1, &updatedSeller)

	require.NoError(t, err)
	require.Equal(t, result, updatedSeller)
}

func TestUnitSeller_Update_SellerNotFound(t *testing.T) {
	updatedSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(internal.Seller{}, nil)

	service := NewSellerService(msr, mlr)

	result, err := service.Update(1, &updatedSeller)

	require.Equal(t, utils.ENotFound("Seller"), err)
	require.Equal(t, result, internal.Seller{})
}

func TestUnitSeller_Update_CidAlreadyInUseByOtherSeller(t *testing.T) {
	updatedSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(internal.Seller{ID: 1, LocalityID: 1}, nil)
	msr.On("GetByCid", mock.Anything).Return(internal.Seller{ID: 2, Cid: 55}, nil)

	service := NewSellerService(msr, mlr)

	result, err := service.Update(1, &updatedSeller)

	require.Equal(t, err, utils.EConflict("Cid", "Seller"))
	require.Equal(t, result, internal.Seller{})
}

func TestUnitSeller_Update_InternalServerError(t *testing.T) {
	updatedSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	internalErr := errors.New("internal server error")

	msr.On("GetByID", mock.Anything).Return(internal.Seller{ID: 1, LocalityID: 1}, nil)
	msr.On("GetByCid", mock.Anything).Return(internal.Seller{ID: 1, Cid: 55}, nil)
	msr.On("Update", mock.Anything).Return(internalErr)

	service := NewSellerService(msr, mlr)

	result, err := service.Update(1, &updatedSeller)

	require.ErrorIs(t, err, internalErr)
	require.Equal(t, result, internal.Seller{})
}

func TestUnitSeller_Delete_Success(t *testing.T) {

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(internal.Seller{ID: 1}, nil)
	msr.On("Delete", mock.Anything).Return(nil)

	service := NewSellerService(msr, mlr)

	err := service.Delete(1)

	require.NoError(t, err)
}

func TestUnitSeller_Delete_SellerNotFound(t *testing.T) {

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	msr.On("GetByID", mock.Anything).Return(internal.Seller{}, nil)
	msr.On("Delete", mock.Anything).Return(utils.ErrNotFound)

	service := NewSellerService(msr, mlr)

	err := service.Delete(1)

	require.Equal(t, utils.ENotFound("Seller"), err)
}

func TestUnitSeller_Delete_InternalServerError(t *testing.T) {

	msr := new(MockSellerRepository)
	mlr := new(MockLocalityRepository)

	internalErr := errors.New("internal server error")

	msr.On("GetByID", mock.Anything).Return(internal.Seller{ID: 1}, nil)
	msr.On("Delete", mock.Anything).Return(internalErr)

	service := NewSellerService(msr, mlr)

	err := service.Delete(1)

	require.ErrorIs(t, err, internalErr)
}
