package status

import "testing"

func BenchmarkPercentageCalculation(b *testing.B) {
	calculatePercentage(93210.0, 412941942.0)
}
