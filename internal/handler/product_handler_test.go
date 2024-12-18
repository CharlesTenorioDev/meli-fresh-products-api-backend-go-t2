package handler

import (
	"net/http"
	"testing"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/pkg"
)

func TestProductHandler_GetProducts(t *testing.T) {
	type fields struct {
		service pkg.ProductService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProductHandler{
				service: tt.fields.service,
			}
			p.GetProducts(tt.args.w, tt.args.r)
		})
	}
}
