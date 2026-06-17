package talib

import (
	"math"
	"testing"
)

func TestRSI_Basic(t *testing.T) {
	// 简单上升趋势：RSI 应较高
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	period := 14

	result, err := RSI(input, period)
	if err != nil {
		t.Fatalf("RSI returned error: %v", err)
	}

	// 纯上升趋势的 RSI 应为 100
	lastRSI := result[len(result)-1]
	if math.Abs(lastRSI-100.0) > 1e-10 {
		t.Errorf("RSI of pure uptrend: expected 100, got %f", lastRSI)
	}
}

func TestRSI_PureDowntrend(t *testing.T) {
	input := []float64{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	period := 14

	result, err := RSI(input, period)
	if err != nil {
		t.Fatalf("RSI returned error: %v", err)
	}

	// 纯下降趋势的 RSI 应为 0
	lastRSI := result[len(result)-1]
	if math.Abs(lastRSI-0.0) > 1e-10 {
		t.Errorf("RSI of pure downtrend: expected 0, got %f", lastRSI)
	}
}

func TestRSI_OutputLength(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	period := 14

	result, err := RSI(input, period)
	if err != nil {
		t.Fatalf("RSI returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}

	// 前 'period' 个值应为 NaN
	for i := 0; i < period; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("expected NaN at index %d, got %f", i, result[i])
		}
	}
}

func TestRSI_InvalidInput(t *testing.T) {
	_, err := RSI(nil, 5)
	if err == nil {
		t.Error("expected error for nil input")
	}
	_, err = RSI([]float64{1, 2}, 10)
	if err == nil {
		t.Error("expected error when period > input length")
	}
}

func TestRSI_Range(t *testing.T) {
	// 随机数据：RSI 应始终在 [0, 100] 范围内
	input := []float64{44.34, 44.09, 44.15, 43.61, 44.33, 44.83, 45.10, 45.42, 45.84,
		46.08, 45.89, 46.03, 45.61, 46.28, 46.28, 46.00, 46.03, 46.41,
		46.22, 45.64, 46.21, 46.25, 45.71, 46.45, 45.78, 45.35, 44.03,
		44.18, 44.22, 44.57, 43.42, 42.66, 43.13}

	result, err := RSI(input, 14)
	if err != nil {
		t.Fatalf("RSI returned error: %v", err)
	}

	for i, rsi := range result {
		if math.IsNaN(rsi) {
			continue
		}
		if rsi < 0 || rsi > 100 {
			t.Errorf("RSI[%d] = %f, out of [0,100] range", i, rsi)
		}
	}
}

func TestRSILookback(t *testing.T) {
	lb := RSILookback(14)
	if lb != 14 {
		t.Errorf("RSI lookback: expected 14, got %d", lb)
	}
}
