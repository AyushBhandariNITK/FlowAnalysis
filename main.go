package main

import (
	"flowanalysis/pkg/inmemory"
	"flowanalysis/pkg/service"

	"github.com/labstack/echo/v4"
)

func init() {
	go inmemory.StartMap()
}

func main() {

	e := echo.New()
	e.GET("/api/verve/accept", service.GetAcceptHandler)
	e.Start(":5010")
}
