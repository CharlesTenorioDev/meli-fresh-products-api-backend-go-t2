package warehouse

import (
	"errors"
	"strings"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWarehouseRepository struct {
	mock.Mock
}

type mockLocalityValidation struct {
	mock.Mock
}

func (m *mockLocalityValidation) GetByID(id int) (locality internal.Locality, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.Locality), args.Error(1)
}

func (m *mockWarehouseRepository) GetAll() (listWarehouses []internal.Warehouse, err error) {
	args := m.Called()
	return args.Get(0).([]internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseRepository) GetByID(id int) (warehouse internal.Warehouse, err error) {
	args := m.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseRepository) Save(newWarehouse internal.Warehouse) (warehouse internal.Warehouse, err error) {
	args := m.Called(newWarehouse)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseRepository) Update(updatedWarehouse internal.Warehouse) (warehouse internal.Warehouse, err error) {
	args := m.Called(updatedWarehouse)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (m *mockWarehouseRepository) Delete(id int) (err error) {
	args := m.Called(id)
	return args.Error(0)
}

func TestUnitWarehouse_GetAll(t *testing.T) {
	type fields struct {
		repo               internal.WarehouseRepository
		validationLocality internal.LocalityValidation
	}

	setupMock := func(repo *mockWarehouseRepository, expectedListWarehouses []internal.Warehouse, expectedErr error) {
		if expectedErr != nil {
			repo.On("GetAll").Return(expectedListWarehouses, utils.ErrNotFound)
		} else {
			repo.On("GetAll").Return(expectedListWarehouses, nil)
		}
	}

	tests := []struct {
		name                   string
		fields                 fields
		expectedListWarehouses []internal.Warehouse
		expectedErr            error
	}{
		{
			name: "GetWarehouse OK",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			expectedListWarehouses: []internal.Warehouse{
				{ID: 1, Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20},
			},
			expectedErr: nil,
		},
		{
			name: "GetWarehouses Error",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			expectedListWarehouses: nil,
			expectedErr:            utils.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicWarehouseService{
				repo:             tt.fields.repo,
				validateLocality: tt.fields.validationLocality,
			}

			setupMock(tt.fields.repo.(*mockWarehouseRepository), tt.expectedListWarehouses, tt.expectedErr)

			gotListWarehouses, err := s.GetAll()

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedListWarehouses, gotListWarehouses)
		})
	}
}

func TestUnitWarehouse_GetById(t *testing.T) {
	type fields struct {
		repo               internal.WarehouseRepository
		validationLocality internal.LocalityValidation
	}

	setupMock := func(repo *mockWarehouseRepository, id int, expectedWarehouse internal.Warehouse, expectedErr error) {
		if expectedErr != nil {
			if errors.Is(expectedErr, utils.ErrNotFound) {
				repo.On("GetByID", id).Return(expectedWarehouse, utils.ErrNotFound)
			} else {
				repo.On("GetByID", id).Return(expectedWarehouse, expectedErr)
			}
		} else {
			repo.On("GetByID", id).Return(expectedWarehouse, nil)
		}

	}

	tests := []struct {
		name              string
		fields            fields
		expectedWarehouse internal.Warehouse
		id                int
		expectedErr       error
	}{
		{
			name: "GetWarehouseByID OK",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			expectedWarehouse: internal.Warehouse{
				ID:                 1,
				Address:            "1234 Cold Storage St, LA",
				Telephone:          "555-3456",
				WarehouseCode:      "WH001",
				LocalityID:         1,
				MinimumCapacity:    30,
				MinimumTemperature: 20,
			},
			id:          1,
			expectedErr: nil,
		},
		{
			name: "GetWarehouseByID Error Not Found",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			id:                1,
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrNotFound,
		},
		{
			name: "GetWarehouseByID Internal Server Error",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			id:                1,
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       errors.New("internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicWarehouseService{
				repo:             tt.fields.repo,
				validateLocality: tt.fields.validationLocality,
			}

			setupMock(tt.fields.repo.(*mockWarehouseRepository), tt.id, tt.expectedWarehouse, tt.expectedErr)

			gotListWarehouses, err := s.GetByID(tt.id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedWarehouse, gotListWarehouses)
		})
	}
}

func TestUnitWarehouse_Save(t *testing.T) {
	type fields struct {
		repo               internal.WarehouseRepository
		validationLocality internal.LocalityValidation
	}

	setupMock := func(fields fields, expectedWarehouse internal.Warehouse, newWarehouse internal.Warehouse, expectedErr error) {
		if expectedErr != nil {
			if strings.EqualFold(expectedErr.Error(), utils.EConflict("Warehouse", "Locality").Error()) {
				fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, utils.ErrConflict)
			} else if strings.EqualFold(expectedErr.Error(), utils.EConflict("Warehouse", "Warehouse code").Error()) {
				fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, nil)
				fields.repo.(*mockWarehouseRepository).On("GetAll").Return([]internal.Warehouse{newWarehouse}, nil)
			} else {
				fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, nil)
				fields.repo.(*mockWarehouseRepository).On("GetAll").Return([]internal.Warehouse{expectedWarehouse}, nil)
				fields.repo.(*mockWarehouseRepository).On("Save", newWarehouse).Return(internal.Warehouse{}, errors.New("internal Server Error"))
			}

		} else {
			fields.repo.(*mockWarehouseRepository).On("Save", newWarehouse).Return(expectedWarehouse, nil)
			fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, nil)
			fields.repo.(*mockWarehouseRepository).On("GetAll").Return([]internal.Warehouse{expectedWarehouse}, nil)
		}
	}

	tests := []struct {
		name              string
		fields            fields
		newWarehouse      internal.Warehouse
		expectedWarehouse internal.Warehouse
		expectedErr       error
	}{
		{
			name: "SaveWarehouse OK",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH002", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{
				ID: 1, Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedErr: nil,
		},
		{
			name: "SaveWarehouse Invalid WarehouseCode",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrInvalidArguments,
		},
		{
			name: "SaveWarehouse Invalid Address",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrInvalidArguments,
		},
		{
			name: "SaveWarehouse Invalid Telephone",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrInvalidArguments,
		},
		{
			name: "SaveWarehouse Invalid MinimumCapacity",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: -30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrInvalidArguments,
		},
		{
			name: "SaveWarehouse Invalid MinimumTemperature",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: -275,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrInvalidArguments,
		},
		{
			name: "SaveWarehouse Invalid LocalityID",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 1,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.EConflict("Warehouse", "Locality"),
		},
		{
			name: "SaveWarehouse when the Warehouse code already exists",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH002", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.EConflict("Warehouse", "Warehouse code"),
		},
		{
			name: "SaveWarehouse internal Server error",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH002", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       errors.New("internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicWarehouseService{
				repo:             tt.fields.repo,
				validateLocality: tt.fields.validationLocality,
			}

			setupMock(tt.fields, tt.expectedWarehouse, tt.newWarehouse, tt.expectedErr)

			gotListWarehouses, err := s.Save(tt.newWarehouse)

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedWarehouse, gotListWarehouses)
		})
	}
}

func TestUnitWarehouse_Update(t *testing.T) {
	type fields struct {
		repo               internal.WarehouseRepository
		validationLocality internal.LocalityValidation
	}

	setupMock := func(fields fields, expectedWarehouse internal.Warehouse, idRequest int, newWarehouse internal.Warehouse, expectedErr error) {

		if expectedErr != nil {
			if strings.EqualFold(expectedErr.Error(), utils.ENotFound("Get By id Not Found").Error()) {
				fields.repo.(*mockWarehouseRepository).On("GetByID", idRequest).Return(internal.Warehouse{}, utils.ErrNotFound)
			} else if strings.EqualFold(expectedErr.Error(), utils.ENotFound("Not Found").Error()) {
				fields.repo.(*mockWarehouseRepository).On("GetByID", idRequest).Return(internal.Warehouse{}, nil)
			} else if errors.Is(expectedErr, utils.ErrInvalidArguments) {
				fields.repo.(*mockWarehouseRepository).On("GetByID", idRequest).Return(newMockWarehouse(), nil)
			} else if strings.EqualFold(expectedErr.Error(), utils.EConflict("Warehouse", "Warehouse code").Error()) {
				fields.repo.(*mockWarehouseRepository).On("GetByID", idRequest).Return(newMockWarehouse(), nil)
				fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, nil)
				fields.repo.(*mockWarehouseRepository).On("GetAll").Return([]internal.Warehouse{newWarehouse}, nil)
			} else {
				fields.repo.(*mockWarehouseRepository).On("GetByID", idRequest).Return(newMockWarehouse(), nil)
				fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, nil)
				fields.repo.(*mockWarehouseRepository).On("GetAll").Return([]internal.Warehouse{expectedWarehouse}, nil)
				newWarehouse.ID = idRequest
				fields.repo.(*mockWarehouseRepository).On("Update", newWarehouse).Return(internal.Warehouse{}, errors.New("internal Server Error"))
			}
		} else {

			fields.repo.(*mockWarehouseRepository).On("GetByID", idRequest).Return(expectedWarehouse, nil)

			fields.validationLocality.(*mockLocalityValidation).On("GetByID", newWarehouse.LocalityID).Return(internal.Locality{}, nil)

			fields.repo.(*mockWarehouseRepository).On("GetAll").Return([]internal.Warehouse{expectedWarehouse}, nil)

			newWarehouse.ID = idRequest
			fields.repo.(*mockWarehouseRepository).On("Update", newWarehouse).Return(expectedWarehouse, nil)
		}
	}

	tests := []struct {
		name              string
		fields            fields
		newWarehouse      internal.Warehouse
		expectedWarehouse internal.Warehouse
		idRequest         int
		expectedErr       error
	}{
		{
			name: "SaveWarehouse OK",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			idRequest: 1,
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{
				ID: 1, Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedErr: nil,
		},
		{
			name: "SaveWarehouse Get By id Not Found",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			idRequest: 2,
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ENotFound("Get By id Not Found"),
		},
		{
			name: "SaveWarehouse Not Found",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			idRequest: 1,
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ENotFound("Not Found"),
		},
		{
			name: "SaveWarehouse invalid warehouseCode",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			idRequest: 1,
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.ErrInvalidArguments,
		},
		{
			name: "SaveWarehouse when the Warehouse code already exists",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			idRequest: 1,
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH001", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       utils.EConflict("Warehouse", "Warehouse code"),
		},
		{
			name: "SaveWarehouse internal server error",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
				validationLocality: &mockLocalityValidation{
					mock.Mock{},
				},
			},
			idRequest: 1,
			newWarehouse: internal.Warehouse{
				Address: "1234 Cold Storage St, LA", Telephone: "555-3456", WarehouseCode: "WH005", LocalityID: 1, MinimumCapacity: 30, MinimumTemperature: 20,
			},
			expectedWarehouse: internal.Warehouse{},
			expectedErr:       errors.New("internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicWarehouseService{
				repo:             tt.fields.repo,
				validateLocality: tt.fields.validationLocality,
			}

			setupMock(tt.fields, tt.expectedWarehouse, tt.idRequest, tt.newWarehouse, tt.expectedErr)

			gotListWarehouses, err := s.Update(tt.idRequest, NewWarehousePointers(tt.newWarehouse))

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedWarehouse, gotListWarehouses)
		})
	}
}

func TestUnitWarehouse_Delete(t *testing.T) {
	type fields struct {
		repo               internal.WarehouseRepository
		validationLocality internal.LocalityValidation
	}

	setupMock := func(repo *mockWarehouseRepository, id int, expectedErr error) {
		if expectedErr != nil {
			if errors.Is(expectedErr, utils.ErrNotFound) {
				repo.On("GetByID", id).Return(internal.Warehouse{}, utils.ErrNotFound)
			} else {
				repo.On("GetByID", id).Return(internal.Warehouse{}, nil)
				repo.On("Delete", id).Return(expectedErr)
			}
		} else {
			repo.On("GetByID", id).Return(internal.Warehouse{}, nil)
			repo.On("Delete", id).Return(nil)
		}

	}

	tests := []struct {
		name        string
		fields      fields
		id          int
		expectedErr error
	}{
		{
			name: "DeleteWarehouseByID OK",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			id:          1,
			expectedErr: nil,
		},
		{
			name: "DeleteWarehouseByID Warehouse not found",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			id:          5,
			expectedErr: utils.ENotFound("Warehouse"),
		},
		{
			name: "DeleteWarehouseByID Warehouse internal server error",
			fields: fields{
				repo: &mockWarehouseRepository{
					mock.Mock{},
				},
			},
			id:          5,
			expectedErr: errors.New("internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BasicWarehouseService{
				repo:             tt.fields.repo,
				validateLocality: tt.fields.validationLocality,
			}

			setupMock(tt.fields.repo.(*mockWarehouseRepository), tt.id, tt.expectedErr)

			err := s.Delete(tt.id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func newMockWarehouse() internal.Warehouse {
	return internal.Warehouse{
		ID:        1,
		Address:   "1234 Cold Storage St, LA",
		Telephone: "555-3456", WarehouseCode: "WH001",
		LocalityID:         1,
		MinimumCapacity:    30,
		MinimumTemperature: 20,
	}
}

func NewWarehousePointers(w internal.Warehouse) internal.WarehousePointers {
	return internal.WarehousePointers{
		Address:            &w.Address,
		Telephone:          &w.Telephone,
		WarehouseCode:      &w.WarehouseCode,
		LocalityID:         &w.LocalityID,
		MinimumCapacity:    &w.MinimumCapacity,
		MinimumTemperature: &w.MinimumTemperature,
	}
}

func TestNewWarehouseService(t *testing.T) {
	repo := new(mockWarehouseRepository)
	localityValidation := new(mockLocalityValidation)

	service := NewWarehouseService(repo, localityValidation)
	expectedService := &BasicWarehouseService{
		repo:             repo,
		validateLocality: localityValidation,
	}
	assert.Equal(t, expectedService, service)
}
