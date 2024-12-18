package handler

import (
	"net/http"

	"github.com/meli-fresh-products-api-backend-go-t2/internal/service"
	"github.com/meli-fresh-products-api-backend-go-t2/internal/utils"
)

func GetBuyersHandler(w http.ResponseWriter, r *http.Request) {
	buyers := service.GetAllBuyers()
	utils.JSON(w, http.StatusOK, buyers)
}
