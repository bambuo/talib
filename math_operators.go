package talib

import (
	"math"
)

// ADD 两个向量逐元素相加。
//
//	out[i] = inReal0[i] + inReal1[i]
func ADD(inReal0, inReal1 []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal0, "ADD:inReal0"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(inReal1, "ADD:inReal1"); err != nil {
		return nil, err
	}
	if len(inReal0) != len(inReal1) {
		return nil, ErrInputLengthMismatch("ADD: inReal0, inReal1")
	}

	n := len(inReal0)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = inReal0[i] + inReal1[i]
	}
	return out, nil
}

// DIV 两个向量逐元素相除。
//
//	out[i] = inReal0[i] / inReal1[i]
//
// 除以零时，结果为 NaN。
func DIV(inReal0, inReal1 []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal0, "DIV:inReal0"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(inReal1, "DIV:inReal1"); err != nil {
		return nil, err
	}
	if len(inReal0) != len(inReal1) {
		return nil, ErrInputLengthMismatch("DIV: inReal0, inReal1")
	}

	n := len(inReal0)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		if inReal1[i] == 0 {
			out[i] = math.NaN()
		} else {
			out[i] = inReal0[i] / inReal1[i]
		}
	}
	return out, nil
}

// MULT 两个向量逐元素相乘。
//
//	out[i] = inReal0[i] * inReal1[i]
func MULT(inReal0, inReal1 []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal0, "MULT:inReal0"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(inReal1, "MULT:inReal1"); err != nil {
		return nil, err
	}
	if len(inReal0) != len(inReal1) {
		return nil, ErrInputLengthMismatch("MULT: inReal0, inReal1")
	}

	n := len(inReal0)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = inReal0[i] * inReal1[i]
	}
	return out, nil
}

// SUB 两个向量逐元素相减。
//
//	out[i] = inReal0[i] - inReal1[i]
func SUB(inReal0, inReal1 []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal0, "SUB:inReal0"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(inReal1, "SUB:inReal1"); err != nil {
		return nil, err
	}
	if len(inReal0) != len(inReal1) {
		return nil, ErrInputLengthMismatch("SUB: inReal0, inReal1")
	}

	n := len(inReal0)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = inReal0[i] - inReal1[i]
	}
	return out, nil
}

// SUM 计算滑动窗口内的总和。
//
//	out[i] = sum(input[j]) 对于 j 在 [i-周期+1, i] 范围内
//
// 前 (周期-1) 个元素为 NaN。
func SUM(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "SUM:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "SUM:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	startIdx := period - 1
	// 初始窗口
	var sum float64
	for j := 0; j < period; j++ {
		sum += input[j]
	}
	out[startIdx] = sum

	// 滑动窗口
	for i := startIdx + 1; i < n; i++ {
		sum = sum - input[i-period] + input[i]
		out[i] = sum
	}

	return out, nil
}

// SUMLookback 返回 SUM 的 lookback。
func SUMLookback(period int) int {
	return period - 1
}

// MAX 计算滑动窗口内的最高值。
//
//	out[i] = max(input[j]) 对于 j 在 [i-周期+1, i] 范围内
//
// 前 (周期-1) 个元素为 NaN。
func MAX(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "MAX:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MAX:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		maxVal := input[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if input[j] > maxVal {
				maxVal = input[j]
			}
		}
		out[i] = maxVal
	}

	return out, nil
}

// MAXLookback 返回 MAX 的 lookback。
func MAXLookback(period int) int {
	return period - 1
}

// MAXINDEX 计算滑动窗口内最高值的索引（0-based，相对于输入起始）。
//
//	out[i] = 最高值在 input[j] 中的索引，j 在 [i-周期+1, i] 范围内
//
// 前 (周期-1) 个元素为 NaN。
func MAXINDEX(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "MAXINDEX:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MAXINDEX:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		maxVal := input[i-period+1]
		maxIdx := i - period + 1
		for j := i - period + 2; j <= i; j++ {
			if input[j] > maxVal {
				maxVal = input[j]
				maxIdx = j
			}
		}
		out[i] = float64(maxIdx)
	}

	return out, nil
}

// MAXINDEXLookback 返回 MAXINDEX 的 lookback。
func MAXINDEXLookback(period int) int {
	return period - 1
}

// MIN 计算滑动窗口内的最低值。
//
//	out[i] = min(input[j]) 对于 j 在 [i-周期+1, i] 范围内
//
// 前 (周期-1) 个元素为 NaN。
func MIN(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "MIN:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MIN:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		minVal := input[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if input[j] < minVal {
				minVal = input[j]
			}
		}
		out[i] = minVal
	}

	return out, nil
}

// MINLookback 返回 MIN 的 lookback。
func MINLookback(period int) int {
	return period - 1
}

// MININDEX 计算滑动窗口内最低值的索引（0-based，相对于输入起始）。
//
//	out[i] = 最低值在 input[j] 中的索引，j 在 [i-周期+1, i] 范围内
//
// 前 (周期-1) 个元素为 NaN。
func MININDEX(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "MININDEX:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MININDEX:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		minVal := input[i-period+1]
		minIdx := i - period + 1
		for j := i - period + 2; j <= i; j++ {
			if input[j] < minVal {
				minVal = input[j]
				minIdx = j
			}
		}
		out[i] = float64(minIdx)
	}

	return out, nil
}

// MININDEXLookback 返回 MININDEX 的 lookback。
func MININDEXLookback(period int) int {
	return period - 1
}

// MinMaxResult 包含滑动窗口的最小值和最大值结果。
type MinMaxResult struct {
	Min []float64 // 滑动窗口最小值
	Max []float64 // 滑动窗口最大值
}

// MINMAX 计算滑动窗口内的最低值和最高值。
//
// 前 (周期-1) 个元素为 NaN。
func MINMAX(input []float64, period int) (*MinMaxResult, error) {
	if err := ValidateNumericInput(input, "MINMAX:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MINMAX:period"); err != nil {
		return nil, err
	}

	n := len(input)
	minOut := MakeOutput(n)
	maxOut := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		minVal := input[i-period+1]
		maxVal := input[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if input[j] < minVal {
				minVal = input[j]
			}
			if input[j] > maxVal {
				maxVal = input[j]
			}
		}
		minOut[i] = minVal
		maxOut[i] = maxVal
	}

	return &MinMaxResult{
		Min: minOut,
		Max: maxOut,
	}, nil
}

// MINMAXLookback 返回 MINMAX 的 lookback。
func MINMAXLookback(period int) int {
	return period - 1
}

// MinMaxIndexResult 包含滑动窗口的 MinIndex 和 MaxIndex 结果。
type MinMaxIndexResult struct {
	MinIdx []float64 // 最小值索引（0-based）
	MaxIdx []float64 // 最大值索引（0-based）
}

// MINMAXINDEX 计算滑动窗口内最低值和最高值的索引。
//
// 前 (周期-1) 个元素为 NaN。
func MINMAXINDEX(input []float64, period int) (*MinMaxIndexResult, error) {
	if err := ValidateNumericInput(input, "MINMAXINDEX:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MINMAXINDEX:period"); err != nil {
		return nil, err
	}

	n := len(input)
	minIdxOut := MakeOutput(n)
	maxIdxOut := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		minVal := input[i-period+1]
		maxVal := input[i-period+1]
		minIdx := i - period + 1
		maxIdx := i - period + 1
		for j := i - period + 2; j <= i; j++ {
			if input[j] < minVal {
				minVal = input[j]
				minIdx = j
			}
			if input[j] > maxVal {
				maxVal = input[j]
				maxIdx = j
			}
		}
		minIdxOut[i] = float64(minIdx)
		maxIdxOut[i] = float64(maxIdx)
	}

	return &MinMaxIndexResult{
		MinIdx: minIdxOut,
		MaxIdx: maxIdxOut,
	}, nil
}

// MINMAXINDEXLookback 返回 MINMAXINDEX 的 lookback。
func MINMAXINDEXLookback(period int) int {
	return period - 1
}
