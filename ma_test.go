package talib

import (
	"math"
	"testing"
)

func TestSMA(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	period := 3

	result, err := SMA(input, period)
	if err != nil {
		t.Fatalf("SMA returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}

	// 前导 NaN
	for i := 0; i < period-1; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("expected NaN at index %d, got %f", i, result[i])
		}
	}

	// 期望值
	expected := []float64{math.NaN(), math.NaN(), 2, 3, 4, 5, 6, 7, 8, 9}
	for i := period - 1; i < len(input); i++ {
		if math.Abs(result[i]-expected[i]) > 1e-10 {
			t.Errorf("SMA[%d]: expected %f, got %f", i, expected[i], result[i])
		}
	}
}

func TestSMA_SinglePeriod(t *testing.T) {
	input := []float64{5, 10, 15}
	result, err := SMA(input, 1)
	if err != nil {
		t.Fatalf("SMA returned error: %v", err)
	}
	for i, v := range input {
		if math.Abs(result[i]-v) > 1e-10 {
			t.Errorf("SMA period=1: expected %f at %d, got %f", v, i, result[i])
		}
	}
}

func TestSMA_FullLengthPeriod(t *testing.T) {
	input := []float64{1, 2, 3}
	result, err := SMA(input, 3)
	if err != nil {
		t.Fatalf("SMA returned error: %v", err)
	}
	// 只有最后一个元素有效
	expected := 2.0 // (1+2+3)/3
	if math.Abs(result[2]-expected) > 1e-10 {
		t.Errorf("expected %f, got %f", expected, result[2])
	}
}

func TestSMA_InvalidInput(t *testing.T) {
	_, err := SMA(nil, 5)
	if err == nil {
		t.Error("expected error for nil input")
	}
	_, err = SMA([]float64{}, 5)
	if err == nil {
		t.Error("expected error for empty input")
	}
	_, err = SMA([]float64{1, 2}, 5)
	if err == nil {
		t.Error("expected error when period > input length")
	}
}

func TestEMA(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	period := 3

	result, err := EMA(input, period)
	if err != nil {
		t.Fatalf("EMA returned error: %v", err)
	}

	// 第一个有效值 = [1,2,3] 的 SMA = 2.0
	if math.Abs(result[2]-2.0) > 1e-10 {
		t.Errorf("EMA[2] (first value): expected 2.0, got %f", result[2])
	}

	// 验证单调增长
	for i := 3; i < len(input); i++ {
		if result[i] < result[i-1] {
			t.Errorf("EMA should be monotonic increasing for increasing input, index %d: %f → %f",
				i, result[i-1], result[i])
		}
	}
}

func TestEMALookback(t *testing.T) {
	lb := EMALookback(10)
	if lb != 9 {
		t.Errorf("EMA lookback for period 10: expected 9, got %d", lb)
	}
}

func TestWMA(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5}
	period := 3

	result, err := WMA(input, period)
	if err != nil {
		t.Fatalf("WMA returned error: %v", err)
	}

	// WMA[2]: (1*1 + 2*2 + 3*3) / (1+2+3) = (1+4+9)/6 = 14/6 = 2.3333...
	expected := 14.0 / 6.0
	if math.Abs(result[2]-expected) > 1e-10 {
		t.Errorf("WMA[2]: expected %f, got %f", expected, result[2])
	}
}

func TestTRIMA(t *testing.T) {
	// TRIMA[4] with period 4:
	// n = (4+1)/2 = 2 (floor)
	// Period is even, so: SMA(SMA(input, 2), 3)
	// SMA(input, 2) = [NaN, 1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5]
	// SMA(sma1, 3) = [NaN, NaN, NaN, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5]
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result, err := TRIMA(input, 4)
	if err != nil {
		t.Fatalf("TRIMA returned error: %v", err)
	}

	if !math.IsNaN(result[0]) || !math.IsNaN(result[1]) || !math.IsNaN(result[2]) {
		t.Error("TRIMA: first 3 values should be NaN for period 4")
	}
}

func TestMA_Dispatcher(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5}
	period := 3

	// 验证每种 MA 类型通过调度器正常工作
	for _, mt := range []MAType{MASMA, MAEMA, MAWMA} {
		result, err := MA(input, period, mt)
		if err != nil {
			t.Errorf("MA(%s) returned error: %v", mt, err)
		}
		if len(result) != len(input) {
			t.Errorf("MA(%s): wrong output length", mt)
		}
	}
}
