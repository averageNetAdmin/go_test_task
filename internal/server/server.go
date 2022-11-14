package server

import "sync"

type Server struct {
}

var (
	server *Server
	once   sync.Once
)

func GetInstance() *Server {
	once.Do(func() {
		server = &Server{}
	})
	return server
}
