package talib

import (
	"math"
)

// 注意：这些函数对 []float64 向量逐元素应用数学运算。
// 它们匹配 TA-Lib 的数学变换类别。

// ACOS 计算每个元素的反余弦。
//
//	out[i] = acos(inReal[i])
func ACOS(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "ACOS:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Acos(inReal[i])
	}
	return out, nil
}

// ASIN 计算每个元素的反正弦。
//
//	out[i] = asin(inReal[i])
func ASIN(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "ASIN:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Asin(inReal[i])
	}
	return out, nil
}

// ATAN 计算每个元素的反正切。
//
//	out[i] = atan(inReal[i])
func ATAN(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "ATAN:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Atan(inReal[i])
	}
	return out, nil
}

// CEIL 计算每个元素的上取整（大于等于值的最小整数）。
//
//	out[i] = ceil(inReal[i])
func CEIL(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "CEIL:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Ceil(inReal[i])
	}
	return out, nil
}

// COS 计算每个元素的余弦。
//
//	out[i] = cos(inReal[i])
func COS(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "COS:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Cos(inReal[i])
	}
	return out, nil
}

// COSH 计算每个元素的双曲余弦。
//
//	out[i] = cosh(inReal[i])
func COSH(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "COSH:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Cosh(inReal[i])
	}
	return out, nil
}

// EXP 计算 e 的每个元素次幂：e^inReal[i]。
//
//	out[i] = exp(inReal[i])
func EXP(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "EXP:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Exp(inReal[i])
	}
	return out, nil
}

// FLOOR 计算每个元素的下取整（小于等于值的最大整数）。
//
//	out[i] = floor(inReal[i])
func FLOOR(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "FLOOR:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Floor(inReal[i])
	}
	return out, nil
}

// LN 计算每个元素的自然对数。
//
//	out[i] = ln(inReal[i])
func LN(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "LN:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Log(inReal[i])
	}
	return out, nil
}

// LOG10 计算每个元素的以 10 为底的对数。
//
//	out[i] = log10(inReal[i])
func LOG10(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "LOG10:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Log10(inReal[i])
	}
	return out, nil
}

// SIN 计算每个元素的正弦。
//
//	out[i] = sin(inReal[i])
func SIN(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "SIN:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Sin(inReal[i])
	}
	return out, nil
}

// SINH 计算每个元素的双曲正弦。
//
//	out[i] = sinh(inReal[i])
func SINH(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "SINH:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Sinh(inReal[i])
	}
	return out, nil
}

// SQRT 计算每个元素的平方根。
//
//	out[i] = sqrt(inReal[i])
func SQRT(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "SQRT:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Sqrt(inReal[i])
	}
	return out, nil
}

// TAN 计算每个元素的正切。
//
//	out[i] = tan(inReal[i])
func TAN(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "TAN:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Tan(inReal[i])
	}
	return out, nil
}

// TANH 计算每个元素的双曲正切。
//
//	out[i] = tanh(inReal[i])
func TANH(inReal []float64) ([]float64, error) {
	if err := ValidateNumericInput(inReal, "TANH:inReal"); err != nil {
		return nil, err
	}

	n := len(inReal)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = math.Tanh(inReal[i])
	}
	return out, nil
}
