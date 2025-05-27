package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/kirshov/otus-go/hw12_13_14_15_calendar/internal/storage"
)

const handlerTimeout time.Duration = 5 * time.Second

type Server struct {
	http        *http.Server
	application Application
}

type Application interface {
	GetLogger() logger.Logger
	GetStorage() storage.Storage
}

func NewServer(app Application) *Server {
	mux := http.NewServeMux()

	routes := getRoutes(app)
	for _, route := range routes {
		mux.HandleFunc(route.URL, route.Handler)
	}

	muxWithLogging := loggingMiddleware(app, mux)
	return &Server{
		http: &http.Server{
			Handler:           muxWithLogging,
			ReadHeaderTimeout: handlerTimeout,
		},
		application: app,
	}
}

func (s *Server) Start(address string) error {
	s.http.Addr = address
	s.application.GetLogger().Info("http server started")

	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
