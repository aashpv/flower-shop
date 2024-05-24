package list

import (
	"flower-shop/product-service/internal/http-server/handlers"
	"flower-shop/product-service/internal/model"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	Products []model.Product `json:"products"`
	Status   int             `json:"status"`
	Error    string          `json:"error,omitempty"`
}

func ProductsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "product-service/web/templates/products.html")
}

func New(log *slog.Logger, hrs handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := hrs.GetAllProducts()
		if err != nil {
			log.Error("failed to get all products", err)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		log.Info("got products")

		render.JSON(w, r, Response{
			Products: products,
			Status:   http.StatusOK,
		})
	}
}
