package perf

import (
	"runtime"

	"github.com/pbnjay/memory"
)

// GetMemUsage gets memory usage
func GetMemUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	totalMem := memory.TotalMemory()
	if totalMem > 0 {
		return float64(m.Alloc) / float64(totalMem)
	}

	return 1
}
