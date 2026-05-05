package stats

import "sync"

type SafeStatsCounter struct {
	mu    sync.RWMutex
	stats map[string]int
}

func NewSafeStatsCounter() *SafeStatsCounter {
	return &SafeStatsCounter{
		stats: make(map[string]int),
	}
}

func (c *SafeStatsCounter) IncrementProcessed(imageType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats[imageType]++
}

func (c *SafeStatsCounter) GetProcessed(imageType string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.stats[imageType]
}

func (c *SafeStatsCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.stats = make(map[string]int)
}
