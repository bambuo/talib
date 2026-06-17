package talib

import (
	"math"
)

// CCI 计算商品通道指数。
//
// 典型价格 = (最高价 + 最低价 + 收盘价) / 3
// 平均偏差 = 均值(|TP - SMA(TP, 周期)|)
// CCI = (TP - SMA(TP, 周期)) / (0.015 * 平均偏差)
func CCI(high, low, close []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(close), "CCI:period"); err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)

	// 计算典型价格
	tp := make([]float64, n)
	for i := 0; i < n; i++ {
		tp[i] = (high[i] + low[i] + close[i]) / 3.0
	}

	// 计算 TP 的 SMA
	smaTP, err := SMA(tp, period)
	if err != nil {
		return nil, err
	}

	// 每个周期的 CCI
	const factor = 0.015
	startIdx := period - 1
	for i := startIdx; i < n; i++ {
		// 平均偏差 = 均值(|TP - SMA(TP)|) 在窗口内
		var meanDev float64
		for j := i - period + 1; j <= i; j++ {
			meanDev += math.Abs(tp[j] - smaTP[i])
		}
		meanDev /= float64(period)

		if meanDev == 0 {
			out[i] = 0
		} else {
			out[i] = (tp[i] - smaTP[i]) / (factor * meanDev)
		}
	}

	return out, nil
}

// CCILookback 返回 CCI 输出中前导 NaN 值的数量。
func CCILookback(period int) int {
	return period - 1
}
