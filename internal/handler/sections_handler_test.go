package handler_test

// import (
// 	"errors"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal"
// 	"io"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// type MockSectionService struct {
// 	mock.Mock
// }

// func (m *MockSectionService) GetAll() ([]internal.Section, error) {
// 	args := m.Called()
// 	return args.Get(0).([]internal.Section), args.Error(1)
// }

// func (m *MockSectionService) Save(section internal.Section) (internal.Section, error) {
// 	args := m.Called(section)
// 	return args.Get(0).(internal.Section), args.Error(1)
// }

// func (m *MockSectionService) Update(id int, toUpdate internal.SectionPointers) (internal.Section, error) {
// 	args := m.Called(id, toUpdate)
// 	return args.Get(0).(internal.Section), args.Error(1)
// }

// func (m *MockSectionService) GetById(id int) (internal.Section, error) {
// 	args := m.Called(id)
// 	return args.Get(0).(internal.Section), args.Error(1)
// }

// func (m *MockSectionService) Delete(id int) error {
// 	args := m.Called(id)
// 	return args.Error(0)
// }

// var simpleSection = internal.Section{
// 	ID:                 1,
// 	SectionNumber:      1,
// 	CurrentTemperature: 1,
// 	MinimumTemperature: 1,
// 	CurrentCapacity:    1,
// 	MinimumCapacity:    1,
// 	MaximumCapacity:    1,
// 	WarehouseID:        1,
// 	ProductTypeID:      1,
// }

// var simpleSectionPointers = internal.SectionPointers{
// 	SectionNumber:      &simpleSection.SectionNumber,
// 	CurrentTemperature: &simpleSection.CurrentTemperature,
// 	MinimumTemperature: &simpleSection.MinimumTemperature,
// 	CurrentCapacity:    &simpleSection.CurrentCapacity,
// 	MinimumCapacity:    &simpleSection.MinimumCapacity,
// 	MaximumCapacity:    &simpleSection.MaximumCapacity,
// 	WarehouseID:        &simpleSection.WarehouseID,
// 	ProductTypeID:      &simpleSection.ProductTypeID,
// }

// func Test_GetAll(t *testing.T) {
// 	t.Run("when exist", func(t *testing.T) {
// 		service := new(MockSectionService)
// 		handler := handler.NewSectionHandler(service)

// 		service.On("GetAll").Return([]internal.Section{simpleSection}, nil)

// 		req := httptest.NewRequest("GET", "/api/v1/sections", nil)
// 		res := httptest.NewRecorder()
// 		handler.GetAll()(res, req)

// 		require.Equal(t, res.Result().StatusCode, 200)
// 	})

// 	t.Run("when interal error", func(t *testing.T) {
// 		service := new(MockSectionService)
// 		handler := handler.NewSectionHandler(service)

// 		service.On("GetAll").Return([]internal.Section{}, errors.New("internal error"))

// 		req := httptest.NewRequest("GET", "/api/v1/sections", nil)
// 		res := httptest.NewRecorder()
// 		handler.GetAll()(res, req)

// 		require.Equal(t, res.Result().StatusCode, 500)
// 	})
// }

// func Test_GetById(t *testing.T) {
// 	scenarios := []struct {
// 		Name               string
// 		Path               string
// 		ExpectedBody       string
// 		ExpectedStatusCode int
// 		MockArgs           int
// 		MockReturn1        internal.Section
// 		MockReturn2        error
// 	}{
// 		{
// 			Name:               "when id exists",
// 			Path:               "/api/v1/sections/1",
// 			ExpectedBody:       `{"data":{"id":1,"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}}`,
// 			ExpectedStatusCode: 200,
// 			MockArgs:           1,
// 			MockReturn1:        simpleSection,
// 			MockReturn2:        nil,
// 		},
// 		{
// 			Name:               "when id does not exist",
// 			Path:               "/api/v1/sections/9",
// 			ExpectedBody:       `{"status":"Not Found","message":"no section for id 9"}`,
// 			ExpectedStatusCode: 404,
// 			MockArgs:           9,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        utils.ErrNotFound,
// 		},
// 		{
// 			Name:               "when id is invalid",
// 			Path:               "/api/v1/sections/r1kslw",
// 			ExpectedBody:       `{"status":"Bad Request","message":"invalid id"}`,
// 			ExpectedStatusCode: 400,
// 			MockArgs:           0,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        nil,
// 		},
// 		{
// 			Name:               "when internal server error",
// 			Path:               "/api/v1/sections/9",
// 			ExpectedBody:       `{"status":"Internal Server Error","message":"Some error occurs"}`,
// 			ExpectedStatusCode: 500,
// 			MockArgs:           9,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        errors.New("Internal error occurs"),
// 		},
// 	}
// 	for _, scenario := range scenarios {
// 		t.Run(scenario.Name, func(s *testing.T) {
// 			service := new(MockSectionService)
// 			mux := chi.NewRouter()
// 			routes.RegisterSectionRoutes(mux, service)

// 			service.On("GetById", scenario.MockArgs).Return(scenario.MockReturn1, scenario.MockReturn2)

// 			req := httptest.NewRequest("GET", scenario.Path, nil)
// 			res := httptest.NewRecorder()

// 			mux.ServeHTTP(res, req)

// 			bContent, _ := io.ReadAll(res.Body)
// 			responseBody := strings.TrimSpace(string(bContent))

// 			require.Equal(s, scenario.ExpectedStatusCode, res.Result().StatusCode)
// 			require.Equal(s, scenario.ExpectedBody, responseBody)
// 			require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
// 		})
// 	}
// }

// func Test_Post(t *testing.T) {
// 	scenarios := []struct {
// 		Name               string
// 		Path               string
// 		Body               string
// 		ExpectedBody       string
// 		ExpectedStatusCode int
// 		MockArgs           int
// 		MockReturn1        internal.Section
// 		MockReturn2        error
// 	}{
// 		{
// 			Name:               "when no error occurs",
// 			Path:               "/api/v1/sections",
// 			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"data":{"id":1,"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}}`,
// 			ExpectedStatusCode: 201,
// 			MockReturn1:        simpleSection,
// 			MockReturn2:        nil,
// 		},
// 		{
// 			Name:               "when invalid body",
// 			Path:               "/api/v1/sections",
// 			Body:               `{section_number:1current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Bad Request","message":"invalid format"}`,
// 			ExpectedStatusCode: 400,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        nil,
// 		},
// 		{
// 			Name:               "when section already exists for section_number",
// 			Path:               "/api/v1/sections",
// 			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Conflict","message":"entity already exists"}`,
// 			ExpectedStatusCode: 409,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        utils.ErrConflict,
// 		},
// 		{
// 			Name:               "when fields are invalid",
// 			Path:               "/api/v1/sections",
// 			Body:               `{"section_number":0,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments"}`,
// 			ExpectedStatusCode: 422,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        utils.ErrInvalidArguments,
// 		},
// 		{
// 			Name:               "when internal server error",
// 			Path:               "/api/v1/sections",
// 			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Internal Server Error","message":"Some error occurs"}`,
// 			ExpectedStatusCode: 500,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        errors.New("Internal error occurs"),
// 		},
// 	}

// 	for _, scenario := range scenarios {
// 		t.Run(scenario.Name, func(s *testing.T) {
// 			service := new(MockSectionService)
// 			mux := chi.NewRouter()
// 			routes.RegisterSectionRoutes(mux, service)

// 			service.On("Save", mock.Anything).Return(scenario.MockReturn1, scenario.MockReturn2)

// 			req := httptest.NewRequest("POST", scenario.Path, strings.NewReader(scenario.Body))
// 			res := httptest.NewRecorder()

// 			mux.ServeHTTP(res, req)

// 			bContent, _ := io.ReadAll(res.Body)
// 			responseBody := strings.TrimSpace(string(bContent))

// 			require.Equal(s, scenario.ExpectedStatusCode, res.Result().StatusCode)
// 			require.Equal(s, scenario.ExpectedBody, responseBody)
// 			require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
// 		})
// 	}
// }

// func Test_Patch(t *testing.T) {
// 	scenarios := []struct {
// 		Name               string
// 		Path               string
// 		Body               string
// 		ExpectedBody       string
// 		ExpectedStatusCode int
// 		MockArgs           int
// 		MockReturn1        internal.Section
// 		MockReturn2        error
// 	}{
// 		{
// 			Name:               "when no error occurs",
// 			Path:               "/api/v1/sections/1",
// 			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"data":{"id":1,"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}}`,
// 			ExpectedStatusCode: 200,
// 			MockReturn1:        simpleSection,
// 			MockReturn2:        nil,
// 		},
// 		{
// 			Name:               "when invalid body",
// 			Path:               "/api/v1/sections/1",
// 			Body:               `{section_number:1current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Bad Request","message":"invalid format"}`,
// 			ExpectedStatusCode: 400,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        nil,
// 		},
// 		{
// 			Name:               "when section already exists for section_number",
// 			Path:               "/api/v1/sections/1",
// 			Body:               `{"section_number":2,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Conflict","message":"entity already exists"}`,
// 			ExpectedStatusCode: 409,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        utils.ErrConflict,
// 		},
// 		{
// 			Name:               "when fields are invalid",
// 			Path:               "/api/v1/sections/1",
// 			Body:               `{"section_number":0,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Unprocessable Entity","message":"invalid arguments"}`,
// 			ExpectedStatusCode: 422,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        utils.ErrInvalidArguments,
// 		},
// 		{
// 			Name:               "when internal server error",
// 			Path:               "/api/v1/sections/1",
// 			Body:               `{"section_number":1,"current_capacity":1,"maximum_capacity":1,"minimum_capacity":1,"current_temperature":1,"minimum_temperature":1,"warehouse_id":1,"product_type_id":1}`,
// 			ExpectedBody:       `{"status":"Internal Server Error","message":"Some error occurs"}`,
// 			ExpectedStatusCode: 500,
// 			MockReturn1:        internal.Section{},
// 			MockReturn2:        errors.New("Internal error occurs"),
// 		},
// 	}

// 	for _, scenario := range scenarios {
// 		t.Run(scenario.Name, func(s *testing.T) {
// 			service := new(MockSectionService)
// 			mux := chi.NewRouter()
// 			routes.RegisterSectionRoutes(mux, service)

// 			service.On("Update", 1, mock.Anything).Return(scenario.MockReturn1, scenario.MockReturn2)

// 			req := httptest.NewRequest("PATCH", scenario.Path, strings.NewReader(scenario.Body))
// 			res := httptest.NewRecorder()

// 			mux.ServeHTTP(res, req)

// 			bContent, _ := io.ReadAll(res.Body)
// 			responseBody := strings.TrimSpace(string(bContent))

// 			require.Equal(s, scenario.ExpectedStatusCode, res.Result().StatusCode)
// 			require.Equal(s, scenario.ExpectedBody, responseBody)
// 			require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
// 		})
// 	}
// }

// func Test_Delete(t *testing.T) {
// 	scenarios := []struct {
// 		Name               string
// 		Path               string
// 		ExpectedBody       string
// 		ExpectedStatusCode int
// 		MockArgs           int
// 		MockReturn         error
// 	}{
// 		{
// 			Name:               "when id exists",
// 			Path:               "/api/v1/sections/1",
// 			ExpectedBody:       ``,
// 			ExpectedStatusCode: 204,
// 			MockArgs:           1,
// 			MockReturn:         nil,
// 		},
// 		{
// 			Name:               "when id does not exist",
// 			Path:               "/api/v1/sections/9",
// 			ExpectedBody:       `{"status":"Not Found","message":"no section for id 9"}`,
// 			ExpectedStatusCode: 404,
// 			MockArgs:           9,
// 			MockReturn:         utils.ErrNotFound,
// 		},
// 		{
// 			Name:               "when id is invalid",
// 			Path:               "/api/v1/sections/r1kslw",
// 			ExpectedBody:       `{"status":"Bad Request","message":"invalid id"}`,
// 			ExpectedStatusCode: 400,
// 			MockArgs:           0,
// 			MockReturn:         nil,
// 		},
// 		{
// 			Name:               "when internal server error",
// 			Path:               "/api/v1/sections/9",
// 			ExpectedBody:       `{"status":"Internal Server Error","message":"Some error occurs"}`,
// 			ExpectedStatusCode: 500,
// 			MockArgs:           9,
// 			MockReturn:         errors.New("Internal error occurs"),
// 		},
// 	}
// 	for _, scenario := range scenarios {
// 		t.Run(scenario.Name, func(s *testing.T) {
// 			service := new(MockSectionService)
// 			mux := chi.NewRouter()
// 			routes.RegisterSectionRoutes(mux, service)

// 			service.On("Delete", scenario.MockArgs).Return(scenario.MockReturn)

// 			req := httptest.NewRequest("DELETE", scenario.Path, nil)
// 			res := httptest.NewRecorder()

// 			mux.ServeHTTP(res, req)

// 			bContent, _ := io.ReadAll(res.Body)
// 			responseBody := strings.TrimSpace(string(bContent))

// 			require.Equal(s, scenario.ExpectedStatusCode, res.Result().StatusCode)
// 			require.Equal(s, scenario.ExpectedBody, responseBody)
// 			if scenario.ExpectedStatusCode != 204 {
// 				require.Equal(s, "application/json", res.Result().Header.Get("Content-Type"))
// 			}
// 		})
// 	}
// }
