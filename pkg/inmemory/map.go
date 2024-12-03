package inmemory

import (
	"context"
	"flowanalysis/pkg/kafka"
	"flowanalysis/pkg/log"
	"flowanalysis/pkg/service"
	"sync"
	"time"

	"sync/atomic"

	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	MAP_SCHEDULAR_IN_SECONDS = 60 * time.Second
)

func init() {
	MapA = NewFlowMap()
	MapB = NewFlowMap()
	activeMap = MapA
	nonactiveMap = MapB
	cond = sync.NewCond(&sync.Mutex{})
	atomic.StoreInt32(&isActive, 1)
}

var (
	MapA         *FlowMap
	MapB         *FlowMap
	activeMap    *FlowMap
	nonactiveMap *FlowMap
	isActive     int32 // Atomic flag to indicate the active map (1 for active, 0 for inactive)
	cond         *sync.Cond
)

type IFlowMap interface {
	Set(string, string)
	Clear()
	Count() int
}

type FlowMap struct {
	data cmap.ConcurrentMap[string, string]
}

func NewFlowMap() *FlowMap {
	return &FlowMap{
		data: cmap.New[string](),
	}
}

func (f *FlowMap) Set(key, value string) {
	for atomic.LoadInt32(&isActive) == 0 {
		cond.Wait()
	}
	f.data.Set(key, value)
}

func (f *FlowMap) Clear() {
	for atomic.LoadInt32(&isActive) == 0 {
		cond.Wait()
	}
	f.data.Clear()
}

func (f *FlowMap) Count() int {
	for atomic.LoadInt32(&isActive) == 0 {
		cond.Wait()
	}
	return f.data.Count()
}

func GetActiveFlowMap() *FlowMap {
	return activeMap
}

func GetNonActiveFlowMap() *FlowMap {
	return nonactiveMap
}

func StartMap() {
	ticker := time.NewTicker(MAP_SCHEDULAR_IN_SECONDS)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			atomic.StoreInt32(&isActive, 0)

			cond.L.Lock()
			if activeMap == MapA {
				activeMap = MapB
				nonactiveMap = MapA
			} else {
				activeMap = MapA
				nonactiveMap = MapB
			}
			atomic.StoreInt32(&isActive, 1)
			cond.Broadcast()
			cond.L.Unlock()

			count, _ := service.CountEntries(time.Now())
			// count := 10
			data := map[string]interface{}{
				"timestamp":    time.Now().Format(time.RFC3339),
				"unique_count": count,
			}
			kafka.SendMessage(context.Background(), kafka.FLOW_UNIQUE_TOPIC, data)
			log.Print(log.File, "No.of unique entries %d \n", count)
			nonactiveMap.Clear()
		}
	}
}
