package processor

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

var optimizedRegexp = regexp.MustCompile(`^image_data_\d+_timestamp_\d+$`)

type SafeCache struct {
	mu      sync.Mutex
	data    map[string][]byte
	maxSize int
}

func NewSafeCache(maxSize int) *SafeCache {
	return &SafeCache{
		data:    make(map[string][]byte),
		maxSize: maxSize,
	}
}

func (c *SafeCache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) >= c.maxSize {
		for oldKey := range c.data {
			delete(c.data, oldKey)
			break
		}
	}

	c.data[key] = value
}

func (c *SafeCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.data)
}

var FixedCache = NewSafeCache(1000)

func ProcessImageFixed(workerID int) {
	data := fmt.Sprintf("image_data_%d_timestamp_%d", workerID, time.Now().UnixNano())

	matched := optimizedRegexp.MatchString(data)

	if matched {
		key := fmt.Sprintf("key_%d", time.Now().UnixNano())
		FixedCache.Set(key, make([]byte, 1024*10))
	}
}

func RunWorkerPoolFixed(count int) {
	for i := 0; i < count; i++ {
		go func(id int) {
			for {
				ProcessImageFixed(id)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	select {}
}
