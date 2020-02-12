package goutils

// Returns the greater of int values.
func MaxInt(a ...int) int {
	if len(a) == 0 {
		panic("empty data")
	}
	max := a[0]
	for i := 1; i < len(a); i++ {
		if max < a[i] {
			max = a[i]
		}
	}
	return max
}

// Returns the smaller of int values.
func MinInt(a ...int) int {
	if len(a) == 0 {
		panic("empty data")
	}
	max := a[0]
	for i := 1; i < len(a); i++ {
		if max > a[i] {
			max = a[i]
		}
	}
	return max
}
