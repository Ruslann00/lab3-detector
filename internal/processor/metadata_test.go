package processor

import "testing"

func TestProcessedCount(t *testing.T) {
	ResetProcessedCount()

	if got := GetProcessedCount(); got != 0 {
		t.Fatalf("expected initial processed count 0, got %d", got)
	}

	processImage(1)
	processImage(2)
	processImage(3)

	if got := GetProcessedCount(); got != 3 {
		t.Fatalf("expected processed count 3, got %d", got)
	}
}

func TestBuildImageDataMatchesRegexp(t *testing.T) {
	data := buildImageData(1)

	if !imageRegexp.MatchString(data) {
		t.Fatalf("expected image data to match regexp, got %q", data)
	}
}

func TestSaveToCacheLimit(t *testing.T) {
	leakCacheMu.Lock()
	LeakCache = make(map[string][]byte, maxCacheSize)
	leakCacheMu.Unlock()

	for i := 0; i < maxCacheSize+100; i++ {
		saveToCache()
	}

	leakCacheMu.Lock()
	cacheLength := len(LeakCache)
	leakCacheMu.Unlock()

	if cacheLength > maxCacheSize {
		t.Fatalf("expected cache length to be <= %d, got %d", maxCacheSize, cacheLength)
	}
}
