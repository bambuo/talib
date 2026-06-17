package talib

import (
	"math"
	"testing"
)

func TestMACD_Basic(t *testing.T) {
	input := make([]float64, 40)
	for i := range input {
		input[i] = float64(i + 1)
	}

	result, err := MACD(input, 12, 26, 9)
	if err != nil {
		t.Fatalf("MACD returned error: %v", err)
	}

	if len(result.MACD) != len(input) {
		t.Errorf("expected MACD length %d, got %d", len(input), len(result.MACD))
	}
	if len(result.Signal) != len(input) {
		t.Errorf("expected Signal length %d, got %d", len(input), len(result.Signal))
	}
	if len(result.Histogram) != len(input) {
		t.Errorf("expected Histogram length %d, got %d", len(input), len(result.Histogram))
	}

	// 验证前导 NaN
	startIdx := MACDLookback(12, 26, 9)
	for i := 0; i < startIdx; i++ {
		if !math.IsNaN(result.Histogram[i]) {
			t.Errorf("expected NaN Histogram at index %d, got %f", i, result.Histogram[i])
		}
	}
}

func TestMACD_InvalidFastEqualsSlow(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	_, err := MACD(input, 12, 12, 9)
	if err == nil {
		t.Error("expected error when fast == slow")
	}
}

func TestMACD_InvalidFastGreaterThanSlow(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	_, err := MACD(input, 26, 12, 9)
	if err == nil {
		t.Error("expected error when fast > slow")
	}
}

func TestMACD_OutputRelation(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
		16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
		31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

	result, err := MACD(input, 12, 26, 9)
	if err != nil {
		t.Fatalf("MACD returned error: %v", err)
	}

	// 直方图应对所有有效索引等于 MACD - Signal
	for i := 0; i < len(input); i++ {
		if math.IsNaN(result.Histogram[i]) {
			continue
		}
		expected := result.MACD[i] - result.Signal[i]
		if math.Abs(result.Histogram[i]-expected) > 1e-10 {
			t.Errorf("Histogram[%d]: expected %f (MACD-Signal), got %f",
				i, expected, result.Histogram[i])
		}
	}
}
