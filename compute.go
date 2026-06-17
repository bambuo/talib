package talib

// Sum 返回浮点数切片的总和。
func Sum(data []float64) float64 {
	var total float64
	for _, v := range data {
		total += v
	}
	return total
}

// Mean 返回浮点数切片的算术平均值。
func Mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	return Sum(data) / float64(len(data))
}

// MaxFloat64 返回两个 float64 值中的最大值。
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// MinFloat64 返回两个 float64 值中的最小值。
func MinFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
