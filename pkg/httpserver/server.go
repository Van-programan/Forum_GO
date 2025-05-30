package httpserver

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine  *gin.Engine
	address string
}

func New(address string) *Server {
	s := &Server{address: address}
	s.Engine = gin.Default()
	return s
}

func (s *Server) Run() {
	fmt.Println(s.address)
	go s.Engine.Run(s.address)
}
