package locality

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocalityService is a mock implementation of the LocalityService interface
type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) GetSellersByLocalityID(id int) ([]internal.SellersByLocality, error) {
	args := m.Called(id)
	return args.Get(0).([]internal.SellersByLocality), args.Error(1)
}

func (m *MockLocalityService) CreateLocality(locality internal.Locality) (internal.Locality, error) {
	args := m.Called(locality)
	return args.Get(0).(internal.Locality), args.Error(1)
}

func (m *MockLocalityService) GetCarriesByLocalityID(id int) ([]internal.CarriesByLocality, error) {
	args := m.Called(id)
	return args.Get(0).([]internal.CarriesByLocality), args.Error(1)
}

func (m *MockLocalityService) Save(*internal.Locality, *internal.Province, *internal.Country) error {
	args := m.Called()
	return args.Error(0)
}

func TestNewLocalityRoutes(t *testing.T) {
	mux := chi.NewMux()
	mockService := new(MockLocalityService)

	err := NewLocalityRoutes(mux, mockService)
	assert.NoError(t, err)

	tests := []struct {
		method       string
		target       string
		expectedCode int
	}{
		{"GET", "/api/v1/localities/reportSellers", http.StatusOK},
		{"POST", "/api/v1/localities", http.StatusOK},
		{"GET", "/api/v1/localities/reportCarries", http.StatusOK},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.target, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		assert.Equal(t, tt.expectedCode, rr.Code)
	}
}
