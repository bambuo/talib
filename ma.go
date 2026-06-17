package talib

import (
	"math"
)

// SMA 计算简单移动平均。
//
// 对于每个输出元素 i（从 period-1 开始）：
//
//	SMA[i] = (input[i] + input[i-1] + ... + input[i-period+1]) / period
//
// 前 (period-1) 个元素为 NaN，匹配 TA-Lib 约定。
func SMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "SMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "SMA:period"); err != nil {
		return nil, err
	}

	out := MakeOutput(len(input))

	// 查找第一个有效窗口（全部非 NaN）
	startIdx := period - 1
	for startIdx < len(input) {
		allValid := true
		var sum float64
		for j := startIdx - period + 1; j <= startIdx; j++ {
			if IsNaN(input[j]) {
				allValid = false
				break
			}
			sum += input[j]
		}
		if allValid {
			out[startIdx] = sum / float64(period)
			break
		}
		startIdx++
	}

	if startIdx >= len(input) {
		return out, nil
	}

	// 滑动窗口
	for i := startIdx + 1; i < len(input); i++ {
		if IsNaN(input[i]) {
			out[i] = out[i-1] // 向前延续
			continue
		}
		prev := input[i-period]
		if IsNaN(prev) {
			// 从窗口中的有效值重新计算
			var sum float64
			var count int
			for j := i - period + 1; j <= i; j++ {
				if !IsNaN(input[j]) {
					sum += input[j]
					count++
				}
			}
			if count > 0 {
				out[i] = sum / float64(count)
			} else {
				out[i] = out[i-1]
			}
		} else {
			out[i] = out[i-1] + (input[i]-prev)/float64(period)
		}
	}

	return out, nil
}

// SMALookback 返回 SMA 输出中前导 NaN 值的数量。
func SMALookback(period int) int {
	return period - 1
}

// EMA 使用 Wilder 方法计算指数移动平均。
//
// 第一个有效值（索引 period-1）是前 period 个元素的 SMA。
// 后续值：EMA[i] = EMA[i-1] + k * (input[i] - EMA[i-1])
//
// 其中 k = 2 / (period + 1)。
func EMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "EMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "EMA:period"); err != nil {
		return nil, err
	}

	out := MakeOutput(len(input))

	// 乘数因子：k = 2 / (period + 1)
	k := 2.0 / float64(period+1)

	// 查找第一个有效 SMA 种子：跳过输入中的 NaN 值
	startIdx := period - 1
	for startIdx < len(input) {
		// 检查窗口 [startIdx-period+1 .. startIdx] 中的所有值是否非 NaN
		allValid := true
		for j := startIdx - period + 1; j <= startIdx; j++ {
			if IsNaN(input[j]) {
				allValid = false
				break
			}
		}
		if allValid {
			break
		}
		startIdx++
	}

	if startIdx >= len(input) {
		return out, nil // 全为 NaN
	}

	// 第一个值是第一个有效窗口的 SMA
	var sum float64
	for j := startIdx - period + 1; j <= startIdx; j++ {
		sum += input[j]
	}
	out[startIdx] = sum / float64(period)

	// 递归 EMA（跳过 NaN 输入）
	for i := startIdx + 1; i < len(input); i++ {
		if IsNaN(input[i]) {
			out[i] = out[i-1] // 向前延续
		} else {
			out[i] = out[i-1] + k*(input[i]-out[i-1])
		}
	}

	return out, nil
}

// EMALookback 返回 EMA 输出中前导 NaN 值的数量。
func EMALookback(period int) int {
	return period - 1
}

// WMA 计算加权移动平均。
//
// 窗口内的每个值乘以一个权重：
//
//	WMA[i] = sum(w[j] * input[i-period+1+j]) / sum(w)
//
// 其中 w = [1, 2, 3, ..., period]（线性加权）。
func WMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "WMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "WMA:period"); err != nil {
		return nil, err
	}

	out := MakeOutput(len(input))

	// 分母：权重之和 1 + 2 + ... + period
	denom := float64(period * (period + 1) / 2)

	for i := period - 1; i < len(input); i++ {
		var sum float64
		weight := 1.0
		for j := i - period + 1; j <= i; j++ {
			sum += input[j] * weight
			weight++
		}
		out[i] = sum / denom
	}

	return out, nil
}

// WMALookback 返回 WMA 输出中前导 NaN 值的数量。
func WMALookback(period int) int {
	return period - 1
}

// DEMA 计算双重指数移动平均。
//
// DEMA = 2 * EMA(input, period) - EMA(EMA(input, period), period)
func DEMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "DEMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "DEMA:period"); err != nil {
		return nil, err
	}

	ema1, err := EMA(input, period)
	if err != nil {
		return nil, err
	}

	ema2, err := EMA(ema1, period)
	if err != nil {
		return nil, err
	}

	out := MakeOutput(len(input))
	startIdx := 2*period - 2 // 需要 2 个完整的 EMA lookback
	for i := startIdx; i < len(input); i++ {
		if !math.IsNaN(ema1[i]) && !math.IsNaN(ema2[i]) {
			out[i] = 2*ema1[i] - ema2[i]
		}
	}

	return out, nil
}

// DEMALookback 返回 DEMA 输出中前导 NaN 值的数量。
func DEMALookback(period int) int {
	return 2*period - 2
}

// TEMA 计算三重指数移动平均。
//
// EMA1 = EMA(input, period)
// EMA2 = EMA(EMA1, period)
// EMA3 = EMA(EMA2, period)
// TEMA = 3*EMA1 - 3*EMA2 + EMA3
func TEMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "TEMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "TEMA:period"); err != nil {
		return nil, err
	}

	ema1, err := EMA(input, period)
	if err != nil {
		return nil, err
	}

	ema2, err := EMA(ema1, period)
	if err != nil {
		return nil, err
	}

	ema3, err := EMA(ema2, period)
	if err != nil {
		return nil, err
	}

	out := MakeOutput(len(input))
	startIdx := 3*period - 3 // 需要 3 个完整的 EMA lookback
	for i := startIdx; i < len(input); i++ {
		if !math.IsNaN(ema1[i]) && !math.IsNaN(ema2[i]) && !math.IsNaN(ema3[i]) {
			out[i] = 3*ema1[i] - 3*ema2[i] + ema3[i]
		}
	}

	return out, nil
}

// TEMALookback 返回 TEMA 输出中前导 NaN 值的数量。
func TEMALookback(period int) int {
	return 3*period - 3
}

// TRIMA 计算三角移动平均。
//
// TRIMA 是 SMA 的 SMA：
//   - N = floor((period + 1) / 2)
//   - 如果 period 为奇数：SMA(SMA(input, N), N)
//   - 如果 period 为偶数：SMA(SMA(input, N), N+1)
func TRIMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "TRIMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "TRIMA:period"); err != nil {
		return nil, err
	}

	n := (period + 1) / 2 // floor((period+1)/2)

	sma1, err := SMA(input, n)
	if err != nil {
		return nil, err
	}

	var sma2 []float64
	if period%2 != 0 {
		// 奇数周期：第二个 SMA 使用相同的 n
		sma2, err = SMA(sma1, n)
	} else {
		// 偶数周期：第二个 SMA 使用 n+1
		sma2, err = SMA(sma1, n+1)
	}
	if err != nil {
		return nil, err
	}

	return sma2, nil
}

// TRIMALookback 返回 TRIMA 输出中前导 NaN 值的数量。
func TRIMALookback(period int) int {
	return period - 1
}

// MA 使用指定的 MA 类型计算移动平均。
func MA(input []float64, period int, maType MAType) ([]float64, error) {
	switch maType {
	case MASMA:
		return SMA(input, period)
	case MAEMA:
		return EMA(input, period)
	case MAWMA:
		return WMA(input, period)
	case MADEMA:
		return DEMA(input, period)
	case MATEMA:
		return TEMA(input, period)
	case MATRIMA:
		return TRIMA(input, period)
	case MAKAMA:
		return KAMA(input, period)
	case MAMAMA:
		result, err := MAMA(input, 0.5, 0.05)
		if err != nil {
			return nil, err
		}
		return result.MAMA, nil
	default:
		return SMA(input, period) // 回退
	}
}

// MALookback 返回给定 MA 类型和周期的 lookback。
func MALookback(period int, maType MAType) int {
	switch maType {
	case MASMA:
		return SMALookback(period)
	case MAEMA:
		return EMALookback(period)
	case MAWMA:
		return WMALookback(period)
	case MADEMA:
		return DEMALookback(period)
	case MATEMA:
		return TEMALookback(period)
	case MATRIMA:
		return TRIMALookback(period)
	case MAKAMA:
		return KAMALookback(period)
	case MAMAMA:
		return MAMALookback()
	default:
		return SMALookback(period)
	}
}
