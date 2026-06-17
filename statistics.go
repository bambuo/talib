package talib

import (
	"math"
)

// STDDEV 计算滑动窗口内的总体标准差。
//
//	σ[i] = sqrt( Σ(input[j] - μ)² / period )
//
// 其中 μ = 在点 i 处的 mean(input, period)，且 j ∈ [i-period+1, i]。
// 前 (period-1) 个元素为 NaN。
func STDDEV(input []float64, period int, nbDev float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "STDDEV:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "STDDEV:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		// 计算窗口均值
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += input[j]
		}
		mean := sum / float64(period)

		// Compute variance
		var sumSq float64
		for j := i - period + 1; j <= i; j++ {
			diff := input[j] - mean
			sumSq += diff * diff
		}

		out[i] = math.Sqrt(sumSq/float64(period)) * nbDev
	}

	return out, nil
}

// STDDEVLookback 返回 STDDEV 的 lookback。
func STDDEVLookback(period int) int {
	return period - 1
}

// VAR 计算滑动窗口内的总体方差。
//
//	σ²[i] = Σ(input[j] - μ)² / period
//
// 其中 μ = 在点 i 处的 mean(input, period)。
// 前 (period-1) 个元素为 NaN。
func VAR(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "VAR:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "VAR:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		// 计算均值
		var sum float64
		for j := i - period + 1; j <= i; j++ {
			sum += input[j]
		}
		mean := sum / float64(period)

		// 计算方差
		var sumSq float64
		for j := i - period + 1; j <= i; j++ {
			diff := input[j] - mean
			sumSq += diff * diff
		}

		out[i] = sumSq / float64(period)
	}

	return out, nil
}

// VARLookback 返回 VAR 的 lookback。
func VARLookback(period int) int {
	return period - 1
}

// CORREL 计算滑动窗口内两个序列的皮尔逊相关系数。
//
//	r[i] = Cov(X, Y) / (σx * σy)
//
// 前 (period-1) 个元素为 NaN。
func CORREL(inReal0, inReal1 []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(inReal0, "CORREL:inReal0"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(inReal1, "CORREL:inReal1"); err != nil {
		return nil, err
	}
	if len(inReal0) != len(inReal1) {
		return nil, ErrInputLengthMismatch("CORREL: inReal0, inReal1")
	}
	if err := ValidatePeriod(period, len(inReal0), "CORREL:period"); err != nil {
		return nil, err
	}

	n := len(inReal0)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		// Compute means
		var sumX, sumY float64
		for j := i - period + 1; j <= i; j++ {
			sumX += inReal0[j]
			sumY += inReal1[j]
		}
		meanX := sumX / float64(period)
		meanY := sumY / float64(period)

		// Compute covariance and variances
		var cov, varX, varY float64
		for j := i - period + 1; j <= i; j++ {
			dx := inReal0[j] - meanX
			dy := inReal1[j] - meanY
			cov += dx * dy
			varX += dx * dx
			varY += dy * dy
		}

		if varX == 0 || varY == 0 {
			out[i] = 0
		} else {
			out[i] = cov / math.Sqrt(varX*varY)
		}
	}

	return out, nil
}

// CORRELLookback 返回 CORREL 的 lookback。
func CORRELLookback(period int) int {
	return period - 1
}

// BETA 计算滑动窗口内两个序列的贝塔系数。
//
//	β[i] = Cov(X, Y) / Var(Y)
//
// 衡量 X 相对于 Y 的敏感度。
// 前 (period-1) 个元素为 NaN。
func BETA(inReal0, inReal1 []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(inReal0, "BETA:inReal0"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(inReal1, "BETA:inReal1"); err != nil {
		return nil, err
	}
	if len(inReal0) != len(inReal1) {
		return nil, ErrInputLengthMismatch("BETA: inReal0, inReal1")
	}
	if err := ValidatePeriod(period, len(inReal0), "BETA:period"); err != nil {
		return nil, err
	}

	n := len(inReal0)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		// Compute means
		var sumX, sumY float64
		for j := i - period + 1; j <= i; j++ {
			sumX += inReal0[j]
			sumY += inReal1[j]
		}
		meanX := sumX / float64(period)
		meanY := sumY / float64(period)

		// Compute covariance and variance of Y
		var cov, varY float64
		for j := i - period + 1; j <= i; j++ {
			cov += (inReal0[j] - meanX) * (inReal1[j] - meanY)
			dy := inReal1[j] - meanY
			varY += dy * dy
		}

		if varY == 0 {
			out[i] = 0
		} else {
			out[i] = cov / varY
		}
	}

	return out, nil
}

// BETALookback 返回 BETA 的 lookback。
func BETALookback(period int) int {
	return period - 1
}

// linearReg 从 (x, y) 对的切片计算 OLS 斜率和截距。
// x 值假定为 0, 1, ..., N-1。
func linearReg(data []float64) (slope, intercept float64) {
	n := float64(len(data))

	var sumX, sumY, sumXY, sumX2 float64
	for i, y := range data {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	denom := n*sumX2 - sumX*sumX
	if denom == 0 {
		return 0, sumY / n
	}

	slope = (n*sumXY - sumX*sumY) / denom
	intercept = (sumY - slope*sumX) / n
	return
}

// LINEARREG 计算每个滑动窗口最后一点的线性回归预测值。
//
// 对于每个长度为 'period' 的窗口，拟合 OLS 回归线 y = a + b*x
// 其中 x = 0, 1, ..., period-1。返回 x = period-1 处的预测 y
// （窗口中的最后一点）。
//
// 前 (period-1) 个元素为 NaN。
func LINEARREG(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "LINEARREG:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "LINEARREG:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		window := input[i-period+1 : i+1]
		slope, intercept := linearReg(window)
		out[i] = intercept + slope*float64(period-1)
	}

	return out, nil
}

// LINEARREGLookback 返回 LINEARREG 的 lookback。
func LINEARREGLookback(period int) int {
	return period - 1
}

// LINEARREG_ANGLE 计算滑动窗口内线性回归线的角度（度数）。
//
//	角度 = arctan(斜率) * 180 / π
//
// 前 (period-1) 个元素为 NaN。
func LINEARREG_ANGLE(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "LINEARREG_ANGLE:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "LINEARREG_ANGLE:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		window := input[i-period+1 : i+1]
		slope, _ := linearReg(window)
		out[i] = math.Atan(slope) * 180.0 / math.Pi
	}

	return out, nil
}

// LINEARREG_ANGLELookback 返回 LINEARREG_ANGLE 的 lookback。
func LINEARREG_ANGLELookback(period int) int {
	return period - 1
}

// LINEARREG_INTERCEPT 计算滑动窗口内 OLS 回归线的 y 截距。
//
// 前 (period-1) 个元素为 NaN。
func LINEARREG_INTERCEPT(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "LINEARREG_INTERCEPT:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "LINEARREG_INTERCEPT:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		window := input[i-period+1 : i+1]
		_, intercept := linearReg(window)
		out[i] = intercept
	}

	return out, nil
}

// LINEARREG_INTERCEPTLookback 返回 LINEARREG_INTERCEPT 的 lookback。
func LINEARREG_INTERCEPTLookback(period int) int {
	return period - 1
}

// LINEARREG_SLOPE 计算滑动窗口内 OLS 回归线的斜率。
//
// 前 (period-1) 个元素为 NaN。
func LINEARREG_SLOPE(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "LINEARREG_SLOPE:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "LINEARREG_SLOPE:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		window := input[i-period+1 : i+1]
		slope, _ := linearReg(window)
		out[i] = slope
	}

	return out, nil
}

// LINEARREG_SLOPELookback 返回 LINEARREG_SLOPE 的 lookback。
func LINEARREG_SLOPELookback(period int) int {
	return period - 1
}

// TSF 计算时间序列预测：下一周期（向前一步）的线性回归预测值。
//
// 对于每个窗口，预测 x = period 处的 y（窗口之后的一个点），
// 即 TSF = 截距 + 斜率 * period。
//
// 前 (period-1) 个元素为 NaN。
func TSF(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "TSF:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "TSF:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		window := input[i-period+1 : i+1]
		slope, intercept := linearReg(window)
		out[i] = intercept + slope*float64(period)
	}

	return out, nil
}

// TSFLookback 返回 TSF 的 lookback。
func TSFLookback(period int) int {
	return period - 1
}
