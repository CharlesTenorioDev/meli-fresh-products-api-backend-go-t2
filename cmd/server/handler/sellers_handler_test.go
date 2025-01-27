package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-go-t2/internal"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockSellerService struct {
	mock.Mock
}

func (s *MockSellerService) GetAll() ([]internal.Seller, error) {
	args := s.Called()
	return args.Get(0).([]internal.Seller), args.Error(1)
}

func (s *MockSellerService) GetByID(id int) (internal.Seller, error) {
	args := s.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (s *MockSellerService) Create(newSeller *internal.Seller) error {
	args := s.Called(newSeller)
	return args.Error(0)
}

func (s *MockSellerService) Update(id int, newSeller *internal.Seller) (internal.Seller, error) {
	args := s.Called(id, newSeller)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (s *MockSellerService) Delete(id int) error {
	args := s.Called(id)
	return args.Error(0)
}

func TestUnitSeller_GetAll_Success(t *testing.T) {
	sellers := []internal.Seller{
		{ID: 1, Cid: 55, CompanyName: "Company", Address: "Address", Telephone: "1199999999", LocalityID: 1},
		{ID: 2, Cid: 45, CompanyName: "Company", Address: "Address", Telephone: "1199999999", LocalityID: 2}}

	service := new(MockSellerService)
	service.On("GetAll").Return(sellers, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellers", nil)

	handler := NewSellerHandler(service)
	handler.GetAll()(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestUnitSeller_GetAll_InternalServerError(t *testing.T) {

	service := new(MockSellerService)
	service.On("GetAll").Return([]internal.Seller{}, errors.New("some error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellers", nil)

	handler := NewSellerHandler(service)
	handler.GetAll()(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUnitSeller_GetByID_Success(t *testing.T) {
	seller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1}

	service := new(MockSellerService)
	service.On("GetByID", mock.Anything).Return(seller, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.GetById()(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestUnitSeller_GetByID_BadRequest(t *testing.T) {

	service := new(MockSellerService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "a")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.GetById()(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnitSeller_GetByID_NotFound(t *testing.T) {

	service := new(MockSellerService)
	service.On("GetByID", mock.Anything).Return(internal.Seller{}, utils.ENotFound("Seller"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.GetById()(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestUnitSeller_GetByID_InternalServerError(t *testing.T) {

	service := new(MockSellerService)
	service.On("GetByID", mock.Anything).Return(internal.Seller{}, errors.New("some error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.GetById()(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUnitSeller_Create_Success(t *testing.T) {

	newSeller := internal.SellerRequest{
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	service := new(MockSellerService)
	service.On("Create", mock.Anything).Return(nil)

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sellers", bodyReader)
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Create()(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
}

func TestUnitSeller_Create_BadRequest(t *testing.T) {

	service := new(MockSellerService)

	bodyByte, _ := json.Marshal("a")
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sellers", bodyReader)
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Create()(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnitSeller_Create_Conflict(t *testing.T) {

	newSeller := internal.SellerRequest{
		CompanyName: "Company",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	service := new(MockSellerService)
	service.On("Create", mock.Anything).Return(utils.EConflict("Cid", "Seller"))

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sellers", bodyReader)
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Create()(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}

func TestUnitSeller_Create_InvalidOrEmptyArguments(t *testing.T) {

	newSeller := internal.SellerRequest{
		CompanyName: "",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	service := new(MockSellerService)
	service.On("Create", mock.Anything).Return(utils.EZeroValue("Company name"))

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sellers", bodyReader)
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Create()(w, req)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestUnitSeller_Create_InternalServerError(t *testing.T) {

	newSeller := internal.SellerRequest{
		CompanyName: "",
		Address:     "Address",
		Telephone:   "1199999999",
		LocalityID:  1,
	}

	service := new(MockSellerService)
	service.On("Create", mock.Anything).Return(errors.New("some error"))

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/sellers", bodyReader)
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Create()(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUnitSeller_Update_Success(t *testing.T) {

	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company changed",
		Address:     "Address changed",
		Telephone:   "1199999222",
		LocalityID:  1}

	service := new(MockSellerService)
	service.On("Update", mock.Anything, mock.Anything).Return(newSeller, nil)

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/sellers", bodyReader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Update()(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestUnitSeller_Update_BadRequestBody(t *testing.T) {

	service := new(MockSellerService)

	bodyByte, _ := json.Marshal("a")
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/sellers", bodyReader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Update()(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnitSeller_Update_BadRequestInvalidID(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company changed",
		Address:     "Address changed",
		Telephone:   "1199999222",
		LocalityID:  1}

	service := new(MockSellerService)

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/sellers", bodyReader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "a")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Update()(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnitSeller_Update_Conflict(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company changed",
		Address:     "Address changed",
		Telephone:   "1199999222",
		LocalityID:  1}

	service := new(MockSellerService)
	service.On("Update", mock.Anything, mock.Anything).Return(internal.Seller{}, utils.EConflict("Cid", "Seller"))

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/sellers", bodyReader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Update()(w, req)

	require.Equal(t, http.StatusConflict, w.Code)
}

func TestUnitSeller_Update_NotFound(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company changed",
		Address:     "Address changed",
		Telephone:   "1199999222",
		LocalityID:  1}

	service := new(MockSellerService)
	service.On("Update", mock.Anything, mock.Anything).Return(internal.Seller{}, utils.ENotFound("Seller"))

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/sellers", bodyReader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Update()(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestUnitSeller_Update_InternalServerError(t *testing.T) {
	newSeller := internal.Seller{
		ID:          1,
		Cid:         55,
		CompanyName: "Company changed",
		Address:     "Address changed",
		Telephone:   "1199999222",
		LocalityID:  1}

	service := new(MockSellerService)
	service.On("Update", mock.Anything, mock.Anything).Return(internal.Seller{}, errors.New("some error"))

	bodyByte, _ := json.Marshal(newSeller)
	bodyReader := bytes.NewReader(bodyByte)
	req, _ := http.NewRequest(http.MethodPatch, "/api/v1/sellers", bodyReader)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler := NewSellerHandler(service)
	handler.Update()(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUnitSeller_Delete_Success(t *testing.T) {

	service := new(MockSellerService)
	service.On("Delete", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.Delete()(w, req)

	require.Equal(t, http.StatusNoContent, w.Code)
}

func TestUnitSeller_Delete_BadRequest(t *testing.T) {

	service := new(MockSellerService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "a")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.Delete()(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUnitSeller_Delete_NotFound(t *testing.T) {

	service := new(MockSellerService)
	service.On("Delete", mock.Anything).Return(utils.ENotFound("Seller"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.Delete()(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestUnitSeller_Delete_InternalServerError(t *testing.T) {

	service := new(MockSellerService)
	service.On("Delete", mock.Anything).Return(errors.New("some error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/sellers/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler := NewSellerHandler(service)
	handler.Delete()(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)
}
