package httpchi

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

// TODO: подумать над options для сервера
func NewServer(port string, router chi.Router) *Server {
	return &Server{httpServer: &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}}
}

func (s *Server) RunServer() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) GracefulShutdownServer(ctx context.Context) {
	<-ctx.Done()

	slog.Info("shutting down server...")

	shutdownCtx := context.Background()
	if err := s.httpServer.Shutdown(shutdownCtx); err != http.ErrServerClosed && err != nil {
		slog.Error(fmt.Sprintf("error on server shutting down: %s", err.Error()))
		return
	}

	slog.Info("server stoped gracefully")
}
