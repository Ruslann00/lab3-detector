package processor

import "testing"

func BenchmarkProcessImageBefore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processImage(i)
	}
}

func BenchmarkProcessImageAfter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProcessImageFixed(i)
	}
}
