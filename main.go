package main

import (
	"flowanalysis/pkg/db"
	"flowanalysis/pkg/inmemory"
	"flowanalysis/pkg/kafka"
	"flowanalysis/pkg/log"
	"flowanalysis/pkg/schedular"
	"flowanalysis/pkg/service"

	"github.com/labstack/echo/v4"
)

func init() {
	log.Print(log.Info, "Starting Count Schedular!!!")
	go schedular.StartCountSchedular()
	log.Print(log.Info, "Starting DB Cleaner Schedular!!!")
	go db.StartCleanupSchedular()
	log.Print(log.Info, "Starting Inmemory maps!!!")
	go inmemory.StartMap()
	log.Print(log.Info, "Starting Kafka!!!")
	go kafka.InitKafka()

}

func main() {

	e := echo.New()
	e.GET("/api/verve/accept", service.GetAcceptHandler)
	e.Start(":5010")
}
