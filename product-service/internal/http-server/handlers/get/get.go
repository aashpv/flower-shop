package get

import (
	"flower-shop/product-service/internal/http-server/handlers"
	"flower-shop/product-service/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	Product model.Product `json:"product"`
	Status  int           `json:"status"`
	Error   string        `json:"error,omitempty"`
}

func ProductPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "product-service/web/templates/product.html")
}

func New(log *slog.Logger, hrs handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("failed to parse id", err)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		product, err := hrs.GetProduct(id)
		if err != nil {
			log.Error("failed to get product", err)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		log.Info("got success")

		render.JSON(w, r, Response{
			Product: *product,
			Status:  http.StatusOK,
		})
	}
}
