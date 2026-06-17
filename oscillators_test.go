package talib

import (
	"math"
	"testing"
)

func TestAPO_Basic(t *testing.T) {
	input := make([]float64, 40)
	for i := range input {
		input[i] = float64(i + 1)
	}

	result, err := APO(input, 12, 26, MAEMA)
	if err != nil {
		t.Fatalf("APO returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("output length mismatch: %d vs %d", len(result), len(input))
	}
}

func TestPPO_Basic(t *testing.T) {
	input := make([]float64, 50)
	for i := range input {
		input[i] = 100.0 + float64(i)*0.5
	}

	result, err := PPO(input, 12, 26, MAEMA)
	if err != nil {
		t.Fatalf("PPO returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}

func TestMACDEXT_Basic(t *testing.T) {
	input := make([]float64, 50)
	for i := range input {
		input[i] = float64(i + 1)
	}

	result, err := MACDEXT(input, 12, MASMA, 26, MASMA, 9, MAEMA)
	if err != nil {
		t.Fatalf("MACDEXT returned error: %v", err)
	}
	if len(result.MACD) != len(input) {
		t.Errorf("MACD length mismatch")
	}
}

func TestAD_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14}
	low := []float64{8, 9, 10, 11, 12}
	close := []float64{9, 11, 10, 12, 13}
	volume := []float64{100, 200, 150, 300, 250}

	result, err := AD(high, low, close, volume)
	if err != nil {
		t.Fatalf("AD returned error: %v", err)
	}
	if len(result) != len(close) {
		t.Errorf("length mismatch")
	}
}

func TestADOSC_Basic(t *testing.T) {
	high := []float64{10, 12, 11, 13, 14, 15, 14, 13, 16, 15}
	low := []float64{8, 9, 10, 11, 12, 13, 12, 11, 14, 13}
	close := []float64{9, 11, 10, 12, 13, 14, 13, 12, 15, 14}
	volume := []float64{100, 200, 150, 300, 250, 400, 350, 200, 500, 300}

	result, err := ADOSC(high, low, close, volume, 3, 10)
	if err != nil {
		t.Fatalf("ADOSC returned error: %v", err)
	}
	if len(result) != len(close) {
		t.Errorf("length mismatch")
	}
}

func TestADXR_Basic(t *testing.T) {
	high := make([]float64, 60)
	low := make([]float64, 60)
	close := make([]float64, 60)
	for i := range high {
		base := float64(i) * 0.1
		high[i] = 50 + base + math.Sin(float64(i)*0.3)*2
		low[i] = 48 + base + math.Cos(float64(i)*0.3)*2
		close[i] = 49 + base
	}

	result, err := ADXR(high, low, close, 14)
	if err != nil {
		t.Fatalf("ADXR returned error: %v", err)
	}
	if len(result) != len(close) {
		t.Errorf("length mismatch: %d vs %d", len(result), len(close))
	}
}

func TestT3_Basic(t *testing.T) {
	input := make([]float64, 40)
	for i := range input {
		input[i] = float64(i + 1)
	}

	result, err := T3(input, 5, 0.7)
	if err != nil {
		t.Fatalf("T3 returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}

func TestTRIX_Basic(t *testing.T) {
	input := make([]float64, 30)
	for i := range input {
		input[i] = 100.0 + float64(i)*0.3
	}

	result, err := TRIX(input, 5)
	if err != nil {
		t.Fatalf("TRIX returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}

func TestMACDFIX_Basic(t *testing.T) {
	input := make([]float64, 50)
	for i := range input {
		input[i] = float64(i + 1)
	}

	result, err := MACDFIX(input, 9)
	if err != nil {
		t.Fatalf("MACDFIX returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}
