package buyer

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// Implementation of repository mock
type MockBuyerRepository struct {
	mock.Mock
}

func (m *MockBuyerRepository) GetAll() ([]internal.Buyer, error) {
	args := m.Called()
	return args.Get(0).([]internal.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) GetOne(id int) (*internal.Buyer, error) {
	args := m.Called(id)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) CreateBuyer(buyer internal.Buyer) (*internal.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) UpdateBuyer(buyer *internal.Buyer) (*internal.Buyer, error) {
	args := m.Called(buyer)
	return args.Get(0).(*internal.Buyer), args.Error(1)
}

func (m *MockBuyerRepository) DeleteBuyer(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestBuyerService_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		want    []internal.Buyer
		wantErr error
	}{
		{
			name: "*Get All* Success to return all buyers",
			want: []internal.Buyer{
				{
					ID: 0,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "9762",
						FirstName:    "Ronaldo",
						LastName:     "Messi",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "*Get All* Fail to return all buyers",
			want: []internal.Buyer{
				{
					ID: 0,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "9762",
						FirstName:    "Ronaldo",
						LastName:     "Messi",
					},
				},
			},
			wantErr: errors.New("error getting all buyers"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBuyerRepository{}
			repo.On("GetAll").Return(tt.want, tt.wantErr)
			defer repo.AssertExpectations(t)

			service := NewBuyer(repo)
			results, err := service.GetAll()

			require.Equal(t, tt.want, results)
			require.Equal(t, tt.wantErr, err)

			repo.AssertNumberOfCalls(t, "GetAll", 1)
		})
	}
}

func TestBuyerService_GetOne(t *testing.T) {
	tests := []struct {
		name     string
		id       int
		repo     []internal.Buyer
		wantResp *internal.Buyer
		wantErr  error
	}{
		{
			name:     "*Get One* Success to return a buyer",
			id:       1,
			repo:     []internal.Buyer{{ID: 1}},
			wantResp: &internal.Buyer{ID: 1},
			wantErr:  nil,
		},
		{
			name:     "*Get One* Fail to return a buyer",
			id:       1,
			repo:     nil,
			wantResp: nil,
			wantErr:  errors.New("error getting buyer"),
		},
		{
			name:     "*Get One* Fail to return a buyer with ID different",
			id:       1,
			repo:     []internal.Buyer{{ID: 2}},
			wantResp: nil,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBuyerRepository{}
			repo.On("GetAll").Return(tt.repo, tt.wantErr)
			defer repo.AssertExpectations(t)

			service := NewBuyer(repo)
			results, err := service.GetOne(tt.id)
			require.Equal(t, tt.wantResp, results)
			require.Equal(t, tt.wantErr, err)

			repo.AssertNumberOfCalls(t, "GetAll", 1)
		})
	}
}

func TestBuyerService_CreateBuyer(t *testing.T) {
	tests := []struct {
		name                string
		buyerAttributes     internal.BuyerAttributes
		repo                []internal.Buyer
		newBuyer            internal.Buyer
		want                *internal.Buyer
		wantErr             error
		wantErrValidation   error
		numberOfCallsAll    int
		numberOfCallsCreate int
	}{
		{
			name: "*Create Buyer* Success to create a buyer",
			buyerAttributes: internal.BuyerAttributes{
				CardNumberID: "282948",
				FirstName:    "Ronaldo",
				LastName:     "Messi",
			},
			repo: []internal.Buyer{},
			newBuyer: internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "282948",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			want: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "282948",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			wantErr:             nil,
			numberOfCallsCreate: 1,
			numberOfCallsAll:    2,
		},
		{
			name:                "*Create Buyer* Fail to create a buyer - Repository error (Get All)",
			buyerAttributes:     internal.BuyerAttributes{},
			repo:                nil,
			newBuyer:            internal.Buyer{},
			want:                nil,
			wantErr:             errors.New("error getting all buyers"),
			numberOfCallsCreate: 0,
			numberOfCallsAll:    1,
		},

		{
			name: "*Create Buyer* Fail to create a buyer - validation error",
			buyerAttributes: internal.BuyerAttributes{
				CardNumberID: "282948",
				FirstName:    "Ronaldo",
				LastName:     "Messi",
			},
			repo: []internal.Buyer{
				{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "282948",
						FirstName:    "Ronaldo",
						LastName:     "Messi",
					},
				},
			},
			newBuyer: internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "282948",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			want:                nil,
			wantErr:             nil,
			wantErrValidation:   errors.New("entity already exists"),
			numberOfCallsCreate: 0,
			numberOfCallsAll:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBuyerRepository{}
			repo.On("GetAll").Return(tt.repo, tt.wantErr)
			if tt.want != nil {
				repo.On("CreateBuyer", tt.newBuyer).Return(tt.want, tt.wantErr)
			}
			defer repo.AssertExpectations(t)

			service := NewBuyer(repo)
			results, err := service.CreateBuyer(tt.buyerAttributes)

			require.Equal(t, tt.want, results)

			if tt.wantErrValidation == nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.Equal(t, tt.wantErrValidation, err)
			}

			repo.AssertNumberOfCalls(t, "GetAll", tt.numberOfCallsAll)
			repo.AssertNumberOfCalls(t, "CreateBuyer", tt.numberOfCallsCreate)
		})
	}
}

func TestBuyerService_UpdateBuyer(t *testing.T) {
	tests := []struct {
		name         string
		repo         []internal.Buyer
		updatedBuyer *internal.Buyer
		want         *internal.Buyer
		wantErr      error
	}{
		{
			name: "*Update Buyer* Success to update a buyer",
			repo: []internal.Buyer{{ID: 1}},
			updatedBuyer: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "2123",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			want:    &internal.Buyer{ID: 1},
			wantErr: nil,
		},
		{
			name: "*Update Buyer* Fail to update a buyer - Repository error (Get All)",
			repo: []internal.Buyer{},
			updatedBuyer: &internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "2123",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			want:    nil,
			wantErr: errors.New("entity not found"),
		},
		{
			name: "*Update Buyer* Fail to update a buyer - Repository error (Get All)",
			repo: []internal.Buyer{
				{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "2123",
						FirstName:    "Ronaldo",
						LastName:     "Messi",
					},
				},
				{
					ID: 2,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "2125",
						FirstName:    "Ronaldo",
						LastName:     "Messi",
					},
				},
			},
			updatedBuyer: &internal.Buyer{
				ID: 2,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "2123",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			want:    nil,
			wantErr: errors.New("entity already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBuyerRepository{}
			repo.On("GetAll").Return(tt.repo, tt.wantErr)
			if tt.want != nil {
				repo.On("UpdateBuyer", tt.updatedBuyer).Return(tt.want, tt.wantErr)
			}
			defer repo.AssertExpectations(t)

			service := NewBuyer(repo)
			results, err := service.UpdateBuyer(tt.updatedBuyer)

			require.Equal(t, tt.want, results)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBuyerService_DeleteBuyer(t *testing.T) {
	tests := []struct {
		name              string
		repo              *internal.Buyer
		id                int
		wantErr           error
		wantErrValidation error
	}{
		{
			name:              "*Delete Buyer* Success to delete a buyer",
			repo:              &internal.Buyer{ID: 1},
			id:                1,
			wantErr:           nil,
			wantErrValidation: nil,
		},
		{
			name:              "*Delete Buyer* Fail to delete a buyer - Repository error (Get One)",
			repo:              nil,
			id:                2,
			wantErr:           errors.New("entity not found"),
			wantErrValidation: errors.New("entity not found"),
		},
		{
			name:              "*Delete Buyer* Fail to delete a buyer - Repository error (DeleteBuyer)",
			repo:              &internal.Buyer{ID: 1},
			id:                2,
			wantErr:           nil,
			wantErrValidation: errors.New("entity not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBuyerRepository{}
			repo.On("GetOne", tt.id).Return(tt.repo, tt.wantErr)
			if tt.wantErr == nil {
				repo.On("DeleteBuyer", tt.id).Return(tt.wantErrValidation)
			}
			defer repo.AssertExpectations(t)

			service := NewBuyer(repo)
			err := service.DeleteBuyer(tt.id)

			if tt.wantErrValidation == nil {
				require.Equal(t, tt.wantErr, err)
			}
			require.Equal(t, tt.wantErrValidation, err)
		})
	}
}

func TestBuyerService_validation(t *testing.T) {
	tests := []struct {
		name              string
		repo              []internal.Buyer
		newBuyer          internal.Buyer
		wantErr           error
		wantErrValidation error
	}{
		{
			name: "Success to validate a buyer",
			repo: []internal.Buyer{
				{
					ID: 1,
					BuyerAttributes: internal.BuyerAttributes{
						CardNumberID: "2123",
						FirstName:    "Ronaldo",
						LastName:     "Messi",
					},
				},
			},
			newBuyer: internal.Buyer{
				ID: 2,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "2124",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			wantErr:           nil,
			wantErrValidation: nil,
		},
		{
			name: "Fail to validate a buyer - Repository error (Get All)",
			repo: []internal.Buyer{},
			newBuyer: internal.Buyer{
				ID: 2,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "2124",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			wantErr:           errors.New("error getting all buyers"),
			wantErrValidation: nil,
		},
		{
			name: "Fail to validate a buyer - Same ID",
			repo: []internal.Buyer{{ID: 1}},
			newBuyer: internal.Buyer{
				ID: 1,
				BuyerAttributes: internal.BuyerAttributes{
					CardNumberID: "2124",
					FirstName:    "Ronaldo",
					LastName:     "Messi",
				},
			},
			wantErr:           nil,
			wantErrValidation: errors.New("entity already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockBuyerRepository{}
			repo.On("GetAll").Return(tt.repo, tt.wantErr)
			defer repo.AssertExpectations(t)

			service := NewBuyer(repo)
			err := service.validation(tt.newBuyer)

			if tt.wantErrValidation == nil {
				require.Equal(t, tt.wantErr, err)
			} else {
				require.Equal(t, tt.wantErrValidation, err)
			}
		})

	}
}
