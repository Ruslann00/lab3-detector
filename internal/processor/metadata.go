package processor

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

const maxCacheSize = 500

var (
	leakCacheMu sync.Mutex
	LeakCache   = make(map[string][]byte, maxCacheSize)

	imageRegexp = regexp.MustCompile(`^image_data_\d+_timestamp_\d+$`)

	processedCount atomic.Uint64
)

func RunWorkerPool(count int) {
	go logProcessedCount()

	for i := 0; i < count; i++ {
		go func(id int) {
			for {
				processImage(id)
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	select {}
}

func processImage(workerID int) {
	data := buildImageData(workerID)

	if !imageRegexp.MatchString(data) {
		return
	}

	saveToCache()
	processedCount.Add(1)
}

func buildImageData(workerID int) string {
	return fmt.Sprintf("image_data_%d_timestamp_%d", workerID, time.Now().UnixNano())
}

func saveToCache() {
	leakCacheMu.Lock()
	defer leakCacheMu.Unlock()

	if len(LeakCache) >= maxCacheSize {
		for key := range LeakCache {
			delete(LeakCache, key)
			break
		}
	}

	key := fmt.Sprintf("key_%d", time.Now().UnixNano())
	LeakCache[key] = make([]byte, 1024*10)
}

func logProcessedCount() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("processed images count: %d", GetProcessedCount())
	}
}

func GetProcessedCount() uint64 {
	return processedCount.Load()
}

func ResetProcessedCount() {
	processedCount.Store(0)
}
