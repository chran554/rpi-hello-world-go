package main

import (
	"fmt"
	"strings"
)

func main() {
	percent := 0.9599

	rate := 1.0 / percent
	rateAccumulator := 0.0

	var text strings.Builder

	for j := 0; j < 1000000; j++ {
		rateAccumulator += 1.0

		if rateAccumulator >= rate {
			rateAccumulator -= rate
			text.WriteString("#")
		} else {
			text.WriteString(" ")
		}
	}

	result := text.String()
	amount := len(result)
	count1 := strings.Count(result, "#")
	count2 := strings.Count(result, " ")
	percent1 := 100 * float64(count1) / float64(amount)
	percent2 := 100 * float64(count2) / float64(amount)

	fmt.Printf("[%s...]\n", result[:100])
	fmt.Printf("Amount samples: %d\n", amount)
	fmt.Println("Percent target: ", percent*100)
	fmt.Printf("'#' %d  (%.2f%%)\n", count1, percent1)
	fmt.Printf("' ' %d  (%.2f%%)\n", count2, percent2)

	fmt.Println()
}
