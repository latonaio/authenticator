package server

import (
	"context"
	"fmt"

	"bitbucket.org/latonaio/authenticator/configs"
	"bitbucket.org/latonaio/authenticator/internal/route"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	Server  *echo.Echo
	Context context.Context
	Port    string
}

type Server interface {
	Start(errC chan error)
	Shutdown(ctx context.Context) error
}

func New(ctx context.Context, cfg configs.Configs) Server {
	cfgs := cfg.Get()
	// Echo instance
	e := echo.New()
	// Routes
	route.RegisterRouting(e)
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// use echo default logger
	e.Use(middleware.Logger())
	// error handling
	e.HTTPErrorHandler = custmres.CustomHTTPErrorHandler

	return &server{
		Server:  e,
		Context: ctx,
		Port:    fmt.Sprintf(":%v", cfgs.Server.Port),
	}
}

func (s *server) Start(errC chan error) {
	err := s.Server.Start(s.Port)
	errC <- err
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
