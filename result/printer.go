package result

import (
	"fmt"
)

func min(items []Item) int64 {
	minVal := items[0].Roundtrip

	for _, item := range items[1:] {
		if item.Roundtrip < minVal {
			minVal = item.Roundtrip
		}
	}
	return minVal
}

func max(items []Item) int64 {
	maxVal := items[0].Roundtrip

	for _, item := range items[1:] {
		if item.Roundtrip > maxVal {
			maxVal = item.Roundtrip
		}
	}
	return maxVal
}

func avg(items []Item) int64 {
	sum := int64(0)
	for _, item := range items {
		sum += item.Roundtrip
	}
	return sum/int64(len(items))
}

func PrintResult(result Result) {
	minVal := min(result.Items)
	maxVal := max(result.Items)
	avgVal := avg(result.Items)

	fmt.Printf("AVG(%vns) MIN(%vns) MAX(%vns)\n", avgVal, minVal, maxVal)
}
