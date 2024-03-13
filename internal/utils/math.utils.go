package utils

func Min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func Clamp(val int, min int, max int) int {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}
