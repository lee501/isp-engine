package server

import (
	"fmt"
	"isp-engine/route"
	"net/http"
	"os"
)

type Server struct {
	Listen string
}

func NewServer(listen string) *Server {
	return &Server{
		Listen: listen,
	}
}

func (s *Server) Run() {
	router := route.InitRouter()
	if err := http.ListenAndServe(s.Listen, router); err != nil {
		fmt.Errorf("server start fail: %s\n", err)
		os.Exit(-1)
	}
}
