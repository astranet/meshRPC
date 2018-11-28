// Package memstatsd implements a simple to use reporting tool that sends
// actual runtime.MemStats values and their diffs to a statsd server.
package memstatsd

import (
	"fmt"
	"runtime"
	"time"
)

type Statter interface {
	Timing(bucket string, d time.Duration)
	Gauge(bucket string, value int)
}

type MemStatsd struct {
	prefix string
	statsd Statter
	debug  bool
	env    string

	previous     *MemStats
	allocLatency time.Duration
}

func New(prefix string, envName string, statsd Statter, debug ...bool) MemStatsd {
	m := MemStatsd{
		prefix: prefix,
		statsd: statsd,
		env:    envName,
	}
	if len(debug) > 0 && debug[0] {
		m.debug = true
	}
	return m
}

func (m *MemStatsd) Run(d time.Duration) {
	tags := ",env=" + m.env
	t := time.NewTicker(d)
	go func() {
		for range t.C {
			m.pushMemStats(tags)
		}
	}()

	t2 := time.NewTicker(d)
	go func() {
		for range t2.C {
			m.pushAllocLatency(tags)
		}
	}()
}

func (m *MemStatsd) pushMemStats(tags string) {
	latest, delta := m.snapshotMemStats()
	if m.debug {
		fmt.Println(tags, "pushMemStats @", time.Now())
	}
	m.statsd.Gauge(m.prefix+"alloc"+tags, int(latest.Alloc))
	m.statsd.Gauge(m.prefix+"sys"+tags, int(latest.Sys))
	m.statsd.Gauge(m.prefix+"lookups"+tags, int(latest.Lookups))
	m.statsd.Gauge(m.prefix+"mallocs"+tags, int(latest.Mallocs))
	m.statsd.Gauge(m.prefix+"frees"+tags, int(latest.Frees))
	m.statsd.Gauge(m.prefix+"heap_alloc"+tags, int(latest.HeapAlloc))
	m.statsd.Gauge(m.prefix+"heap_sys"+tags, int(latest.HeapSys))
	m.statsd.Gauge(m.prefix+"heap_idle"+tags, int(latest.HeapIdle))
	m.statsd.Gauge(m.prefix+"heap_inuse"+tags, int(latest.HeapInuse))
	m.statsd.Gauge(m.prefix+"heap_released"+tags, int(latest.HeapReleased))
	m.statsd.Gauge(m.prefix+"heap_objects"+tags, int(latest.HeapObjects))
	m.statsd.Gauge(m.prefix+"num_gc"+tags, int(latest.NumGC))
	m.statsd.Gauge(m.prefix+"num_goroutine"+tags, latest.NumGoroutine)

	m.statsd.Timing(m.prefix+"pause_gc"+tags, latest.PauseGC)

	m.negativeGauge(m.prefix+"alloc.delta"+tags, int(delta.Alloc))
	m.negativeGauge(m.prefix+"sys.delta"+tags, int(delta.Sys))
	m.negativeGauge(m.prefix+"lookups.delta"+tags, int(delta.Lookups))
	m.negativeGauge(m.prefix+"mallocs.delta"+tags, int(delta.Mallocs))
	m.negativeGauge(m.prefix+"frees.delta"+tags, int(delta.Frees))
	m.negativeGauge(m.prefix+"heap_alloc.delta"+tags, int(delta.HeapAlloc))
	m.negativeGauge(m.prefix+"heap_sys.delta"+tags, int(delta.HeapSys))
	m.negativeGauge(m.prefix+"heap_idle.delta"+tags, int(delta.HeapIdle))
	m.negativeGauge(m.prefix+"heap_inuse.delta"+tags, int(delta.HeapInuse))
	m.negativeGauge(m.prefix+"heap_released.delta"+tags, int(delta.HeapReleased))
	m.negativeGauge(m.prefix+"heap_objects.delta"+tags, int(delta.HeapObjects))
	m.negativeGauge(m.prefix+"num_gc.delta"+tags, int(delta.NumGC))
	m.negativeGauge(m.prefix+"num_goroutine.delta"+tags, delta.NumGoroutine)
}

func (m *MemStatsd) negativeGauge(name string, v int) {
	if v < 0 {
		// https://github.com/cactus/go-statsd-client/issues/16
		m.statsd.Gauge(name, 0)
	}
	m.statsd.Gauge(name, v)
}

func (m *MemStatsd) pushAllocLatency(tags string) {
	latency, _ := m.snapshotAllocLatency()
	if m.debug {
		fmt.Println("pushAllocLatency @", time.Now())
	}
	m.statsd.Timing(m.prefix+"alloc_latency"+tags, latency)
}

type MemStats struct {
	// General stats
	Alloc   uint64 // bytes allocated and not yet freed
	Sys     uint64 // bytes obtained from system (sum of XxxSys below)
	Lookups uint64 // number of pointer lookups
	Mallocs uint64 // number of mallocs
	Frees   uint64 // number of frees

	// Heap stats
	HeapAlloc    uint64 // bytes allocated and not yet freed (same as Alloc above)
	HeapSys      uint64 // bytes obtained from system
	HeapIdle     uint64 // bytes in idle spans
	HeapInuse    uint64 // bytes in non-idle span
	HeapReleased uint64 // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// GC stats
	NumGC   uint32
	PauseGC time.Duration

	// Misc
	NumGoroutine int
}

func (m *MemStatsd) snapshotAllocLatency() (latency time.Duration, delta time.Duration) {
	const wait = 100 * time.Millisecond
	const size = 10 * 1024

	start := time.Now()
	var _ = make([]byte, size)
	time.Sleep(wait)
	var _ = make([]byte, size)
	latency = time.Since(start) - wait
	if m.allocLatency == 0 {
		m.allocLatency = latency
		return
	}
	delta = latency - m.allocLatency
	m.allocLatency = latency
	return
}

func (m *MemStatsd) snapshotMemStats() (latest *MemStats, delta MemStats) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	latest = &MemStats{
		Alloc:        stats.Alloc,
		Sys:          stats.Sys,
		Lookups:      stats.Lookups,
		Mallocs:      stats.Mallocs,
		Frees:        stats.Frees,
		HeapAlloc:    stats.HeapAlloc,
		HeapSys:      stats.HeapSys,
		HeapIdle:     stats.HeapIdle,
		HeapInuse:    stats.HeapInuse,
		HeapReleased: stats.HeapReleased,
		HeapObjects:  stats.HeapObjects,
		NumGC:        stats.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
		PauseGC:      time.Duration(stats.PauseNs[(stats.NumGC+255)%256]),
	}

	if m.previous == nil {
		m.previous = latest
		return
	}
	delta = MemStats{
		Alloc:        latest.Alloc - m.previous.Alloc,
		Sys:          latest.Sys - m.previous.Sys,
		Lookups:      latest.Lookups - m.previous.Lookups,
		Mallocs:      latest.Mallocs - m.previous.Mallocs,
		Frees:        latest.Frees - m.previous.Frees,
		HeapAlloc:    latest.HeapAlloc - m.previous.HeapAlloc,
		HeapSys:      latest.HeapSys - m.previous.HeapSys,
		HeapIdle:     latest.HeapIdle - m.previous.HeapIdle,
		HeapInuse:    latest.HeapInuse - m.previous.HeapInuse,
		HeapReleased: latest.HeapReleased - m.previous.HeapReleased,
		HeapObjects:  latest.HeapObjects - m.previous.HeapObjects,
		NumGC:        latest.NumGC - m.previous.NumGC,
		NumGoroutine: latest.NumGoroutine - m.previous.NumGoroutine,
		PauseGC:      latest.PauseGC - m.previous.PauseGC,
	}
	m.previous = latest
	return
}
