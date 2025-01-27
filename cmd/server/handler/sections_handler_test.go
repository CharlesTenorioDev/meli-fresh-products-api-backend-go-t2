package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/cmd/server/handler"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSectionService struct {
	mock.Mock
}

func (m *MockSectionService) GetAll() ([]internal.Section, error) {
	args := m.Called()
	return args.Get(0).([]internal.Section), args.Error(1)
}

func (m *MockSectionService) Save(section internal.Section) (internal.Section, error) {
	args := m.Called(section)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionService) Update(id int, toUpdate internal.SectionPointers) (internal.Section, error) {
	args := m.Called(id, toUpdate)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionService) GetByID(id int) (internal.Section, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Section), args.Error(1)
}

func (m *MockSectionService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSectionService) GetSectionProductsReport(id int) ([]internal.SectionProductsReport, error) {
	args := m.Called(id)
	return args.Get(0).([]internal.SectionProductsReport), args.Error(1)
}

var mockSection = internal.Section{
	ID:                 1,
	SectionNumber:      1,
	CurrentTemperature: 1,
	MinimumTemperature: 1,
	CurrentCapacity:    1,
	MinimumCapacity:    1,
	MaximumCapacity:    1,
	WarehouseID:        1,
	ProductTypeID:      1,
}
var mockSectionProductsReport = []internal.SectionProductsReport{
	{
		SectionID:     1,
		SectionNumber: 1,
		ProductsCount: 20,
	},
}

func TestUnitHandler_GetAll(t *testing.T) {

	t.Run("READ-FIND_ALL-200", func(t *testing.T) {
		mockService := new(MockSectionService)
		mockService.On("GetAll").Return([]internal.Section{mockSection}, nil)
		sectionHandler := handler.NewSectionHandler(mockService)

		req := &http.Request{Method: "GET"}
		res := httptest.NewRecorder()
		sectionHandler.GetAll()(res, req)

		require.Equal(t, res.Result().StatusCode, 200)
	})

	t.Run("READ-FIND_ALL-500", func(t *testing.T) {
		mockService := new(MockSectionService)
		mockService.On("GetAll").Return([]internal.Section{}, errors.New("internal error"))
		sectionHandler := handler.NewSectionHandler(mockService)

		req := &http.Request{Method: "GET"}
		res := httptest.NewRecorder()
		sectionHandler.GetAll()(res, req)

		require.Equal(t, res.Result().StatusCode, 500)
	})
}

func TestUnitSection_GetById(t *testing.T) {
	scenarios := []struct {
		Name               string
		ID                 string
		ExpectedBody       string
		ExpectedStatusCode int
		MockData           internal.Section
		MockError          error
	}{
		{
			Name:               "READ-FIND_BY_ID-200",
			ID:                 "1",
			ExpectedBody:       `{"data":{"id":1,"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}}`,
			ExpectedStatusCode: 200,
			MockData:           mockSection,
			MockError:          nil,
		},
		{
			Name:               "READ-FIND_BY_ID-400",
			ID:                 "ads",
			ExpectedBody:       `{"status":"Bad Request","message":"invalid format: id with invalid format"}`,
			ExpectedStatusCode: 400,
			MockData:           mockSection,
			MockError:          utils.EBadRequest("id"),
		},
		{
			Name:               "READ-FIND_BY_ID-404",
			ID:                 "1",
			ExpectedBody:       `{"message":"entity not found: section doesn't exist", "status":"Not Found"}`,
			ExpectedStatusCode: 404,
			MockData:           mockSection,
			MockError:          utils.ENotFound("section"),
		},
		{
			Name:               "READ-FIND_BY_ID-500",
			ID:                 "1",
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: 500,
			MockData:           mockSection,
			MockError:          errors.New("internal error"),
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(s *testing.T) {
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", scenario.ID)

			mockService := new(MockSectionService)

			sectionHandler := handler.NewSectionHandler(mockService)

			mockService.On("GetByID", mock.Anything).Return(scenario.MockData, scenario.MockError)

			req := &http.Request{Method: "GET"}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
			res := httptest.NewRecorder()

			sectionHandler.GetById()(res, req)

			require.JSONEq(s, scenario.ExpectedBody, strings.TrimSpace(res.Body.String()))
			require.Equal(s, scenario.ExpectedStatusCode, res.Result().StatusCode)
			require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
		})
	}
}

func TestUnitSection_Delete(t *testing.T) {
	cases := []struct {
		Name               string
		ID                 string
		ExpectedBody       string
		ExpectedStatusCode int
		MockError          error
	}{
		{
			Name:               "DELETE-DELETE_BY_ID-204",
			ID:                 "1",
			ExpectedBody:       ``,
			ExpectedStatusCode: 204,
			MockError:          nil,
		},
		{
			Name:               "DELETE-DELETE_BY_ID-404",
			ID:                 "9",
			ExpectedBody:       `{"message":"entity not found: section doesn't exist", "status":"Not Found"}`,
			ExpectedStatusCode: 404,
			MockError:          utils.ENotFound("section"),
		},
		{
			Name:               "DELETE-DELETE_BY_ID-400",
			ID:                 "asd",
			ExpectedBody:       `{"message":"invalid format: id with invalid format", "status":"Bad Request"}`,
			ExpectedStatusCode: 400,
			MockError:          nil,
		},
		{
			Name:               "DELETE-DELETE_BY_ID-500",
			ID:                 "9",
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: 500,
			MockError:          errors.New("Internal error occurs"),
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(s *testing.T) {
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", c.ID)

			mockService := new(MockSectionService)

			sectionHandler := handler.NewSectionHandler(mockService)

			mockService.On("Delete", mock.Anything).Return(c.MockError)

			req := &http.Request{Method: "DELETE"}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
			res := httptest.NewRecorder()

			sectionHandler.Delete()(res, req)

			if c.ExpectedBody != "" {
				require.JSONEq(s, c.ExpectedBody, strings.TrimSpace(res.Body.String()))
				require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
			}
			require.Equal(s, c.ExpectedStatusCode, res.Result().StatusCode)
		})
	}
}

func TestUnitSection_Post(t *testing.T) {
	cases := []struct {
		Name               string
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
		MockArgs           int
		MockData           internal.Section
		MockError          error
	}{
		{
			Name:               "POST-CREATE-201",
			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"data":{"id":1,"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}}`,
			ExpectedStatusCode: 201,
			MockData:           mockSection,
			MockError:          nil,
		},
		{
			Name:               "POST-CREATE-400",
			Body:               `{section_number:1current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"message":"invalid format: body with invalid format", "status":"Bad Request"}`,
			ExpectedStatusCode: 400,
			MockData:           internal.Section{},
			MockError:          nil,
		},
		{
			Name:               "POST-CREATE-409",
			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"message":"entity already exists: section with attribute 'id: 1' already exists", "status":"Conflict"}`,
			ExpectedStatusCode: 409,
			MockData:           internal.Section{},
			MockError:          utils.EConflict("section", "id: 1"),
		},
		{
			Name:               "POST-CREATE-422",
			Body:               `{"section_number":0,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments"}`,
			ExpectedStatusCode: 422,
			MockData:           internal.Section{},
			MockError:          utils.ErrInvalidArguments,
		},
		{
			Name:               "POST-CREATE-500",
			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: 500,
			MockData:           internal.Section{},
			MockError:          errors.New("Internal error occurs"),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(s *testing.T) {
			mockService := new(MockSectionService)

			sectionHandler := handler.NewSectionHandler(mockService)

			mockService.On("Save", mock.Anything).Return(mockSection, c.MockError)

			req := httptest.NewRequest("POST", "/", strings.NewReader(c.Body))
			res := httptest.NewRecorder()

			sectionHandler.CreateSection()(res, req)

			require.JSONEq(s, c.ExpectedBody, strings.TrimSpace(res.Body.String()))
			require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
			require.Equal(s, c.ExpectedStatusCode, res.Result().StatusCode)
		})
	}
}

func TestUnitSection_Patch(t *testing.T) {
	cases := []struct {
		Name               string
		ID                 string
		Body               string
		ExpectedBody       string
		ExpectedStatusCode int
		MockArgs           int
		MockData           internal.Section
		MockError          error
	}{
		{
			Name:               "PATCH-UPDATE-200",
			ID:                 "1",
			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"data":{"id":1,"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}}`,
			ExpectedStatusCode: 200,
			MockData:           mockSection,
			MockError:          nil,
		},
		{
			Name:               "PATCH-UPDATE-400_BODY",
			ID:                 "1",
			Body:               `{section_number:1current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"message":"invalid format: body with invalid format", "status":"Bad Request"}`,
			ExpectedStatusCode: 400,
			MockData:           internal.Section{},
			MockError:          nil,
		},
		{
			Name:               "PATCH-UPDATE-400_ID",
			ID:                 "asd",
			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"message":"invalid format: id with invalid format", "status":"Bad Request"}`,
			ExpectedStatusCode: 400,
			MockData:           internal.Section{},
			MockError:          nil,
		},
		{
			Name:               "PATCH-UPDATE-409",
			ID:                 "1",
			Body:               `{"section_number":2,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"status":"Conflict","message":"entity already exists"}`,
			ExpectedStatusCode: 409,
			MockData:           internal.Section{},
			MockError:          utils.ErrConflict,
		},
		{
			Name:               "PATCH-UPDATE-422",
			ID:                 "1",
			Body:               `{"section_number":0,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments"}`,
			ExpectedStatusCode: 422,
			MockData:           internal.Section{},
			MockError:          utils.ErrInvalidArguments,
		},
		{
			Name:               "PATCH-UPDATE-500",
			ID:                 "1",
			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: 500,
			MockData:           internal.Section{},
			MockError:          errors.New("Internal error occurs"),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(s *testing.T) {
			routeContext := chi.NewRouteContext()
			routeContext.URLParams.Add("id", c.ID)

			mockService := new(MockSectionService)

			sectionHandler := handler.NewSectionHandler(mockService)

			mockService.On("Update", mock.Anything, mock.Anything).Return(mockSection, c.MockError)

			req := httptest.NewRequest("PATCH", "/", strings.NewReader(c.Body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
			res := httptest.NewRecorder()

			sectionHandler.Update()(res, req)

			require.JSONEq(s, c.ExpectedBody, strings.TrimSpace(res.Body.String()))
			require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
			require.Equal(s, c.ExpectedStatusCode, res.Result().StatusCode)
		})
	}
}

func TestUnitSection_GetSectionProductsReport(t *testing.T) {
	cases := []struct {
		Name               string
		RawQuery           string
		ExpectedBody       string
		ExpectedStatusCode int
		MockError          error
		MockData           []internal.SectionProductsReport
	}{
		{
			Name:               "GET-GET_SECTION_BY_PRODUCTS-200",
			RawQuery:           "",
			ExpectedBody:       `{"data":[{"products_count":20, "section_id":1, "section_number":1}]}`,
			ExpectedStatusCode: 200,
			MockError:          nil,
			MockData:           mockSectionProductsReport,
		},
		{
			Name:               "GET-GET_SECTION_BY_PRODUCTS-400",
			RawQuery:           "id=asd",
			ExpectedBody:       `{"message":"invalid format: id with invalid format", "status":"Bad Request"}`,
			ExpectedStatusCode: 400,
			MockError:          utils.EBadRequest("id"),
			MockData:           nil,
		},
		{
			Name:               "GET-GET_SECTION_BY_PRODUCTS-400",
			RawQuery:           "id=9999",
			ExpectedBody:       `{"message":"entity not found: section doesn't exist", "status":"Not Found"}`,
			ExpectedStatusCode: 404,
			MockError:          utils.ENotFound("section"),
			MockData:           nil,
		},
		{
			Name:               "GET-GET_SECTION_BY_PRODUCTS-500",
			RawQuery:           "id=1",
			ExpectedBody:       `{"message":"internal server error", "status":"Internal Server Error"}`,
			ExpectedStatusCode: 500,
			MockError:          errors.New("Internal error occurs"),
			MockData:           nil,
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(s *testing.T) {
			mockService := new(MockSectionService)

			sectionHandler := handler.NewSectionHandler(mockService)

			mockService.On("GetSectionProductsReport", mock.Anything).Return(c.MockData, c.MockError)

			req := httptest.NewRequest("GET", "/?"+c.RawQuery, nil)
			res := httptest.NewRecorder()

			sectionHandler.GetSectionProductsReport()(res, req)

			if c.ExpectedBody != "" {
				require.JSONEq(s, c.ExpectedBody, strings.TrimSpace(res.Body.String()))
				require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
			}
			require.Equal(s, c.ExpectedStatusCode, res.Result().StatusCode)
		})
	}
}
