package talib

import (
	"math"
)

// BBandsResult 包含布林带的值。
type BBandsResult struct {
	Upper  []float64 // 上轨：中轨 + nbDevUp * 标准差
	Middle []float64 // 中轨：输入数据的移动平均
	Lower  []float64 // 下轨：中轨 - nbDevDn * 标准差
}

// BBands 计算布林带。
//
// 中轨 = MA(输入, 周期, maType)
// 标准差 = sqrt(sum((输入[i] - 中轨[i])^2) / 周期)
// 上轨 = 中轨 + nbDevUp * 标准差
// 下轨 = 中轨 - nbDevDn * 标准差
//
// 当 nbDevUp == nbDevDn 时，使用单一标准差计算（TA-Lib 约定）。
func BBands(input []float64, period int, nbDevUp, nbDevDn float64, maType MAType) (*BBandsResult, error) {
	if err := ValidateNumericInput(input, "BBands:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "BBands:period"); err != nil {
		return nil, err
	}

	n := len(input)

	// 计算中轨（移动平均）
	middle, err := MA(input, period, maType)
	if err != nil {
		return nil, err
	}

	upper := MakeOutput(n)
	lower := MakeOutput(n)

	startIdx := period - 1
	for i := startIdx; i < n; i++ {
		// 计算窗口内的总体标准差
		var sumSq float64
		for j := i - period + 1; j <= i; j++ {
			diff := input[j] - middle[i]
			sumSq += diff * diff
		}
		stdDev := math.Sqrt(sumSq / float64(period))

		upper[i] = middle[i] + nbDevUp*stdDev
		lower[i] = middle[i] - nbDevDn*stdDev
	}

	return &BBandsResult{
		Upper:  upper,
		Middle: middle,
		Lower:  lower,
	}, nil
}

// BBandsLookback 返回 BBands 输出中前导 NaN 值的总数。
func BBandsLookback(period int, maType MAType) int {
	return period - 1
}
