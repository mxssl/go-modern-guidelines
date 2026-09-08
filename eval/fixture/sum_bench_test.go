package legacy

import "testing"

func BenchmarkSum(b *testing.B) {
	// One distinct input per iteration so the compiler cannot hoist the call.
	inputs := make([][]int, b.N)
	for i := 0; i < b.N; i++ {
		inputs[i] = []int{i, i + 1, i + 2, i + 3}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Sum(inputs[i])
	}
}
