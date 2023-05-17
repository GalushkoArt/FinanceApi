package server

import (
	"FinanceApi/pkg/config"
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + config.Conf.Server.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Duration(config.Conf.Server.ReadTimeoutSec) * time.Second,
		WriteTimeout:   time.Duration(config.Conf.Server.WriteTimeoutSec) * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
