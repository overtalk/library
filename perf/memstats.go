package perf

import (
	"fmt"
	"runtime"
	"time"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigDefault

// GetMemStats gets runtime memstats
func GetMemStats(isPretty bool) ([]byte, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	if isPretty {
		memStats := fmt.Sprintf("Alloc: %.2fM\nTotal Alloc: %.2fM\n", float64(m.Alloc)/1024/1024/1024, float64(m.TotalAlloc)/1024/1024/1024)
		memStats += fmt.Sprintf("Sys: %.2fM\nHeapAlloc: %.2fM\n", float64(m.Sys)/1024/1024/1024, float64(m.HeapAlloc)/1024/1024/1024)
		memStats += fmt.Sprintf("Lookups: %d\nMallocs: %.2fM\n", m.Lookups, float64(m.Mallocs)/1024/1024/1024)
		memStats += fmt.Sprintf("Frees: %.2fM\n", float64(m.Frees)/1024/1024/1024)
		memStats += fmt.Sprintf("HeapSys: %.2fM\nHeapIdle: %.2fM\n", float64(m.HeapSys)/1024/1024/1024, float64(m.HeapIdle)/1024/1024/1024)
		memStats += fmt.Sprintf("HeapInuse: %.2fM\nHeapReleased: %.2fM\n", float64(m.HeapInuse)/1024/1024/1024, float64(m.HeapReleased)/1024/1024/1024)
		memStats += fmt.Sprintf("HeapObjects: %d\nStackInuse: %.2fM\n", m.HeapObjects, float64(m.StackInuse)/1024/1024/1024)
		memStats += fmt.Sprintf("StackSys: %.2fM\nGCSys: %.2fM\n", float64(m.StackSys)/1024/1024/1024, float64(m.GCSys)/1024/1024/1024)
		memStats += fmt.Sprintf("NextGC: %.2fM\nLastGC: %s\n", float64(m.NextGC)/1024/1024/1024, time.Unix(0, int64(m.LastGC)).Format("2006-01-02 15:04:05"))
		return []byte(memStats), nil
	}

	return json.Marshal(m)
}
