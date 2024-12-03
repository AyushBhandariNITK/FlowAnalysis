package utils

import (
	"flowanalysis/pkg/log"
	"os"
	"strings"
)

var (
	inMemoryUse bool
)

func init() {
	useInMemory := GetEnv("INMEMORY", "false")
	if strings.Contains(useInMemory, "false") {
		inMemoryUse = false
	} else {
		inMemoryUse = true
	}
	log.Print(log.Info, "Using inmemory: %+v", inMemoryUse)
}

func UseInMemory() bool {
	return inMemoryUse
}
func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
