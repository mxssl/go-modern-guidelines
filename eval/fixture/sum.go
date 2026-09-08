package legacy

// Sum adds the values.
func Sum(values []int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}
