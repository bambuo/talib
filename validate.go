package talib

import (
	"errors"
	"math"
)

// ErrFastPeriodGreaterOrEqualSlow 快速周期 >= 慢速周期时返回的错误。
var ErrFastPeriodGreaterOrEqualSlow = errors.New("MACD: fast period must be less than slow period")

// ErrInputLengthMismatch 返回一个格式化的错误，当两个切片长度不一致时使用。
func ErrInputLengthMismatch(labels string) error {
	return errors.New(labels + " must have equal length")
}

// ValidateNumericInput 检查输入切片非 nil 且非空。
func ValidateNumericInput(input []float64, label string) error {
	if input == nil {
		return errors.New(label + ": input is nil")
	}
	if len(input) == 0 {
		return errors.New(label + ": input is empty")
	}
	return nil
}

// ValidateOHLCInput 检查 OHLC 切片有效且长度相等。
func ValidateOHLCInput(high, low, close []float64) error {
	if err := ValidateNumericInput(high, "high"); err != nil {
		return err
	}
	if err := ValidateNumericInput(low, "low"); err != nil {
		return err
	}
	if err := ValidateNumericInput(close, "close"); err != nil {
		return err
	}
	if len(high) != len(low) || len(high) != len(close) {
		return errors.New("high, low, close must have equal length")
	}
	return nil
}

// ValidatePeriod 检查周期为正数且不超过输入长度。
func ValidatePeriod(period, inputLen int, label string) error {
	if period <= 0 {
		return errors.New(label + ": period must be positive")
	}
	if period > inputLen {
		return errors.New(label + ": period exceeds input length")
	}
	return nil
}

// IsNaN 如果值为 NaN 则返回 true。
func IsNaN(v float64) bool {
	return math.IsNaN(v)
}

// MakeOutput 创建指定长度的输出切片，全部初始化为 NaN。
// 这匹配 TA-Lib 的约定，前导 lookback 值用 NaN 填充。
func MakeOutput(length int) []float64 {
	out := make([]float64, length)
	for i := range out {
		out[i] = math.NaN()
	}
	return out
}
