package talib

import (
	"math"
	"testing"
)

func TestHT_TRENDLINE_Basic(t *testing.T) {
	input := make([]float64, 100)
	for i := range input {
		input[i] = 100.0 + math.Sin(float64(i)*0.1)*5.0
	}

	result, err := HT_TRENDLINE(input)
	if err != nil {
		t.Fatalf("HT_TRENDLINE returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch: %d vs %d", len(result), len(input))
	}
}

func TestHT_TRENDMODE_Basic(t *testing.T) {
	input := make([]float64, 100)
	for i := range input {
		input[i] = 100.0 + float64(i)*0.1 + math.Sin(float64(i)*0.2)*2.0
	}

	result, err := HT_TRENDMODE(input)
	if err != nil {
		t.Fatalf("HT_TRENDMODE returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}

	// 模式应仅为 0 或 1
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v != 0 && v != 1 {
			t.Errorf("TRENDMODE[%d] = %f, should be 0 or 1", i, v)
		}
	}
}

func TestHT_DCPERIOD_Basic(t *testing.T) {
	input := make([]float64, 100)
	for i := range input {
		input[i] = 100.0 + math.Sin(float64(i)*0.3)*3.0
	}

	result, err := HT_DCPERIOD(input)
	if err != nil {
		t.Fatalf("HT_DCPERIOD returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}

func TestHT_DCPHASE_Basic(t *testing.T) {
	input := make([]float64, 100)
	for i := range input {
		input[i] = 100.0 + math.Sin(float64(i)*0.15)*4.0
	}

	result, err := HT_DCPHASE(input)
	if err != nil {
		t.Fatalf("HT_DCPHASE returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}

func TestHT_PHASOR_Basic(t *testing.T) {
	input := make([]float64, 50)
	for i := range input {
		input[i] = 100.0 + math.Sin(float64(i)*0.2)*3.0
	}

	result, err := HT_PHASOR(input)
	if err != nil {
		t.Fatalf("HT_PHASOR returned error: %v", err)
	}
	if len(result.InPhase) != len(input) {
		t.Errorf("InPhase length mismatch")
	}
	if len(result.Quadrature) != len(input) {
		t.Errorf("Quadrature length mismatch")
	}
}

func TestHT_SINE_Basic(t *testing.T) {
	input := make([]float64, 100)
	for i := range input {
		input[i] = 100.0 + math.Sin(float64(i)*0.12)*5.0
	}

	result, err := HT_SINE(input)
	if err != nil {
		t.Fatalf("HT_SINE returned error: %v", err)
	}
	if len(result.Sine) != len(input) {
		t.Errorf("Sine length mismatch")
	}
	if len(result.LeadSine) != len(input) {
		t.Errorf("LeadSine length mismatch")
	}

	// Sine 应在 [-1, 1] 范围内
	for i, v := range result.Sine {
		if math.IsNaN(v) {
			continue
		}
		if v < -1.0 || v > 1.0 {
			t.Errorf("HT_SINE Sine[%d] = %f, out of range", i, v)
		}
	}
}

func TestHT_DCCOMPONENT_Basic(t *testing.T) {
	input := make([]float64, 100)
	for i := range input {
		input[i] = 100.0 + float64(i)*0.05 + math.Sin(float64(i)*0.1)*2.0
	}

	result, err := HT_DCCOMPONENT(input)
	if err != nil {
		t.Fatalf("HT_DCCOMPONENT returned error: %v", err)
	}
	if len(result) != len(input) {
		t.Errorf("length mismatch")
	}
}

func TestHT_ShortInput(t *testing.T) {
	input := []float64{1, 2, 3}
	_, err := HT_TRENDLINE(input)
	if err == nil {
		t.Error("expected error for short input")
	}
}
