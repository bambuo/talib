package talib

import (
	"math"
)

// 方向运动（DM）类型。
// +DM = 正方向运动
// -DM = 负方向运动

// PLUS_DM 计算正方向运动。
//
// 如果 (high[i] - high[i-1]) > (low[i-1] - low[i]) 且 (high[i] - high[i-1]) > 0，
// 则 +DM = high[i] - high[i-1], 否则为 0。
func PLUS_DM(high, low []float64, period int) ([]float64, error) {
	n := len(high)
	out := MakeOutput(n)

	if n <= 1 {
		return out, nil
	}

	dm := make([]float64, n)
	for i := 1; i < n; i++ {
		upMove := high[i] - high[i-1]
		downMove := low[i-1] - low[i]
		if upMove > downMove && upMove > 0 {
			dm[i] = upMove
		}
	}

	// Wilder 平滑
	if period >= n {
		return out, nil
	}
	startIdx := period
	var sum float64
	for i := 1; i <= period; i++ {
		sum += dm[i]
	}
	avg := sum / float64(period)
	out[startIdx] = avg

	for i := startIdx + 1; i < n; i++ {
		avg = (avg*float64(period-1) + dm[i]) / float64(period)
		out[i] = avg
	}

	return out, nil
}

// PLUS_DMLookback 返回 +DM 的 lookback。
func PLUS_DMLookback(period int) int {
	return period
}

// MINUS_DM 计算负方向运动。
//
// 如果 (low[i-1] - low[i]) > (high[i] - high[i-1]) 且 (low[i-1] - low[i]) > 0，
// 则 -DM = low[i-1] - low[i], 否则为 0。
func MINUS_DM(high, low []float64, period int) ([]float64, error) {
	n := len(high)
	out := MakeOutput(n)

	if n <= 1 {
		return out, nil
	}

	dm := make([]float64, n)
	for i := 1; i < n; i++ {
		upMove := high[i] - high[i-1]
		downMove := low[i-1] - low[i]
		if downMove > upMove && downMove > 0 {
			dm[i] = downMove
		}
	}

	if period >= n {
		return out, nil
	}
	startIdx := period
	var sum float64
	for i := 1; i <= period; i++ {
		sum += dm[i]
	}
	avg := sum / float64(period)
	out[startIdx] = avg

	for i := startIdx + 1; i < n; i++ {
		avg = (avg*float64(period-1) + dm[i]) / float64(period)
		out[i] = avg
	}

	return out, nil
}

// MINUS_DMLookback 返回 -DM 的 lookback。
func MINUS_DMLookback(period int) int {
	return period
}

// PLUS_DI 计算正方向指标。
//
// +DI = (+DM / 真实波幅) * 100
func PLUS_DI(high, low, close []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}

	plusDM, err := PLUS_DM(high, low, period)
	if err != nil {
		return nil, err
	}

	atr, err := ATR(high, low, close, period)
	if err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)
	startIdx := period
	for i := startIdx; i < n; i++ {
		if atr[i] == 0 {
			out[i] = 0
		} else {
			out[i] = plusDM[i] / atr[i] * 100.0
		}
	}

	return out, nil
}

// PLUS_DILookback 返回 +DI 的 lookback。
func PLUS_DILookback(period int) int {
	return period
}

// MINUS_DI 计算负方向指标。
//
// -DI = (-DM / 真实波幅) * 100
func MINUS_DI(high, low, close []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}

	minusDM, err := MINUS_DM(high, low, period)
	if err != nil {
		return nil, err
	}

	atr, err := ATR(high, low, close, period)
	if err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)
	startIdx := period
	for i := startIdx; i < n; i++ {
		if atr[i] == 0 {
			out[i] = 0
		} else {
			out[i] = minusDM[i] / atr[i] * 100.0
		}
	}

	return out, nil
}

// MINUS_DILookback 返回 -DI 的 lookback。
func MINUS_DILookback(period int) int {
	return period
}

// ADX 计算平均方向运动指数。
//
// DX  = |(+DI) - (-DI)| / ((+DI) + (-DI)) * 100
// ADX = Wilder 的 EMA(DX, period)
func ADX(high, low, close []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}

	plusDI, err := PLUS_DI(high, low, close, period)
	if err != nil {
		return nil, err
	}

	minusDI, err := MINUS_DI(high, low, close, period)
	if err != nil {
		return nil, err
	}

	n := len(close)
	dx := MakeOutput(n)

	// 计算 DX
	for i := 0; i < n; i++ {
		if IsNaN(plusDI[i]) || IsNaN(minusDI[i]) {
			continue
		}
		sum := plusDI[i] + minusDI[i]
		if sum == 0 {
			dx[i] = 0
		} else {
			dx[i] = math.Abs(plusDI[i]-minusDI[i]) / sum * 100.0
		}
	}

	// ADX = Wilder 的 EMA 形式的 DX（与 wma 类似平滑）
	startIdx := 2 * period
	if startIdx >= n {
		return dx, nil
	}

	out := MakeOutput(n)
	// 第一个 ADX = 从 'period' 开始的前 'period' 个 DX 值的平均值
	var sumDX float64
	for i := period; i < startIdx; i++ {
		sumDX += dx[i]
	}
	out[startIdx] = sumDX / float64(period)

	// Wilder's smooth
	for i := startIdx + 1; i < n; i++ {
		if IsNaN(dx[i]) {
			out[i] = out[i-1]
		} else {
			out[i] = (out[i-1]*float64(period-1) + dx[i]) / float64(period)
		}
	}

	return out, nil
}

// ADXLookback 返回 ADX 的 lookback。
func ADXLookback(period int) int {
	return 2 * period
}
