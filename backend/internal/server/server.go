package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
)

type Server struct {
	port    int
	storage store.Storage
}

func NewServer(storage store.Storage) *http.Server {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// ISO8601 validator
		v.RegisterValidation("iso8601", func(fl validator.FieldLevel) bool {
			_, err := time.Parse(time.RFC3339, fl.Field().String())
			return nil == err
		})

		// Future validator
		v.RegisterValidation("future", func(fl validator.FieldLevel) bool {
			t, ok := fl.Field().Interface().(time.Time)
			if !ok {
				return false
			}
			return t.After(time.Now())
		})
	}

	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	NewServer := &Server{
		port,
		storage,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
