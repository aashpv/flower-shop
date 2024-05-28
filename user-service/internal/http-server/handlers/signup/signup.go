package signup

import (
	"errors"
	"flower-shop/user-service/internal/http-server/handlers"
	"flower-shop/user-service/internal/model"
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
	model.User
}

type Response struct {
	Id     int    `json:"id"`
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func SignUpPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "user-service/web/templates/signup.html")
}

func New(log *slog.Logger, hrs handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.OpenFile("data.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Error("failed to open file", slog.Any("err", err))
		}
		defer file.Close()

		var req model.User

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "failed to decode request body", // update msg
			})

			return
		}

		log.Info("request body is decoded!")

		req.Role = "user"

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusBadRequest,
				Error:  "invalid request", // update msg
			})

			return
		}

		log.Info("request body is valid!")

		hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate hash", slog.Any("err", err))

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error", // update msg
			})

			return
		}

		log.Info("hash is generated!")

		req.Password = string(hashPass)

		id, err := hrs.CreateUser(&req)
		if err != nil {
			if errors.Is(err, storage.ErrUserAlreadyExists) {
				log.Info("user with this email or phone already exists", slog.Any("err", err))

				render.JSON(w, r, Response{
					Status: http.StatusConflict,
					Error:  "user with this email or phone already exists",
				})

				return
			}

			log.Info("internal server error")

			render.JSON(w, r, Response{
				Status: http.StatusInternalServerError,
				Error:  "internal server error",
			})

			return
		}

		date := time.Now().Format("2006-01-02 15:04:05")

		data := fmt.Sprintf("SIGNUP: [%s] user: %s with id: %d sign up successfully", date, req.Email, id)

		_, err = fmt.Fprintf(file, data)
		if err != nil {
			log.Error("Failed to write file", slog.Any("data", data), slog.String("err", err.Error()))
		}
		_, _ = fmt.Fprintf(file, "\n")
		render.JSON(w, r, Response{
			Id:     id,
			Status: http.StatusOK,
		})
	}
}
