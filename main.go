package main

import (
	"flowanalysis/pkg"

	"github.com/labstack/echo/v4"
)

func init() {
	go pkg.StartMap()
}
func main() {

	e := echo.New()
	e.GET("/api/verve/accept", pkg.GetAcceptHandler)
	e.Start(":5010")
}
