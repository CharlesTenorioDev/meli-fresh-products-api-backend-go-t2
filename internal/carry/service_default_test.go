package carry_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/carry"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCarryRepository struct {
	mock.Mock
}

type MockLocalityValidation struct {
	mock.Mock
}

func (m *MockCarryRepository) Save(carry *internal.Carry) error {
	args := m.Called(carry)
	return args.Error(0)
}

func (m *MockCarryRepository) GetAll() ([]internal.Carry, error) {
	args := m.Called()
	return args.Get(0).([]internal.Carry), args.Error(1)
}

func (m *MockCarryRepository) GetByID(id int) (internal.Carry, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Carry), args.Error(1)
}

func (m *MockCarryRepository) Update(carry *internal.Carry) error {
	args := m.Called(carry)
	return args.Error(0)
}

func (m *MockCarryRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLocalityValidation) GetByID(id int) (internal.Locality, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Locality), args.Error(1)
}

func TestMySQLCarryService_GetAll(t *testing.T) {
	tests := []struct {
		name        string
		mockRepo    func() (*MockCarryRepository, *MockLocalityValidation)
		expected    []internal.Carry
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Get all carries",
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("GetAll").Return([]internal.Carry{
					{
						CID:         1,
						CompanyName: "Company 1",
						Address:     "Address 1",
						Telephone:   "123456789",
						LocalityID:  1,
					},
				}, nil)
				localityValidation := new(MockLocalityValidation)
				return repo, localityValidation
			},
			expected: []internal.Carry{
				{
					CID:         1,
					CompanyName: "Company 1",
					Address:     "Address 1",
					Telephone:   "123456789",
					LocalityID:  1,
				},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Get all carries - Error",
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("GetAll").Return([]internal.Carry{}, utils.ENotFound("Carry"))
				localityValidation := new(MockLocalityValidation)
				return repo, localityValidation
			},
			expected:    []internal.Carry{},
			wantErr:     true,
			expectedErr: utils.ENotFound("Carry"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, localityValidation := tt.mockRepo()
			s := carry.NewMySQLCarryService(repo, localityValidation)

			got, err := s.GetAll()
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestMySQLCarryService_GetByID(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		mockRepo    func() (*MockCarryRepository, *MockLocalityValidation)
		expected    internal.Carry
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Get carry by ID",
			id:   1,
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("GetByID", 1).Return(internal.Carry{
					CID:         1,
					CompanyName: "Company 1",
					Address:     "Address 1",
					Telephone:   "123456789",
					LocalityID:  1,
				}, nil)
				localityValidation := new(MockLocalityValidation)
				return repo, localityValidation
			},
			expected: internal.Carry{
				CID:         1,
				CompanyName: "Company 1",
				Address:     "Address 1",
				Telephone:   "123456789",
				LocalityID:  1,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "Get carry by ID - Error",
			id:   1,
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("GetByID", 1).Return(internal.Carry{}, utils.ENotFound("Carry"))
				localityValidation := new(MockLocalityValidation)
				return repo, localityValidation
			},
			expected:    internal.Carry{},
			wantErr:     true,
			expectedErr: utils.ENotFound("Carry"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, localityValidation := tt.mockRepo()
			s := carry.NewMySQLCarryService(repo, localityValidation)

			got, err := s.GetByID(tt.id)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestMySQLCarryService_Save(t *testing.T) {
	tests := []struct {
		name         string
		carry        *internal.Carry
		mockRepo     func() (*MockCarryRepository, *MockLocalityValidation)
		expectedBody internal.Carry
		wantErr      bool
		expectedErr  error
		Assert       func(*testing.T, *MockCarryRepository, *MockLocalityValidation)
	}{
		{
			name: "Save carry",
			carry: &internal.Carry{
				CID:         1,
				CompanyName: "Company 1",
				Address:     "Address 1",
				Telephone:   "123456789",
				LocalityID:  1,
			},
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("Save", &internal.Carry{
					ID:          0,
					CID:         1,
					CompanyName: "Company 1",
					Address:     "Address 1",
					Telephone:   "123456789",
					LocalityID:  1,
				}).Return(nil)
				localityValidation := new(MockLocalityValidation)
				localityValidation.On("GetByID", 1).Return(internal.Locality{}, nil)
				return repo, localityValidation
			},
			wantErr:     false,
			expectedErr: nil,
			expectedBody: internal.Carry{
				ID:          0,
				CID:         1,
				CompanyName: "Company 1",
				Address:     "Address 1",
				Telephone:   "123456789",
				LocalityID:  1,
			},
			Assert: func(t *testing.T, repo *MockCarryRepository, localityValidation *MockLocalityValidation) {
				repo.AssertNumberOfCalls(t, "Save", 1)
				localityValidation.AssertNumberOfCalls(t, "GetByID", 1)
			},
		},
		{
			name:  "Save carry - Error When CID is Zero",
			carry: &internal.Carry{},
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("Save", &internal.Carry{}).Return(&internal.Carry{}, utils.EBadRequest("Carry"))
				localityValidation := new(MockLocalityValidation)
				localityValidation.On("GetByID", 0).Return(internal.Locality{}, nil)
				return repo, localityValidation
			},
			wantErr:      true,
			expectedErr:  utils.EZeroValue("CID"),
			expectedBody: internal.Carry{},
			Assert: func(t *testing.T, repo *MockCarryRepository, localityValidation *MockLocalityValidation) {
				repo.AssertNumberOfCalls(t, "Save", 0)
				localityValidation.AssertNumberOfCalls(t, "GetByID", 1)
			},
		},
		{
			name: "Save carry - Error When Address is Zero",
			carry: &internal.Carry{
				CID:         1,
				CompanyName: "Company 1",
			},
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("Save", &internal.Carry{}).Return(&internal.Carry{}, utils.EBadRequest("Carry"))
				localityValidation := new(MockLocalityValidation)
				localityValidation.On("GetByID", 0).Return(internal.Locality{}, nil)
				return repo, localityValidation
			},
			wantErr:      true,
			expectedErr:  utils.EZeroValue("Address"),
			expectedBody: internal.Carry{},
			Assert: func(t *testing.T, repo *MockCarryRepository, localityValidation *MockLocalityValidation) {
				repo.AssertNumberOfCalls(t, "Save", 0)
				localityValidation.AssertNumberOfCalls(t, "GetByID", 1)
			},
		},
		{
			name: "Save carry - Error When Company Name is Zero",
			carry: &internal.Carry{
				CID:     1,
				Address: "Address 1",
			},
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("Save", &internal.Carry{}).Return(&internal.Carry{}, utils.EBadRequest("Carry"))
				localityValidation := new(MockLocalityValidation)
				localityValidation.On("GetByID", 0).Return(internal.Locality{}, nil)
				return repo, localityValidation
			},
			wantErr:      true,
			expectedErr:  utils.EZeroValue("CompanyName"),
			expectedBody: internal.Carry{},
			Assert: func(t *testing.T, repo *MockCarryRepository, localityValidation *MockLocalityValidation) {
				repo.AssertNumberOfCalls(t, "Save", 0)
				localityValidation.AssertNumberOfCalls(t, "GetByID", 1)
			},
		},
		{
			name: "Save carry - Error when Telephone is Zero",
			carry: &internal.Carry{
				CID:         1,
				Address:     "Address 1",
				CompanyName: "Company 1",
			},
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				repo := new(MockCarryRepository)
				repo.On("Save", &internal.Carry{}).Return(&internal.Carry{}, utils.EBadRequest("Carry"))
				localityValidation := new(MockLocalityValidation)
				localityValidation.On("GetByID", 0).Return(internal.Locality{}, nil)
				return repo, localityValidation
			},
			wantErr:      true,
			expectedErr:  utils.EZeroValue("Telephone"),
			expectedBody: internal.Carry{},
			Assert: func(t *testing.T, repo *MockCarryRepository, localityValidation *MockLocalityValidation) {
				repo.AssertNumberOfCalls(t, "Save", 0)
				localityValidation.AssertNumberOfCalls(t, "GetByID", 1)
			},
		},
		{
			name: "Save carry - Error when LocalityID is Zero",
			carry: &internal.Carry{
				CID:         1,
				Address:     "Address 1",
				CompanyName: "Company 1",
				Telephone:   "123456789",
			},
			mockRepo: func() (*MockCarryRepository, *MockLocalityValidation) {
				localityValidation := new(MockLocalityValidation)
				localityValidation.On("GetByID", 0).Return(internal.Locality{}, utils.ENotFound("Locality"))
				repo := new(MockCarryRepository)
				repo.On("Save", mock.Anything).Return(utils.ENotFound("Locality"))

				return repo, localityValidation
			},
			wantErr:      true,
			expectedErr:  utils.ENotFound("Locality"),
			expectedBody: internal.Carry{},
			Assert: func(t *testing.T, repo *MockCarryRepository, localityValidation *MockLocalityValidation) {
				repo.AssertNumberOfCalls(t, "Save", 0)
				localityValidation.AssertNumberOfCalls(t, "GetByID", 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, localityValidation := tt.mockRepo()
			s := carry.NewMySQLCarryService(repo, localityValidation)
			err := s.Save(tt.carry)
			gotCarry := internal.Carry{}
			if err == nil {
				gotCarry.ID = tt.carry.ID
				gotCarry.CID = tt.carry.CID
				gotCarry.CompanyName = tt.carry.CompanyName
				gotCarry.Address = tt.carry.Address
				gotCarry.Telephone = tt.carry.Telephone
				gotCarry.LocalityID = tt.carry.LocalityID
			}
			tt.Assert(t, repo, localityValidation)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedBody, gotCarry)
		})
	}
}
