package pkg

import (
	"sync"
	"time"

	"flowanalysis/pkg/log"
	"sync/atomic"

	cmap "github.com/orcaman/concurrent-map/v2"
)

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

func init() {
	MapA = NewFlowMap()
	MapB = NewFlowMap()
	activeMap = MapA
	nonactiveMap = MapB
	cond = sync.NewCond(&sync.Mutex{})
	atomic.StoreInt32(&isActive, 1)
}

func StartMap() {
	ticker := time.NewTicker(60 * time.Second)
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
			log.Print(log.File, "No.of unique entries %d \n", nonactiveMap.Count())
			nonactiveMap.Clear()
		}
	}
}
