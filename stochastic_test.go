package talib

import (
	"math"
	"testing"
)

func TestSTOCH_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15, 17, 18, 16, 15, 19}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13, 14, 15, 13, 12, 16}
	close := []float64{9, 11, 10, 12, 13, 14, 13, 12, 15, 14, 16, 17, 14, 13, 18}

	result, err := STOCH(high, low, close, 5, 3, 3, MASMA)
	if err != nil {
		t.Fatalf("STOCH returned error: %v", err)
	}

	if len(result.K) != len(close) || len(result.D) != len(close) {
		t.Errorf("output length mismatch: K=%d, D=%d, expected=%d",
			len(result.K), len(result.D), len(close))
	}

	// %K 和 %D 应在 [0, 100] 范围内
	for i, v := range result.K {
		if math.IsNaN(v) {
			continue
		}
		if v < -1e-10 || v > 100.0+1e-10 {
			t.Errorf("STOCH %%K[%d] = %f, out of range", i, v)
		}
	}
}

func TestSTOCHF_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15, 17, 18}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13, 14, 15}
	close := []float64{9, 11, 10, 12, 13, 14, 13, 12, 15, 14, 16, 17}

	result, err := STOCHF(high, low, close, 5, 3, MASMA)
	if err != nil {
		t.Fatalf("STOCHF returned error: %v", err)
	}

	if len(result.K) != len(close) {
		t.Errorf("K length mismatch: %d vs %d", len(result.K), len(close))
	}
}

func TestADX_Basic(t *testing.T) {
	high := []float64{50, 52, 51, 53, 54, 55, 54, 53, 56, 55, 57, 58, 56, 55, 59,
		60, 58, 57, 61, 60, 62, 63, 61, 60, 64}
	low := []float64{48, 50, 49, 51, 52, 53, 52, 51, 54, 53, 54, 55, 53, 52, 56,
		57, 55, 54, 58, 57, 59, 60, 58, 57, 61}
	close := []float64{49, 51, 50, 52, 53, 54, 53, 52, 55, 54, 56, 57, 54, 53, 58,
		59, 56, 55, 60, 59, 61, 62, 59, 58, 63}
	period := 14

	result, err := ADX(high, low, close, period)
	if err != nil {
		t.Fatalf("ADX returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("length mismatch: %d vs %d", len(result), len(close))
	}

	// ADX 应在 [0, 100] 范围内
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < -1e-10 || v > 100.0+1e-10 {
			t.Errorf("ADX[%d] = %f, out of range", i, v)
		}
	}
}

func TestMFI_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15, 17, 18, 16, 15, 19}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13, 14, 15, 13, 12, 16}
	close := []float64{9, 11, 10, 12, 13, 14, 13, 12, 15, 14, 16, 17, 14, 13, 18}
	volume := []float64{100, 200, 150, 300, 250, 400, 350, 200, 500, 300, 450, 600, 350, 250, 550}
	period := 5

	result, err := MFI(high, low, close, volume, period)
	if err != nil {
		t.Fatalf("MFI returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("length mismatch: %d vs %d", len(result), len(close))
	}

	// MFI 应在 [0, 100] 范围内
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < -1e-10 || v > 100.0+1e-10 {
			t.Errorf("MFI[%d] = %f, out of range", i, v)
		}
	}
}

func TestAROON_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15, 17, 18, 16, 15, 19}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13, 14, 15, 13, 12, 16}
	period := 5

	result, err := AROON(high, low, period)
	if err != nil {
		t.Fatalf("AROON returned error: %v", err)
	}

	// Aroon Up/Down 应在 [0, 100] 范围内
	for i, v := range result.Up {
		if math.IsNaN(v) {
			continue
		}
		if v < -1e-10 || v > 100.0+1e-10 {
			t.Errorf("AroonUp[%d] = %f, out of range", i, v)
		}
	}

	for i, v := range result.Down {
		if math.IsNaN(v) {
			continue
		}
		if v < -1e-10 || v > 100.0+1e-10 {
			t.Errorf("AroonDown[%d] = %f, out of range", i, v)
		}
	}
}

func TestAROONOSC_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15, 17, 18}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13, 14, 15}
	period := 4

	result, err := AROONOSC(high, low, period)
	if err != nil {
		t.Fatalf("AROONOSC returned error: %v", err)
	}

	// AroonOSC 应在 [-100, 100] 范围内
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < -100.0-1e-10 || v > 100.0+1e-10 {
			t.Errorf("AroonOSC[%d] = %f, out of range", i, v)
		}
	}
}

func TestULTOSC_Basic(t *testing.T) {
	high := []float64{50, 52, 51, 53, 54, 55, 54, 53, 56, 55, 57, 58, 56, 55, 59,
		60, 58, 57, 61, 60, 62, 63, 61, 60, 64, 65, 63, 62, 66, 65, 67, 68}
	low := []float64{48, 50, 49, 51, 52, 53, 52, 51, 54, 53, 54, 55, 53, 52, 56,
		57, 55, 54, 58, 57, 59, 60, 58, 57, 61, 62, 60, 59, 63, 62, 64, 65}
	close := []float64{49, 51, 50, 52, 53, 54, 53, 52, 55, 54, 56, 57, 54, 53, 58,
		59, 56, 55, 60, 59, 61, 62, 59, 58, 63, 64, 61, 60, 65, 64, 66, 67}
	period1, period2, period3 := 7, 14, 28

	result, err := ULTOSC(high, low, close, period1, period2, period3)
	if err != nil {
		t.Fatalf("ULTOSC returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("length mismatch: %d vs %d", len(result), len(close))
	}
}
