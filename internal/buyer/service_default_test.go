package buyer

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type BuyerRepositoryMock struct {
	mock.Mock
}

func (m *BuyerRepositoryMock) GetAll() ([]internal.Buyer, error) {
	args := m.Called()
	return args.Get(0).([]internal.Buyer), args.Error(1)
}

func (m *BuyerRepositoryMock) GetOne(id int) (*internal.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *BuyerRepositoryMock) CreateBuyer(buyer internal.Buyer) (*internal.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *BuyerRepositoryMock) UpdateBuyer(buyer *internal.Buyer) (*internal.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *BuyerRepositoryMock) DeleteBuyer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBuyerService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockRepo    func() *BuyerRepositoryMock
		expected    []internal.Buyer
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Get all buyers",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				return m
			},
			expected: []internal.Buyer{
				{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "123456789",
						FirstName:    "John",
						LastName:     "Doe",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Error to get all buyers",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{}, utils.ErrNotFound)
				return m
			},
			expected:    []internal.Buyer{},
			wantErr:     true,
			expectedErr: utils.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			service := NewBuyer(repo)

			got, err := service.GetAll()
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestBuyerService_GetOne(t *testing.T) {
	tests := []struct {
		name        string
		mockRepo    func() *BuyerRepositoryMock
		id          int
		expected    *internal.Buyer
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Get one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				return m
			},
			id: 1,
			expected: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "123456789",
					FirstName:    "John",
					LastName:     "Doe",
				},
			},
			wantErr: false,
		},
		{
			name: "Error to get one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{}, utils.ErrNotFound)
				return m
			},
			id:          1,
			expected:    nil,
			wantErr:     true,
			expectedErr: utils.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			service := NewBuyer(repo)

			got, err := service.GetOne(tt.id)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestUnitBuyerService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		mockRepo    func() *BuyerRepositoryMock
		id          int
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Delete one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				m.On("DeleteBuyer", 1).Return(nil)
				return m
			},
			id:      1,
			wantErr: false,
		},
		{
			name: "Error to delete one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				m.On("DeleteBuyer", 1).Return(utils.ErrNotFound)
				return m
			},
			id:          1,
			wantErr:     true,
			expectedErr: utils.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			service := NewBuyer(repo)

			err := service.DeleteBuyer(tt.id)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUnitBuyerService_Create(t *testing.T) {
	tests := []struct {
		name          string
		mockRepo      func() *BuyerRepositoryMock
		buyer         internal.BuyerAttributes
		expectedBuyer *internal.Buyer
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Create one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{}, nil)
				m.On("CreateBuyer", mock.Anything).Return(&internal.Buyer{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "123456789",
						FirstName:    "John",
						LastName:     "Doe",
					},
				}, nil)
				return m
			},
			buyer: internal.BuyerAttributes{
				CardNumberID: "123456789",
				FirstName:    "John",
				LastName:     "Doe",
			},
			wantErr: false,
			expectedBuyer: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "123456789",
					FirstName:    "John",
					LastName:     "Doe",
				},
			},
		},
		{
			name: "Error to create one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{}, nil)
				m.On("CreateBuyer", mock.Anything).Return(&internal.Buyer{}, utils.ErrNotFound)
				return m
			},
			buyer: internal.BuyerAttributes{
				CardNumberID: "123456789",
				FirstName:    "John",
				LastName:     "Doe",
			},
			wantErr:       true,
			expectedErr:   utils.ErrNotFound,
			expectedBuyer: &internal.Buyer{},
		},
		{
			name: "Error to create one buyer - CardNumberID already exists",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "123456789",
						FirstName:    "John",
						LastName:     "Doe",
					},
				}}, nil)
				m.On("CreateBuyer", mock.Anything).Return(&internal.Buyer{}, utils.ErrConflict)
				return m
			},
			buyer: internal.BuyerAttributes{
				CardNumberID: "123456789",
				FirstName:    "John",
				LastName:     "Doe",
			},
			wantErr:       true,
			expectedErr:   utils.ErrConflict,
			expectedBuyer: &internal.Buyer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			service := NewBuyer(repo)

			_, err := service.CreateBuyer(tt.buyer)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUnitBuyerService_Update(t *testing.T) {
	tests := []struct {
		name          string
		mockRepo      func() *BuyerRepositoryMock
		buyer         internal.Buyer
		expectedBuyer *internal.Buyer
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Update one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				m.On("UpdateBuyer", mock.Anything).Return(&internal.Buyer{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "987654321",
						FirstName:    "Jane",
						LastName:     "Doe",
					},
				}, nil)
				return m
			},
			buyer: internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "987654321",
					FirstName:    "Jane",
					LastName:     "Doe",
				},
			},
			wantErr: false,
			expectedBuyer: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "987654321",
					FirstName:    "Jane",
					LastName:     "Doe",
				},
			},
		},
		{
			name: "Error to update one buyer",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				m.On("UpdateBuyer", mock.Anything).Return(&internal.Buyer{}, utils.ErrNotFound)
				return m
			},
			buyer: internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "987654321",
					FirstName:    "Jane",
					LastName:     "Doe",
				},
			},
			wantErr:       true,
			expectedErr:   utils.ErrNotFound,
			expectedBuyer: &internal.Buyer{},
		},
		{
			name: "Error to update one buyer - Card number already exist",
			mockRepo: func() *BuyerRepositoryMock {
				m := &BuyerRepositoryMock{}
				m.On("GetAll").Return([]internal.Buyer{
					{
						ID: 1,
						BuyerAttributes: internal.BuyerAttributes{
							CardNumberID: "123456789",
							FirstName:    "John",
							LastName:     "Doe",
						},
					},
				}, nil)
				m.On("UpdateBuyer", mock.Anything).Return(&internal.Buyer{}, utils.ErrConflict)
				return m
			},
			buyer: internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "123456789",
					FirstName:    "Jane",
					LastName:     "Doe",
				},
			},
			wantErr:       true,
			expectedErr:   utils.ErrConflict,
			expectedBuyer: &internal.Buyer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.mockRepo()
			service := NewBuyer(repo)

			_, err := service.UpdateBuyer(&tt.buyer)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
