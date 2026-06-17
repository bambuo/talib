package talib

import (
	"math"
	"testing"
)

func TestATR_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13}
	close := []float64{9, 11, 10, 12, 13, 14, 13, 12, 15, 14}
	period := 3

	result, err := ATR(high, low, close, period)
	if err != nil {
		t.Fatalf("ATR returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("expected length %d, got %d", len(close), len(result))
	}

	// ATR 在索引 period 处的首个值应为前 3 个 TR 值的 SMA
	// TR[0] = 10-8 = 2
	// TR[1] = max(12-9=3, |12-9|=3, |9-9|=0) = 3
	// TR[2] = max(11-10=1, |11-11|=0, |10-11|=1) = 1
	expectedFirstATR := (2.0 + 3.0 + 1.0) / 3.0
	if math.Abs(result[period]-expectedFirstATR) > 1e-10 {
		t.Errorf("first ATR: expected %f, got %f", expectedFirstATR, result[period])
	}
}

func TestTRANGE_Basic(t *testing.T) {
	high := []float64{10, 12, 11}
	low := []float64{8, 9, 10}
	close := []float64{9, 11, 10}

	result, err := TRANGE(high, low, close)
	if err != nil {
		t.Fatalf("TRANGE returned error: %v", err)
	}

	// TR[0] = 10 - 8 = 2
	if math.Abs(result[0]-2.0) > 1e-10 {
		t.Errorf("TR[0]: expected 2.0, got %f", result[0])
	}
	// TR[1] = max(12-9=3, |12-9|=3, |9-9|=0) = 3
	if math.Abs(result[1]-3.0) > 1e-10 {
		t.Errorf("TR[1]: expected 3.0, got %f", result[1])
	}
}

func TestOBV_Basic(t *testing.T) {
	close := []float64{10, 12, 11, 13, 12, 14}
	volume := []float64{100, 200, 150, 300, 250, 400}

	result, err := OBV(close, volume)
	if err != nil {
		t.Fatalf("OBV returned error: %v", err)
	}

	// OBV[0] = 0
	if result[0] != 0 {
		t.Errorf("OBV[0]: expected 0, got %f", result[0])
	}
	// close[1] > close[0], so OBV[1] = 0 + 200 = 200
	if result[1] != 200 {
		t.Errorf("OBV[1]: expected 200, got %f", result[1])
	}
	// close[2] < close[1], so OBV[2] = 200 - 150 = 50
	if result[2] != 50 {
		t.Errorf("OBV[2]: expected 50, got %f", result[2])
	}
	// close[3] > close[2], so OBV[3] = 50 + 300 = 350
	if result[3] != 350 {
		t.Errorf("OBV[3]: expected 350, got %f", result[3])
	}
}

func TestCCI_Basic(t *testing.T) {
	high := []float64{50, 52, 51, 53, 54, 55, 54, 53, 56, 55}
	low := []float64{48, 50, 49, 51, 52, 53, 52, 51, 54, 53}
	close := []float64{49, 51, 50, 52, 53, 54, 53, 52, 55, 54}
	period := 5

	result, err := CCI(high, low, close, period)
	if err != nil {
		t.Fatalf("CCI returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("expected length %d, got %d", len(close), len(result))
	}
}

func TestWILLIAMS_R_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13}
	close := []float64{9, 11, 10, 12, 13, 14, 13, 12, 15, 14}
	period := 5

	result, err := WILLIAMS_R(high, low, close, period)
	if err != nil {
		t.Fatalf("WILLIAMS_R returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("expected length %d, got %d", len(close), len(result))
	}

	// %R 范围检查：应在 [-100, 0] 范围内
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < -100.0 || v > 0.0 {
			t.Errorf("WILLIAMS_R[%d] = %f, out of [-100, 0] range", i, v)
		}
	}
}
