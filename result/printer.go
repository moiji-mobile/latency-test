package result

import (
	"fmt"
	"os"
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

func PrintSummary(result Result) {
	minVal := min(result.Items)
	maxVal := max(result.Items)
	avgVal := avg(result.Items)

	fmt.Printf("AVG(%vns) MIN(%vns) MAX(%vns)\n", avgVal, minVal, maxVal)
}

func PrintCSV(pktSize int, result Result, dest string) {
	file, err := os.Create(dest)
	if err != nil {
		panic("Can't open the result file")
	}
	defer file.Close()

	for _, item := range result.Items {
		fmt.Fprintf(file, "%v,%v,%v\n", pktSize, item.Message, item.Roundtrip)
	}
}
