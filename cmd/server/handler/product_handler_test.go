package handler_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockProductService struct {
	mock.Mock
}

func (m *mockProductService) GetProducts() ([]internal.Product, error) {
	args := m.Called()
	return args.Get(0).([]internal.Product), args.Error(1)
}

func (m *mockProductService) GetProductByID(id int) (internal.Product, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *mockProductService) CreateProduct(product internal.ProductAttributes) (internal.Product, error) {
	args := m.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *mockProductService) UpdateProduct(product internal.Product) (internal.Product, error) {
	args := m.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (m *mockProductService) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProductHandler_GetProducts(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "GetProducts_OK",
			ExpectedBody:       `{"data":[{"id":1,"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"freezing_rate":10,"product_type":1,"seller_id":1}]}`,
			ExpectedStatusCode: http.StatusOK,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "GetProducts_InternalServerError",
			ExpectedBody:       `{"status":"Not Found","message":"entity not found: Product doesn't exist"}`,
			ExpectedStatusCode: http.StatusNotFound,
			ErrorToReturn:      utils.ENotFound("Product"),
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductService)
			service.On("GetProducts").Return([]internal.Product{
				{
					ID: 1,
					ProductAttributes: internal.ProductAttributes{
						ProductCode:                    "123",
						Description:                    "product1",
						Width:                          1000,
						Height:                         10,
						Length:                         10,
						NetWeight:                      10,
						ExpirationRate:                 10,
						RecommendedFreezingTemperature: 10,
						FreezingRate:                   10,
						ProductType:                    1,
						SellerID:                       1,
					},
				},
			}, c.ErrorToReturn)

			handler := handler.NewProductHandler(service)

			request := &http.Request{
				Body:   io.NopCloser(strings.NewReader(c.ExpectedBody)),
				Header: http.Header{"Content-Type": []string{"application/json"}},
			}
			response := httptest.NewRecorder()
			handler.GetProducts(response, request)
			require.Equal(t, c.ExpectedStatusCode, response.Result().StatusCode)
			require.Equal(t, "application/json", response.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, response.Body.String())

		})
	}
}

func TestProductHandler_GetProductByID(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedStatusCode int
	}{
		{
			TestName: "GetProductByID_OK",
			Body: `{
						"data": {
							"id": 1,
							"product_code": "PA001",
							"description": "Fresh Apples",
							"width": 5,
							"height": 4.5,
							"length": 7.5,
							"net_weight": 1.2,
							"expiration_rate": 0.1,
							"recommended_freezing_temperature": -2.5,
							"freezing_rate": 0.2,
							"product_type": 2,
							"seller_id": 1
							}
					}`,
			ExpectedStatusCode: http.StatusOK,
			ErrorToReturn:      nil,
		},
		{
			TestName: "GetProductByID_ErrorNotFound",
			Body: `{
					"status": "Not Found",
					"message": "entity not found: Product doesn't exist"
					}`,
			ExpectedStatusCode: http.StatusNotFound,
			ErrorToReturn:      utils.ENotFound("Product"),
		},
		{
			TestName: "GetProductByID_ErrorBadRequest",
			Body: `{
				"status": "Bad Request",
				"message": "invalid format: Invalid ID with invalid format"
			}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductService)
			service.On("GetProductByID", 1).Return(internal.Product{
				ID: 1,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "PA001",
					Description:                    "Fresh Apples",
					Width:                          5,
					Height:                         4.5,
					Length:                         7.5,
					NetWeight:                      1.2,
					ExpirationRate:                 0.1,
					RecommendedFreezingTemperature: -2.5,
					FreezingRate:                   0.2,
					ProductType:                    2,
					SellerID:                       1,
				},
			}, c.ErrorToReturn)
			handler := handler.NewProductHandler(service)

			req, _ := http.NewRequest("GET", "http://localhost:8080/api/v1/products/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if c.TestName == "GetProductByID_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			handler.GetProductByID(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.Body, res.Body.String())

		})
	}
}

func TestProductHandler_CreateProduct(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName: "CreateProduct_OK",
			Body:     `{"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"freezing_rate":10,"product_type":1,"seller_id":1}`,
			ExpectedBody: `{
									"data": {
										"id": 1,
										"product_code": "123",
										"description": "product1",
										"width": 1000,
										"height": 10,
										"length": 10,
										"net_weight": 10,
										"expiration_rate": 10,
										"recommended_freezing_temperature": 10,
										"freezing_rate": 10,
										"product_type": 1,
										"seller_id": 1
									}
								}`,
			ExpectedStatusCode: http.StatusCreated,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "CreateProduct_ErrorUnprocessableEntity",
			Body:               `{"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"product_type":1,"seller_id":1}`,
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments: Product Code cannot be empty/null"}`,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			ErrorToReturn:      utils.EZeroValue("Product Code"),
		},
		{
			TestName:           "CreateProduct_ErrorConflict",
			Body:               `{"product_code":"","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"freezing_rate":10,"product_type":1,"seller_id":1}`,
			ExpectedBody:       `{"status":"Conflict","message":"entity already exists: Product with attribute 'Product Code' already exists", "status":"Conflict"}`,
			ExpectedStatusCode: http.StatusConflict,
			ErrorToReturn:      utils.EConflict("Product", "Product Code"),
		},
		{
			TestName:           "CreateProduct_InternalServerError",
			Body:               "",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedBody:       `{"status":"Internal Server Error","message":"internal server error"}`,
			ErrorToReturn:      utils.EBadRequest("Invalid Message Format"),
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductService)
			service.On("CreateProduct", mock.Anything).Return(internal.Product{
				ID: 1,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "product1",
					Width:                          1000,
					Height:                         10,
					Length:                         10,
					NetWeight:                      10,
					ExpirationRate:                 10,
					RecommendedFreezingTemperature: 10,
					FreezingRate:                   10,
					ProductType:                    1,
					SellerID:                       1,
				},
			}, c.ErrorToReturn)
			handler := handler.NewProductHandler(service)

			req, _ := http.NewRequest("POST", "http://localhost:8080/api/v1/products/", strings.NewReader(c.Body))
			res := httptest.NewRecorder()
			handler.CreateProduct(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())

		})
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
	}{
		{
			TestName:           "UpdateProduct_OK",
			Body:               `{"id":1,"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"freezing_rate":10,"product_type":1,"seller_id":1}`,
			ExpectedBody:       `{"data":{"id":1,"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"freezing_rate":10,"product_type":1,"seller_id":1}}`,
			ExpectedStatusCode: http.StatusOK,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "UpdateProduct_ErrorBadRequest",
			Body:               `{"id":1,"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"product_type":1,"seller_id":1}`,
			ExpectedBody:       `{"status":"Bad Request","message":"invalid format: Invalid ID with invalid format"}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
		},
		{
			TestName:           "UpdateProduct_ErrorUnprocessableEntity",
			Body:               `{"id":1,"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"product_type":1,"seller_id":1}`,
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments: Freezing Rate cannot be empty/null"}`,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
			ErrorToReturn:      utils.EZeroValue("Freezing Rate"),
		},
		{
			TestName:           "UpdateProduct_ErrorConflict",
			Body:               `{"id":1,"product_code":"123","description":"product1","width":1000,"height":10,"length":10,"net_weight":10,"expiration_rate":10,"recommended_freezing_temperature":10,"freezing_rate":10,"product_type":1,"seller_id":1}`,
			ExpectedBody:       `{"status":"Conflict","message":"entity already exists: Product with attribute 'Product Code' already exists"}`,
			ExpectedStatusCode: http.StatusConflict,
			ErrorToReturn:      utils.EConflict("Product", "Product Code"),
		},
		{
			TestName:           "UpdateProduct_InternalServerError",
			Body:               "",
			ExpectedBody:       `{"status":"Internal Server Error","message":"internal server error"}`,
			ExpectedStatusCode: http.StatusInternalServerError,
			ErrorToReturn:      utils.EBadRequest("Invalid Message Format"),
		},
	}
	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductService)
			service.On("UpdateProduct", mock.Anything).Return(internal.Product{
				ID: 1,
				ProductAttributes: internal.ProductAttributes{
					ProductCode:                    "123",
					Description:                    "product1",
					Width:                          1000,
					Height:                         10,
					Length:                         10,
					NetWeight:                      10,
					ExpirationRate:                 10,
					RecommendedFreezingTemperature: 10,
					FreezingRate:                   10,
					ProductType:                    1,
					SellerID:                       1,
				},
			}, c.ErrorToReturn)
			handler := handler.NewProductHandler(service)

			req, _ := http.NewRequest("PATCH", "http://localhost:8080/api/v1/products/", strings.NewReader(c.Body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if c.TestName == "UpdateProduct_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			handler.UpdateProduct(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
			require.JSONEq(t, c.ExpectedBody, res.Body.String())
		})
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	cases := []struct {
		TestName           string
		ErrorToReturn      error
		ExpectedStatusCode int
	}{
		{
			TestName:           "DeleteProduct_OK",
			ExpectedStatusCode: http.StatusNoContent,
			ErrorToReturn:      nil,
		},
		{
			TestName:           "DeleteProduct_ErrorBadRequest",
			ExpectedStatusCode: http.StatusBadRequest,
			ErrorToReturn:      utils.EBadRequest("Invalid ID"),
		},
		{
			TestName:           "DeleteProduct_ErrorNotFound",
			ExpectedStatusCode: http.StatusNotFound,
			ErrorToReturn:      utils.ENotFound("Product"),
		},
	}

	for _, c := range cases {
		t.Run(c.TestName, func(t *testing.T) {
			service := new(mockProductService)
			service.On("DeleteProduct", 1).Return(c.ErrorToReturn)
			handler := handler.NewProductHandler(service)

			req, _ := http.NewRequest("DELETE", "http://localhost:8080/api/v1/products/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", "1")
			if c.TestName == "DeleteProduct_ErrorBadRequest" {
				rctx.URLParams.Add("id", "a")
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			res := httptest.NewRecorder()
			handler.DeleteProduct(res, req)
			require.Equal(t, c.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(t, "application/json", res.Header().Get("Content-Type"))
		})
	}
}
