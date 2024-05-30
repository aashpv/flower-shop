package delete

import (
	"flower-shop/user-service/internal/config"
	"flower-shop/user-service/internal/http-server/handlers"
	"flower-shop/user-service/internal/lib/jwt"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
	Email string `json:"email"`
}

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func DeletePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user-service/web/templates/delete.html")
}

func New(log *slog.Logger, hrs handlers.Handlers, cfg config.JwtConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("failed to open file", slog.Any("err", err))
		}
		defer file.Close()

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request body", // update msg
			})

			return
		}

		log.Info("request body is decoded!")

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "invalid request", // update msg
			})

			return
		}

		log.Info("request body is valid!")

		err = hrs.DeleteUser(req.Email)
		if err != nil {
			log.Error("failed to delete user", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "failed to delete user",
			})

			return
		}

		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		uid := jwt.GetData(tokenString, cfg)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("DELETE: [%s] user: %s with email: %d deleted successfully", date, uid, req.Email)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")
		render.JSON(w, r, Response{
			Status: http.StatusOK,
		})
	}
}
