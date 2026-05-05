package processor

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

var LeakCache = make(map[string][]byte)

var leakMu sync.Mutex

func RunWorkerPool(count int) {
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
	data := fmt.Sprintf("image_data_%d_timestamp_%d", workerID, time.Now().UnixNano())

	matched, _ := regexp.MatchString(`^image_data_\d+_timestamp_\d+$`, data)

	if matched {
		key := fmt.Sprintf("key_%d", time.Now().UnixNano())

		leakMu.Lock()
		LeakCache[key] = make([]byte, 1024*10)
		leakMu.Unlock()
	}
}
