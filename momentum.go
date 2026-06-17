package talib

import (
	"math"
)

// BOP 计算平衡力量指标。
//
//	BOP = (收盘价 - 开盘价) / (最高价 - 最低价)
//
// BOP 衡量每个周期买入与卖出压力的强度。
// 当最高 == 最低时返回 NaN（无区间）。
func BOP(open, high, low, close []float64) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(open, "BOP:open"); err != nil {
		return nil, err
	}
	if len(open) != len(close) {
		return nil, ErrInputLengthMismatch("BOP: open, close")
	}

	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		diff := high[i] - low[i]
		if diff == 0 {
			out[i] = 0
		} else {
			out[i] = (close[i] - open[i]) / diff
		}
	}
	return out, nil
}

// CMO 计算 Chande 动能振荡器。
//
//	涨幅和 = sum(max(价格[i] - 价格[i-1], 0)) 在周期内
//	跌幅和 = sum(max(价格[i-1] - 价格[i], 0)) 在周期内
//	CMO   = (涨幅和 - 跌幅和) / (涨幅和 + 跌幅和) * 100
//
// 前 (周期) 个元素为 NaN。
func CMO(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "CMO:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "CMO:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	if n <= period {
		return out, nil
	}

	// 计算第一个窗口的涨跌幅和
	var upSum, downSum float64
	for i := 1; i <= period; i++ {
		diff := input[i] - input[i-1]
		if diff > 0 {
			upSum += diff
		} else {
			downSum += -diff
		}
	}

	idx := period
	total := upSum + downSum
	if total == 0 {
		out[idx] = 0
	} else {
		out[idx] = (upSum - downSum) / total * 100.0
	}

	// 滑动窗口
	for i := period + 1; i < n; i++ {
		// 添加新的差值
		diff := input[i] - input[i-1]
		if diff > 0 {
			upSum += diff
		} else {
			downSum += -diff
		}

		// 移除旧的差值
		oldDiff := input[i-period] - input[i-period-1]
		if oldDiff > 0 {
			upSum -= oldDiff
		} else {
			downSum -= -oldDiff
		}

		total := upSum + downSum
		if total == 0 {
			out[i] = 0
		} else {
			out[i] = (upSum - downSum) / total * 100.0
		}
	}

	return out, nil
}

// CMOLookback 返回 CMO 的 lookback。
func CMOLookback(period int) int {
	return period
}

// DX 计算方向运动指标（原始值，在 ADX 平滑之前）。
//
//	DX = |+DI - -DI| / (+DI + -DI) * 100
//
// 这是 ADX 平滑的中间值。用于
// 识别每个点的方向运动强度。
func DX(high, low, close []float64, period int) ([]float64, error) {
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
	out := MakeOutput(n)

	for i := 0; i < n; i++ {
		if IsNaN(plusDI[i]) || IsNaN(minusDI[i]) {
			continue
		}
		sum := plusDI[i] + minusDI[i]
		if sum == 0 {
			out[i] = 0
		} else {
			out[i] = math.Abs(plusDI[i]-minusDI[i]) / sum * 100.0
		}
	}

	return out, nil
}

// DXLookback 返回 DX 的 lookback。
func DXLookback(period int) int {
	return period
}

// MOM 计算动量。
//
//	MOM[i] = 价格[i] - 价格[i-周期]
//
// 前 (周期) 个元素为 NaN。
func MOM(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "MOM:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MOM:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period; i < n; i++ {
		out[i] = input[i] - input[i-period]
	}

	return out, nil
}

// MOMLookback 返回 MOM 的 lookback。
func MOMLookback(period int) int {
	return period
}

// ROC 计算变动率。
//
//	ROC[i] = (价格[i] / 价格[i-周期] - 1) * 100
//
// 前 (周期) 个元素为 NaN。
func ROC(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "ROC:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "ROC:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period; i < n; i++ {
		if input[i-period] != 0 {
			out[i] = (input[i]/input[i-period] - 1.0) * 100.0
		}
	}

	return out, nil
}

// ROCLookback 返回 ROC 的 lookback。
func ROCLookback(period int) int {
	return period
}

// ROCP 计算变动率百分比。
//
//	ROCP[i] = (价格[i] - 价格[i-周期]) / 价格[i-周期]
//
// 前 (周期) 个元素为 NaN。
func ROCP(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "ROCP:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "ROCP:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period; i < n; i++ {
		if input[i-period] != 0 {
			out[i] = (input[i] - input[i-period]) / input[i-period]
		}
	}

	return out, nil
}

// ROCPLookback 返回 ROCP 的 lookback。
func ROCPLookback(period int) int {
	return period
}

// ROCR 计算变动率比率。
//
//	ROCR[i] = 价格[i] / 价格[i-周期]
//
// 前 (周期) 个元素为 NaN。
func ROCR(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "ROCR:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "ROCR:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period; i < n; i++ {
		if input[i-period] != 0 {
			out[i] = input[i] / input[i-period]
		}
	}

	return out, nil
}

// ROCRLookback 返回 ROCR 的 lookback。
func ROCRLookback(period int) int {
	return period
}

// ROCR100 计算变动率比率乘以 100。
//
//	ROCR100[i] = (价格[i] / 价格[i-周期]) * 100
//
// 前 (周期) 个元素为 NaN。
func ROCR100(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "ROCR100:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "ROCR100:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period; i < n; i++ {
		if input[i-period] != 0 {
			out[i] = input[i] / input[i-period] * 100.0
		}
	}

	return out, nil
}

// ROCR100Lookback 返回 ROCR100 的 lookback。
func ROCR100Lookback(period int) int {
	return period
}
