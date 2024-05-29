package create

import (
	"errors"
	"flower-shop/product-service/internal/config"
	"flower-shop/product-service/internal/http-server/handlers"
	"flower-shop/product-service/internal/lib/jwt"
	"flower-shop/product-service/internal/model"
	"flower-shop/product-service/internal/storage"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func CreatePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "product-service/web/templates/create.html")
}

func New(log *slog.Logger, hrs handlers.Handlers, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file:", err)
		}
		defer file.Close()

		var req model.Product
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body:", err)

			render.JSON(w, r, Response{Status: http.StatusBadRequest, Error: "failed to decode request body"})

			return
		}

		log.Info("decode success")

		if err := validator.New().Struct(&req); err != nil {
			log.Error("invalid request:", err)

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "invalid request",
			})

			return
		}

		log.Info("validate success")

		id, err := hrs.CreateProduct(&req)
		if err != nil {
			if errors.Is(err, storage.ErrProductNotFound) {
				log.Error("product not found")

				render.JSON(w, r, Response{
					Status: http.StatusNotFound,
					Error:  "product not found",
				})

				return
			}
			log.Error("failed to create product:", err)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "failed to create product",
			})

			return
		}

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetUID(tokenString, cfg)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("CREATE: [%s] user with id:%s create product:%d", date, uid, id)
		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

		log.Info("created success")

		render.JSON(w, r, Response{
			Id:     id,
			Status: http.StatusOK,
		})
	}
}
