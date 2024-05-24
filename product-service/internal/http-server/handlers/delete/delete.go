package delete

import (
	"flower-shop/product-service/internal/config"
	"flower-shop/product-service/internal/http-server/handlers"
	"flower-shop/product-service/internal/lib/jwt"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func New(log *slog.Logger, hrs handlers.Handlers, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("open file:", err)
		}
		defer file.Close()

		idParam := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idParam)
		if err != nil {
			log.Error("failed to parse id:", err)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		err = hrs.DeleteProduct(id)
		if err != nil {
			log.Error("failed to delete product:", err)

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetData(tokenString, cfg)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("DELETE: [%s] user with id:%s delete product:%d", date, uid, id)
		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

		log.Info("deleted success")

		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
