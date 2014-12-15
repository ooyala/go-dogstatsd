package dogstatsd

import (
	"runtime"
)

//Memory in use metrics
// bytes allocated and still in use.
func (c *Client) GetAllocatedMemory(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.Alloc
}

// bytes allocated and still in use (Main allocation Heap Statistics)
func (c *Client) GetHeapAllocation(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.HeapAlloc
}

// mcache structures (Inuse is bytes used now.)
func (c *Client) GetMCacheInuse(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.MCacheInuse
}

// mspan structures
func (c *Client) GetMSpanInuse(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.MSpanInuse
}

// bytes in non-idle span
func (c *Client) GetHeapInuse(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.HeapInuse
}

// Sytem memory allocations

// bytes obtained from system
func (c *Client) GetSys(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.Sys
}

// bytes obtained from system
func (c *Client) GetHeapSys(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.HeapSys
}

func (c *Client) GetMSpanSys(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.MSpanSys
}

func (c *Client) GetMCacheSys(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.MCacheSys
}

// Garbage collector statistics
// next collection will happen when HeapAlloc ≥ this amount
func (c *Client) GetNextGC(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.NextGC
}

// end time of last collection (nanoseconds since 1970)
func (c *Client) GetLastGC(m *runtime.MemStats) uint64 {
	runtime.ReadMemStats(m)
	return m.LastGC
}
