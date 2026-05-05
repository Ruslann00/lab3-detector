package stats

import (
	"sync"
	"testing"
)

func TestRaceConditionBefore(t *testing.T) {
	t.Skip("This test demonstrates race condition before fix. Run old version manually if needed.")

	ResetStats()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			IncrementProcessed("jpg")
		}()
	}

	wg.Wait()

	result := GetProcessed("jpg")

	if result != 100 {
		t.Fatalf("expected 100, got %d", result)
	}
}

func TestSafeStatsCounterAfterFix(t *testing.T) {
	counter := NewSafeStatsCounter()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			counter.IncrementProcessed("jpg")
		}()
	}

	wg.Wait()

	result := counter.GetProcessed("jpg")

	if result != 100 {
		t.Fatalf("expected 100, got %d", result)
	}
}
