package main

import (
	gatewayserver "quiz/internal/gateway-server"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/metrics", echoprometheus.NewHandler())

	srv, err := gatewayserver.New()
	if err != nil {
		e.Logger.Fatal(err)
	}

	srv.RegisterRoute(e.Group("/"))

	e.Logger.Fatal(e.Start(":8080"))
}
