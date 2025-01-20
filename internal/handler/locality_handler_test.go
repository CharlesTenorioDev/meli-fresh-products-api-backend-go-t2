package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLocalityService struct {
	mock.Mock
}

func (m *MockLocalityService) Save(locality *internal.Locality, province *internal.Province, country *internal.Country) error {
	args := m.Called(locality, province, country)
	return args.Error(0)
}
func (m *MockLocalityService) GetSellersByLocalityID(localityId int) ([]internal.SellersByLocality, error) {
	args := m.Called(localityId)
	return args.Get(0).([]internal.SellersByLocality), args.Error(1)
}

func (m *MockLocalityService) GetCarriesByLocalityID(localityId int) ([]internal.CarriesByLocality, error) {
	args := m.Called(localityId)
	return args.Get(0).([]internal.CarriesByLocality), args.Error(1)
}

func TestUnitLocality_CreateLocality(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedStatusCode int
	}{
		{
			TestName:           "CREATED",
			Body:               `{"data":{"id":6701,"locality_name":"Lujan","province_name":"Buenos Aires","country_name":"USA"}}`,
			ExpectedStatusCode: 201,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "BAD_REQUEST",
			Body:               `data:{"id":6701,"locality_name":"Lujan","province_name":"Buenos Aires","country_name":"USA"}}`,
			ExpectedStatusCode: 400,
			ErrorToReturn:      nil, // Error is genereated by JSON parsing process
		},
		{
			TestName:           "CONFLICT",
			Body:               `{"data":{"id":6701,"locality_name":"Lujan","province_name":"Buenos Aires","country_name":"USA"}}`,
			ExpectedStatusCode: 409,
			ErrorToReturn:      utils.ErrConflict,
		},
		{
			TestName:           "UNPROCESSABLE_ENTITY",
			Body:               `{"data":{"id":6701,"locality_name":"","province_name":"","country_name":"USA"}}`,
			ExpectedStatusCode: 422,
			ErrorToReturn:      utils.ErrInvalidArguments,
		},
		{
			TestName:           "INTERNAL_SERVER_ERROR",
			Body:               `{"data":{"id":6701,"locality_name":"Lujan","province_name":"Buenos Aires","country_name":"USA"}}`,
			ExpectedStatusCode: 500,
			ErrorToReturn:      errors.New("Internal Server Error"),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(MockLocalityService)
			service.On("Save", mock.Anything, mock.Anything, mock.Anything).Return(c.ErrorToReturn)
			handler := handler.NewLocalityHandler(service)

			request := &http.Request{
				Body:   io.NopCloser(strings.NewReader(c.Body)),
				Header: http.Header{"Content-Type": []string{"application/json"}},
			}
			response := httptest.NewRecorder()
			handler.CreateLocality()(response, request)

			require.Equal(t, c.ExpectedStatusCode, response.Result().StatusCode)
			require.Equal(t, "application/json", response.Header().Get("Content-Type"))
		})
	}
}

func TestUnitLocality_GetSellersByLocalityId(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		DataToReturn       []internal.SellersByLocality
		RawQuery           string
		ExpectedStatusCode int
	}{
		{
			TestName:           "OK",
			ExpectedStatusCode: 200,
			DataToReturn:       []internal.SellersByLocality{},
			ErrorToReturn:      nil,
			RawQuery:           "",
		},
		{
			TestName:           "BAD_REQUEST",
			ExpectedStatusCode: 400,
			DataToReturn:       []internal.SellersByLocality{},
			ErrorToReturn:      nil,
			RawQuery:           "id=asd",
		},
		{
			TestName:           "NOT_FOUND",
			ExpectedStatusCode: 404,
			DataToReturn:       []internal.SellersByLocality{},
			ErrorToReturn:      utils.ErrNotFound,
			RawQuery:           "id=99",
		},
		{
			TestName:           "INTERNAL_SERVER_ERROR",
			ExpectedStatusCode: 500,
			DataToReturn:       []internal.SellersByLocality{},
			ErrorToReturn:      errors.New("Internal Server Error"),
			RawQuery:           "id=1",
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(MockLocalityService)
			service.On("GetSellersByLocalityID", mock.Anything).Return(c.DataToReturn, c.ErrorToReturn)
			handler := handler.NewLocalityHandler(service)
			request := &http.Request{
				URL:    &url.URL{RawQuery: c.RawQuery},
				Header: http.Header{"Content-Type": []string{"application/json"}},
			}
			response := httptest.NewRecorder()
			handler.GetSellersByLocalityID()(response, request)
			require.Equal(t, c.ExpectedStatusCode, response.Result().StatusCode)
		})
	}

}
