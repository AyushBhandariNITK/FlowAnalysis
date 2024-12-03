package schedular

import (
	"context"
	"flowanalysis/pkg/inmemory"
	"flowanalysis/pkg/kafka"
	"flowanalysis/pkg/log"
	"flowanalysis/pkg/service"
	"flowanalysis/pkg/utils"

	"time"
)

const (
	MAP_SCHEDULAR_IN_SECONDS = 60 * time.Second
)

func StartCountSchedular() {
	ticker := time.NewTicker(MAP_SCHEDULAR_IN_SECONDS)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if utils.UseInMemory() {
				inmemory.InMemoryCount()
			} else {
				count, _ := service.CountEntries(time.Now())
				data := map[string]interface{}{
					"timestamp":    time.Now().Format(time.RFC3339),
					"unique_count": count,
				}
				log.Print(log.Info, "Send message to kafka %+v", data)
				kafka.SendMessage(context.Background(), kafka.FLOW_UNIQUE_TOPIC, data)
			}

		}
	}
}
