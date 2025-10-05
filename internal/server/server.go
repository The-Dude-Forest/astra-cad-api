package server

import (
	"fmt"
	"go-auth/internal/database"
	"os"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *database.Service
	router *gin.Engine
	port   string
}

func NewServer() *Server {
	db := database.New()
	router := gin.Default()

	s := &Server{
		db:     db,
		router: router,
		port:   os.Getenv("PORT"),
	}

	s.setupMiddlewares()
	s.registerRoutes()

	return s
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) Router() *gin.Engine {
	return s.router
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%s", s.port)
	return s.router.Run(addr)
}
