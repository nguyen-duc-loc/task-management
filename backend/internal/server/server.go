package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/nguyen-duc-loc/task-management/backend/internal/store"
	"github.com/nguyen-duc-loc/task-management/backend/internal/token"
	"github.com/nguyen-duc-loc/task-management/backend/util"
)

type Server struct {
	Port       int
	router     *gin.Engine
	storage    store.Storage
	tokenMaker token.Maker
}

func (s *Server) RegisterRoutes() http.Handler {
	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%s", os.Getenv("FRONTEND_PORT"))},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	s.router.GET("/health", s.healthHandler)

	s.router.POST("/users", s.createUserHandler)
	s.router.POST("/users/login", s.loginUserHandler)

	authRoutes := s.router.Group("/").Use(authMiddleware(s.tokenMaker))
	authRoutes.GET("/tasks", s.getTasksHandler)
	authRoutes.GET("/tasks/:id", s.getTaskByIDHandler)
	authRoutes.POST("/tasks", s.createTaskHandler)

	return s.router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func NewServer(storage store.Storage) (*Server, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// ISO8601 validator
		v.RegisterValidation("iso8601", func(fl validator.FieldLevel) bool {
			_, err := time.Parse(time.RFC3339, fl.Field().String())
			return nil == err
		})
	}

	jwtConfig, err := util.LoadJWTConfig()
	if err != nil {
		return nil, err
	}

	tokenMaker, err := token.NewJWTMaker(jwtConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	newServer := &Server{
		Port:       port,
		router:     gin.Default(),
		storage:    storage,
		tokenMaker: tokenMaker,
	}
	return newServer, nil
}
