package http

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"hash-api/internal/domain"
	"log"
	"net/http"
)

type hashGetter interface {
	Handle(ctx context.Context) (domain.Hash, error)
}

type Server struct {
	addr string
	srv  *echo.Echo
}

func NewServer(addr string, hashGetter hashGetter) *Server {
	var server = Server{addr: addr, srv: echo.New()}

	server.srv.GET("/hash", func(c echo.Context) error {
		hash, err := hashGetter.Handle(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Response().Header().Set("Expires", hash.ExpiresAt().Format(http.TimeFormat))
		return c.JSON(http.StatusOK, hash)
	})

	return &server
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		if err := s.srv.Start(s.addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("echo start: %s", err)
		}
	}()

	<-ctx.Done()

	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("echo shutdown: %s", err)
	} else {
		log.Println("echo server stopped")
	}
}
