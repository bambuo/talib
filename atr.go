package talib

import (
	"math"
)

// ATR 使用 Wilder 平滑方法计算平均真实波幅。
//
// 真实波幅 = max(最高 - 最低, |最高 - 前收|, |最低 - 前收|)
// 第一个 ATR = SMA(真实波幅, 周期)
// 后续 ATR = (前一个ATR * (周期-1) + TR) / 周期
func ATR(high, low, close []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(close), "ATR:period"); err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)

	// 计算每个周期的真实波幅（TR[0] 未定义，使用最高-最低）
	trSlice := make([]float64, n)
	for i := 0; i < n; i++ {
		tr := high[i] - low[i]
		if i > 0 {
			prevClose := close[i-1]
			tr = math.Max(tr, math.Abs(high[i]-prevClose))
			tr = math.Max(tr, math.Abs(low[i]-prevClose))
		}
		trSlice[i] = tr
	}

	// 第一个 ATR：前 period 个 TR 值的 SMA
	startIdx := period
	var sum float64
	for i := 0; i < period; i++ {
		sum += trSlice[i]
	}
	out[startIdx] = sum / float64(period)

	// 其余使用 Wilder 平滑 ATR
	for i := startIdx + 1; i < n; i++ {
		out[i] = (out[i-1]*float64(period-1) + trSlice[i]) / float64(period)
	}

	return out, nil
}

// ATRLookback 返回 ATR 输出中前导 NaN 值的数量。
func ATRLookback(period int) int {
	return period
}

// TRANGE 计算每个周期的真实波幅。
//
// TR[i] = max(最高-最低, |最高-前收|, |最低-前收|)
// TR[0] = high[0] - low[0]（无前收盘价可用）
func TRANGE(high, low, close []float64) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}

	n := len(close)
	out := make([]float64, n)

	for i := 0; i < n; i++ {
		tr := high[i] - low[i]
		if i > 0 {
			prevClose := close[i-1]
			tr = MaxFloat64(tr, math.Abs(high[i]-prevClose))
			tr = MaxFloat64(tr, math.Abs(low[i]-prevClose))
		}
		out[i] = tr
	}

	return out, nil
}

// NATR 计算归一化平均真实波幅。
//
// NATR = (ATR / 收盘价) * 100
func NATR(high, low, close []float64, period int) ([]float64, error) {
	atr, err := ATR(high, low, close, period)
	if err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)
	startIdx := period
	for i := startIdx; i < n; i++ {
		if close[i] != 0 && !math.IsNaN(atr[i]) {
			out[i] = (atr[i] / close[i]) * 100.0
		}
	}

	return out, nil
}

// NATRLookback 返回 NATR 输出中前导 NaN 值的数量。
func NATRLookback(period int) int {
	return period
}
