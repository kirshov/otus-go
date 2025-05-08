package internalhttp

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	http        *http.Server
	logger      Logger
	application Application
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	return &Server{
		http: &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger:      logger,
		application: app,
	}
}

func (s *Server) Start(ctx context.Context, address string) error {
	<-ctx.Done()
	_ = address

	/*s.http.Addr = address
	return s.http.ListenAndServe()*/
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func handler(w http.ResponseWriter, r *http.Request) {
	msg := strings.Builder{}
	msg.WriteString(r.RemoteAddr + " ")
	msg.WriteString(r.Method + " ")
	msg.WriteString(r.URL.String() + " ")
	msg.WriteString(r.Proto + " ")
	msg.WriteString(r.Header.Get("User-Agent"))

	w.Write([]byte(msg.String()))
	// todo залогировать в middleware после соответствующего урока
}
