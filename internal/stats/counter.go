package stats

var GlobalStats = make(map[string]int)

func IncrementProcessed(imageType string) {
	GlobalStats[imageType]++
}

func GetProcessed(imageType string) int {
	return GlobalStats[imageType]
}

func ResetStats() {
	GlobalStats = make(map[string]int)
}
