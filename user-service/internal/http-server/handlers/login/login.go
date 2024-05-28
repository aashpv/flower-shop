package login

import (
	"errors"
	"flower-shop/user-service/internal/config"
	"flower-shop/user-service/internal/http-server/handlers"
	"flower-shop/user-service/internal/lib/jwt"
	"flower-shop/user-service/internal/storage"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	Token  string `json:"token,omitempty"`
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user-service/web/templates/login.html")
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
				Error:  "failed to decode request body",
			})

			return
		}

		log.Info("request body is decoded!")

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "invalid request",
			})

			return
		}

		log.Info("request body is valid!")

		id, hashPass, role, err := hrs.CheckUser(req.Email)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				log.Error("user not found", slog.Any("err", err))

				render.JSON(w, r, Response{
					Status: http.StatusNotFound,
					Error:  "user not found",
				})

				return
			}

			log.Error("internal server error", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		log.Info("user is checked!")

		if err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(req.Password)); err != nil {
			log.Error("Invalid password", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusNotFound,
				Error:  "Invalid data",
			})

			return
		}

		log.Info("password has been verified!")

		token, err := jwt.NewToken(id, role, cfg)
		if err != nil {
			log.Error("internal server error", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		r.Header.Set("Authorization", "Bearer "+token)

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("AUTH: [%s] user:%s with id:%d logged in successfully", date, req.Email, id)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")

		log.Info("token: ", r.Header.Get("Authorization"))

		render.JSON(w, r, Response{
			Token:  token,
			Status: http.StatusOK,
		})
	}
}
