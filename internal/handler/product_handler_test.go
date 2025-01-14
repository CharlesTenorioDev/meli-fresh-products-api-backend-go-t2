package handler_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/handler"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/repository"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/routes"
// 	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
// 	"github.com/stretchr/testify/require"
// )

// var mockProduct = internal.Product{
// 	ID: 1,
// 	ProductAttributes: internal.ProductAttributes{
// 		ProductCode:                    "123",
// 		Description:                    "test",
// 		Width:                          1,
// 		Height:                         1,
// 		Length:                         1,
// 		NetWeight:                      1,
// 		ExpirationRate:                 1,
// 		RecommendedFreezingTemperature: 1,
// 		FreezingRate:                   1,
// 		ProductType:                    1,
// 		SellerID:                       1,
// 	},
// }
// var mockProductTypeService = service.NewProductTypeService(repository.NewProductTypeDB(map[int]internal.ProductType{
// 	1: {ID: 1, Description: "test"},
// }))

// var mockService = service.NewProductService(repository.NewProductDB(map[int]internal.Product{
// 	1: mockProduct,
// }), mockProductTypeService)

// func TestProductHandler_GetProducts(t *testing.T) {
// 	h := handler.NewProductHandler(mockService)
// 	req := httptest.NewRequest(http.MethodGet, "/products", nil)
// 	w := httptest.NewRecorder()
// 	h.GetProducts(w, req)
// 	require.Equal(t, http.StatusOK, w.Code)
// }

// func TestProductHandler_GetProductByID_WhenExists(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusOK, w.Code)
// }

// func TestProductHandler_GetProductByID_WhenNotExists(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodGet, "/products/99", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusNotFound, w.Code)
// }

// func TestProductHandler_CreateProduct(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	mockProduct.ProductCode = "1234"
// 	body, err := json.Marshal(mockProduct)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", bytes.NewReader(body))
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusCreated, w.Code)
// }

// func TestProductHandler_UpdateProduct(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	body, err := json.Marshal(mockProduct)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodPatch, "/api/v1/products/1", bytes.NewReader(body))
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusOK, w.Code)
// }

// func TestProductHandler_DeleteProduct(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusNoContent, w.Code)
// }

// func TestProductHandler_DeleteProduct_WhenNotExists(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodDelete, "/products/99", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusNotFound, w.Code)
// }

// func TestProductHandler_CreateProduct_WhenEmptyFields(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", nil)
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusBadRequest, w.Code)
// }

// func TestProductHandler_CreateProduct_WhenDuplicated(t *testing.T) {
// 	router := chi.NewRouter()
// 	err := routes.NewProductRoutes(router, mockService)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	body, err := json.Marshal(mockProduct)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req := httptest.NewRequest(http.MethodPost, "/api/v1/products/", bytes.NewReader(body))
// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, req)
// 	require.Equal(t, http.StatusConflict, w.Code)
// }
// func TestProductHandler_CreateProduct_WhenEmptyFieldsAndDuplicated(t *testing.T) {
// 	h := handler.NewProductHandler(mockService)
// 	req := httptest.NewRequest(http.MethodPost, "/products", nil)
// 	w := httptest.NewRecorder()
// 	h.CreateProduct(w, req)
// 	require.Equal(t, http.StatusBadRequest, w.Code)
// }
