package rest

import (
	"auth/internal/config"
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	server *http.Server
	log    *logrus.Logger
}

func NewServer(log *logrus.Logger, handler http.Handler, cfg *config.Configs) *Server {
	return &Server{
		server: &http.Server{
			Addr:    cfg.GetHTTPAddr(),
			Handler: handler,
		},
		log: log,
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
