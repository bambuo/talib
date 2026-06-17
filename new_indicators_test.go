package talib

import (
	"math"
	"testing"
)

// =============================================================================
// 价格变换测试
// =============================================================================

func TestAVGPRICE_Basic(t *testing.T) {
	open := []float64{10, 11, 12, 13}
	high := []float64{12, 13, 14, 15}
	low := []float64{9, 10, 11, 12}
	close := []float64{11, 12, 13, 14}

	result, err := AVGPRICE(open, high, low, close)
	if err != nil {
		t.Fatalf("AVGPRICE returned error: %v", err)
	}
	// (10 + 12 + 9 + 11) / 4 = 10.5
	expected := 10.5
	if math.Abs(result[0]-expected) > 1e-10 {
		t.Errorf("AVGPRICE[0]: expected %f, got %f", expected, result[0])
	}
}

func TestMEDPRICE_Basic(t *testing.T) {
	high := []float64{12, 14, 16}
	low := []float64{8, 10, 12}

	result, err := MEDPRICE(high, low)
	if err != nil {
		t.Fatalf("MEDPRICE returned error: %v", err)
	}
	// (12 + 8) / 2 = 10
	if math.Abs(result[0]-10.0) > 1e-10 {
		t.Errorf("MEDPRICE[0]: expected 10, got %f", result[0])
	}
}

func TestTYPPRICE_Basic(t *testing.T) {
	high := []float64{12, 14, 16}
	low := []float64{8, 10, 12}
	close := []float64{10, 12, 14}

	result, err := TYPPRICE(high, low, close)
	if err != nil {
		t.Fatalf("TYPPRICE returned error: %v", err)
	}
	// (12 + 8 + 10) / 3 = 10
	if math.Abs(result[0]-10.0) > 1e-10 {
		t.Errorf("TYPPRICE[0]: expected 10, got %f", result[0])
	}
}

func TestWCLPRICE_Basic(t *testing.T) {
	high := []float64{12, 14}
	low := []float64{8, 10}
	close := []float64{10, 12}

	result, err := WCLPRICE(high, low, close)
	if err != nil {
		t.Fatalf("WCLPRICE returned error: %v", err)
	}
	// (12 + 8 + 2*10) / 4 = 40/4 = 10
	if math.Abs(result[0]-10.0) > 1e-10 {
		t.Errorf("WCLPRICE[0]: expected 10, got %f", result[0])
	}
}

// =============================================================================
// 数学运算测试
// =============================================================================

func TestADD_Basic(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	result, err := ADD(a, b)
	if err != nil {
		t.Fatalf("ADD returned error: %v", err)
	}
	for i, v := range result {
		if v != a[i]+b[i] {
			t.Errorf("ADD[%d]: expected %f, got %f", i, a[i]+b[i], v)
		}
	}
}

func TestSUB_Basic(t *testing.T) {
	a := []float64{10, 20, 30}
	b := []float64{1, 2, 3}
	result, err := SUB(a, b)
	if err != nil {
		t.Fatalf("SUB returned error: %v", err)
	}
	for i, v := range result {
		if v != a[i]-b[i] {
			t.Errorf("SUB[%d]: expected %f, got %f", i, a[i]-b[i], v)
		}
	}
}

func TestMULT_Basic(t *testing.T) {
	a := []float64{2, 3, 4}
	b := []float64{5, 6, 7}
	result, err := MULT(a, b)
	if err != nil {
		t.Fatalf("MULT returned error: %v", err)
	}
	for i, v := range result {
		if v != a[i]*b[i] {
			t.Errorf("MULT[%d]: expected %f, got %f", i, a[i]*b[i], v)
		}
	}
}

func TestDIV_Basic(t *testing.T) {
	a := []float64{10, 20, 30}
	b := []float64{2, 4, 0}
	result, err := DIV(a, b)
	if err != nil {
		t.Fatalf("DIV returned error: %v", err)
	}
	if result[0] != 5.0 {
		t.Errorf("DIV[0]: expected 5, got %f", result[0])
	}
	if result[1] != 5.0 {
		t.Errorf("DIV[1]: expected 5, got %f", result[1])
	}
	if !math.IsNaN(result[2]) {
		t.Errorf("DIV[2]: expected NaN, got %f", result[2])
	}
}

func TestSUM_Basic(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	period := 3
	result, err := SUM(input, period)
	if err != nil {
		t.Fatalf("SUM returned error: %v", err)
	}
	// SUM[2] = 1+2+3 = 6
	if math.Abs(result[2]-6.0) > 1e-10 {
		t.Errorf("SUM[2]: expected 6, got %f", result[2])
	}
	// SUM[3] = 2+3+4 = 9
	if math.Abs(result[3]-9.0) > 1e-10 {
		t.Errorf("SUM[3]: expected 9, got %f", result[3])
	}
}

func TestMAX_Basic(t *testing.T) {
	input := []float64{1, 5, 3, 7, 2, 9, 4}
	period := 3
	result, err := MAX(input, period)
	if err != nil {
		t.Fatalf("MAX returned error: %v", err)
	}
	if math.Abs(result[2]-5.0) > 1e-10 {
		t.Errorf("MAX[2]: expected 5, got %f", result[2])
	}
	if math.Abs(result[3]-7.0) > 1e-10 {
		t.Errorf("MAX[3]: expected 7, got %f", result[3])
	}
	if math.Abs(result[5]-9.0) > 1e-10 {
		t.Errorf("MAX[5]: expected 9, got %f", result[5])
	}
}

func TestMIN_Basic(t *testing.T) {
	input := []float64{5, 3, 7, 2, 9, 4, 1}
	period := 3
	result, err := MIN(input, period)
	if err != nil {
		t.Fatalf("MIN returned error: %v", err)
	}
	if math.Abs(result[2]-3.0) > 1e-10 {
		t.Errorf("MIN[2]: expected 3, got %f", result[2])
	}
	if math.Abs(result[3]-2.0) > 1e-10 {
		t.Errorf("MIN[3]: expected 2, got %f", result[3])
	}
}

func TestMAXINDEX_Basic(t *testing.T) {
	input := []float64{1, 5, 3, 7, 2}
	period := 3
	result, err := MAXINDEX(input, period)
	if err != nil {
		t.Fatalf("MAXINDEX returned error: %v", err)
	}
	// max of [1,5,3] is 5 at index 1
	if result[2] != 1.0 {
		t.Errorf("MAXINDEX[2]: expected 1, got %f", result[2])
	}
}

func TestMININDEX_Basic(t *testing.T) {
	input := []float64{5, 3, 7, 2, 9}
	period := 3
	result, err := MININDEX(input, period)
	if err != nil {
		t.Fatalf("MININDEX returned error: %v", err)
	}
	// min of [5,3,7] is 3 at index 1
	if result[2] != 1.0 {
		t.Errorf("MININDEX[2]: expected 1, got %f", result[2])
	}
}

func TestMINMAX_Basic(t *testing.T) {
	input := []float64{1, 5, 3, 7, 2, 9, 4}
	period := 3
	result, err := MINMAX(input, period)
	if err != nil {
		t.Fatalf("MINMAX returned error: %v", err)
	}
	if math.Abs(result.Min[2]-1.0) > 1e-10 {
		t.Errorf("MIN[2]: expected 1, got %f", result.Min[2])
	}
	if math.Abs(result.Max[2]-5.0) > 1e-10 {
		t.Errorf("MAX[2]: expected 5, got %f", result.Max[2])
	}
}

// =============================================================================
// 数学变换测试
// =============================================================================

func TestSQRT_Basic(t *testing.T) {
	input := []float64{4, 9, 16}
	result, err := SQRT(input)
	if err != nil {
		t.Fatalf("SQRT returned error: %v", err)
	}
	if math.Abs(result[0]-2.0) > 1e-10 {
		t.Errorf("SQRT[0]: expected 2, got %f", result[0])
	}
}

func TestLN_Basic(t *testing.T) {
	input := []float64{1, math.E}
	result, err := LN(input)
	if err != nil {
		t.Fatalf("LN returned error: %v", err)
	}
	if math.Abs(result[0]-0.0) > 1e-10 {
		t.Errorf("LN[0]: expected 0, got %f", result[0])
	}
	if math.Abs(result[1]-1.0) > 1e-10 {
		t.Errorf("LN[1]: expected 1, got %f", result[1])
	}
}

func TestEXP_Basic(t *testing.T) {
	input := []float64{0, 1}
	result, err := EXP(input)
	if err != nil {
		t.Fatalf("EXP returned error: %v", err)
	}
	if math.Abs(result[0]-1.0) > 1e-10 {
		t.Errorf("EXP[0]: expected 1, got %f", result[0])
	}
	if math.Abs(result[1]-math.E) > 1e-10 {
		t.Errorf("EXP[1]: expected e, got %f", result[1])
	}
}

func TestSIN_Basic(t *testing.T) {
	input := []float64{0, math.Pi / 2}
	result, err := SIN(input)
	if err != nil {
		t.Fatalf("SIN returned error: %v", err)
	}
	if math.Abs(result[0]-0.0) > 1e-10 {
		t.Errorf("SIN[0]: expected 0, got %f", result[0])
	}
	if math.Abs(result[1]-1.0) > 1e-10 {
		t.Errorf("SIN[1]: expected 1, got %f", result[1])
	}
}

func TestCOS_Basic(t *testing.T) {
	input := []float64{0}
	result, err := COS(input)
	if err != nil {
		t.Fatalf("COS returned error: %v", err)
	}
	if math.Abs(result[0]-1.0) > 1e-10 {
		t.Errorf("COS[0]: expected 1, got %f", result[0])
	}
}

func TestFLOOR_Basic(t *testing.T) {
	input := []float64{1.5, 2.9, -1.3}
	result, err := FLOOR(input)
	if err != nil {
		t.Fatalf("FLOOR returned error: %v", err)
	}
	if result[0] != 1.0 || result[1] != 2.0 || result[2] != -2.0 {
		t.Errorf("FLOOR: unexpected values %v", result)
	}
}

func TestCEIL_Basic(t *testing.T) {
	input := []float64{1.5, 2.1, -1.3}
	result, err := CEIL(input)
	if err != nil {
		t.Fatalf("CEIL returned error: %v", err)
	}
	if result[0] != 2.0 || result[1] != 3.0 || result[2] != -1.0 {
		t.Errorf("CEIL: unexpected values %v", result)
	}
}

// =============================================================================
// 动量测试
// =============================================================================

func TestBOP_Basic(t *testing.T) {
	open := []float64{10, 11, 10}
	high := []float64{12, 14, 13}
	low := []float64{9, 10, 9}
	close := []float64{11, 12, 12}

	result, err := BOP(open, high, low, close)
	if err != nil {
		t.Fatalf("BOP returned error: %v", err)
	}
	// BOP[0] = (11-10)/(12-9) = 1/3
	if math.Abs(result[0]-1.0/3.0) > 1e-10 {
		t.Errorf("BOP[0]: expected 1/3, got %f", result[0])
	}
}

func TestMOM_Basic(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7}
	period := 3

	result, err := MOM(input, period)
	if err != nil {
		t.Fatalf("MOM returned error: %v", err)
	}
	// MOM[3] = 4 - 1 = 3
	if math.Abs(result[3]-3.0) > 1e-10 {
		t.Errorf("MOM[3]: expected 3, got %f", result[3])
	}
	// MOM[6] = 7 - 4 = 3
	if math.Abs(result[6]-3.0) > 1e-10 {
		t.Errorf("MOM[6]: expected 3, got %f", result[6])
	}
}

func TestROC_Basic(t *testing.T) {
	input := []float64{10, 12, 11, 13, 14, 15}
	period := 3

	result, err := ROC(input, period)
	if err != nil {
		t.Fatalf("ROC returned error: %v", err)
	}
	// ROC[3] = (13/10 - 1) * 100 = 30
	if math.Abs(result[3]-30.0) > 1e-10 {
		t.Errorf("ROC[3]: expected 30, got %f", result[3])
	}
}

func TestROCR_Basic(t *testing.T) {
	input := []float64{10, 12, 11, 13, 14, 15}
	period := 3

	result, err := ROCR(input, period)
	if err != nil {
		t.Fatalf("ROCR returned error: %v", err)
	}
	// ROCR[3] = 13/10 = 1.3
	if math.Abs(result[3]-1.3) > 1e-10 {
		t.Errorf("ROCR[3]: expected 1.3, got %f", result[3])
	}
}

func TestCMO_Basic(t *testing.T) {
	input := []float64{10, 12, 11, 14, 13, 15, 14, 16, 15, 17}
	period := 3

	result, err := CMO(input, period)
	if err != nil {
		t.Fatalf("CMO returned error: %v", err)
	}

	// Check output length
	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}

	// CMO range check
	for i, v := range result {
		if math.IsNaN(v) {
			continue
		}
		if v < -100 || v > 100 {
			t.Errorf("CMO[%d] = %f, out of [-100, 100] range", i, v)
		}
	}
}

// =============================================================================
// 统计测试
// =============================================================================

func TestSTDDEV_Basic(t *testing.T) {
	input := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	period := 4

	result, err := STDDEV(input, period, 1.0)
	if err != nil {
		t.Fatalf("STDDEV returned error: %v", err)
	}

	// 窗口 [2,4,4,4]: mean=3.5, std = sqrt(((2-3.5)²+(4-3.5)²+(4-3.5)²+(4-3.5)²)/4)
	// = sqrt((2.25+0.25+0.25+0.25)/4) = sqrt(3/4) = sqrt(0.75) ≈ 0.866
	if math.IsNaN(result[3]) {
		t.Error("STDDEV[3] should not be NaN")
	}
}

func TestVAR_Basic(t *testing.T) {
	input := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	period := 4

	result, err := VAR(input, period)
	if err != nil {
		t.Fatalf("VAR returned error: %v", err)
	}

	// 窗口 [2,4,4,4]: mean=3.5, var = 0.75
	expected := 0.75
	if math.Abs(result[3]-expected) > 1e-10 {
		t.Errorf("VAR[3]: expected %f, got %f", expected, result[3])
	}
}

func TestCORREL_Basic(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{2, 4, 6, 8, 10, 12, 14, 16}
	period := 4

	result, err := CORREL(x, y, period)
	if err != nil {
		t.Fatalf("CORREL returned error: %v", err)
	}

	// 完美正相关
	if !math.IsNaN(result[3]) && math.Abs(result[3]-1.0) > 0.1 {
		t.Errorf("CORREL[3]: expected near 1.0, got %f", result[3])
	}
}

func TestBETA_Basic(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	y := []float64{2, 4, 6, 8, 10, 12, 14, 16}
	period := 4

	result, err := BETA(x, y, period)
	if err != nil {
		t.Fatalf("BETA returned error: %v", err)
	}

	if len(result) != len(x) {
		t.Errorf("expected length %d, got %d", len(x), len(result))
	}
}

func TestLINEARREG_Basic(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	period := 4

	result, err := LINEARREG(input, period)
	if err != nil {
		t.Fatalf("LINEARREG returned error: %v", err)
	}

	// 窗口 [1,2,3,4]: 最后一个预测值
	if math.IsNaN(result[3]) {
		t.Error("LINEARREG[3] should not be NaN")
	}
}

func TestLINEARREG_SLOPE_Basic(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	period := 4

	result, err := LINEARREG_SLOPE(input, period)
	if err != nil {
		t.Fatalf("LINEARREG_SLOPE returned error: %v", err)
	}
	// 对于 [1,2,3,4], 斜率应为 1.0
	if math.Abs(result[3]-1.0) > 1e-10 {
		t.Errorf("LINEARREG_SLOPE[3]: expected 1.0, got %f", result[3])
	}
}

func TestTSF_Basic(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	period := 4

	result, err := TSF(input, period)
	if err != nil {
		t.Fatalf("TSF returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}

	// 检查输出长度和 NaN 约定
	for i := 0; i < period-1; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("TSF[%d]: expected NaN, got %f", i, result[i])
		}
	}
}

// =============================================================================
// 重叠研究测试
// =============================================================================

func TestMIDPOINT_Basic(t *testing.T) {
	input := []float64{10, 14, 12, 18, 13}
	period := 3

	result, err := MIDPOINT(input, period)
	if err != nil {
		t.Fatalf("MIDPOINT returned error: %v", err)
	}
	// Window [10,14,12]: highest=14, lowest=10, midpoint = (14+10)/2 = 12
	if math.Abs(result[2]-12.0) > 1e-10 {
		t.Errorf("MIDPOINT[2]: expected 12, got %f", result[2])
	}
	// Window [14,12,18]: highest=18, lowest=12, midpoint = (18+12)/2 = 15
	if math.Abs(result[3]-15.0) > 1e-10 {
		t.Errorf("MIDPOINT[3]: expected 15, got %f", result[3])
	}
}

func TestMIDPRICE_Basic(t *testing.T) {
	high := []float64{12, 14, 13, 18, 16}
	low := []float64{9, 10, 11, 12, 13}
	period := 3

	result, err := MIDPRICE(high, low, period)
	if err != nil {
		t.Fatalf("MIDPRICE returned error: %v", err)
	}
	// Window high=[12,14,13] -> highest=14, low=[9,10,11] -> lowest=9
	// midpoint = (14+9)/2 = 11.5
	if math.Abs(result[2]-11.5) > 1e-10 {
		t.Errorf("MIDPRICE[2]: expected 11.5, got %f", result[2])
	}
}

func TestSAR_Basic(t *testing.T) {
	high := []float64{10, 12, 13, 15, 14, 16, 15, 18, 17, 19}
	low := []float64{8, 9, 10, 12, 11, 13, 12, 14, 13, 15}

	result, err := SAR(high, low, 0.02, 0.20)
	if err != nil {
		t.Fatalf("SAR returned error: %v", err)
	}

	if len(result) != len(high) {
		t.Errorf("expected length %d, got %d", len(high), len(result))
	}

	// SAR 应在第一个 period 之后产生非 NaN 值
	sarCount := 0
	for _, v := range result {
		if !math.IsNaN(v) {
			sarCount++
		}
	}
	if sarCount < 3 {
		t.Errorf("expected at least 3 non-NaN SAR values, got %d", sarCount)
	}
}

func TestSAREXT_Basic(t *testing.T) {
	high := []float64{10, 12, 13, 15, 14, 16, 15, 18, 17, 19}
	low := []float64{8, 9, 10, 12, 11, 13, 12, 14, 13, 15}

	result, err := SAREXT(high, low, 0.02, 0.20, 0.0)
	if err != nil {
		t.Fatalf("SAREXT returned error: %v", err)
	}

	if len(result) != len(high) {
		t.Errorf("expected length %d, got %d", len(high), len(result))
	}
}

func TestMAVP_Basic(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	periods := []float64{2, 2, 3, 3, 4, 4, 5, 5, 3, 3}

	result, err := MAVP(input, periods, 2, 5, MASMA)
	if err != nil {
		t.Fatalf("MAVP returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

// =============================================================================
// KAMA 测试
// =============================================================================

func TestKAMA_Basic(t *testing.T) {
	input := []float64{10, 10.5, 10.3, 10.6, 10.8, 11.0, 11.2, 11.0, 10.8, 11.1, 11.3, 11.5}
	period := 5

	result, err := KAMA(input, period)
	if err != nil {
		t.Fatalf("KAMA returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}

	// 检查 NaN 约定
	for i := 0; i < period; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("KAMA[%d]: expected NaN, got %f", i, result[i])
		}
	}
}

func TestMA_KAMA(t *testing.T) {
	input := []float64{10, 10.5, 10.3, 10.6, 10.8, 11.0, 11.2, 11.0, 10.8, 11.1}
	period := 5

	// 通过 MA 调度器测试 KAMA 类型
	result, err := MA(input, period, MAKAMA)
	if err != nil {
		t.Fatalf("MA(KAMA) returned error: %v", err)
	}

	// 与直接 KAMA 调用比较
	kamaResult, err := KAMA(input, period)
	if err != nil {
		t.Fatalf("KAMA returned error: %v", err)
	}

	for i, v := range result {
		if math.IsNaN(v) != math.IsNaN(kamaResult[i]) {
			t.Errorf("MA(KAMA)[%d]: NaN mismatch", i)
		} else if !math.IsNaN(v) && math.Abs(v-kamaResult[i]) > 1e-10 {
			t.Errorf("MA(KAMA)[%d]: expected %f, got %f", i, kamaResult[i], v)
		}
	}
}

// =============================================================================
// MAMA 测试
// =============================================================================

func TestMAMA_Basic(t *testing.T) {
	// 需要足够的 MAMA 数据（至少 14 根柱线才有有效输出）
	input := []float64{
		10, 10.2, 10.1, 10.3, 10.5, 10.4, 10.6,
		10.8, 10.7, 10.9, 11.0, 10.8, 11.1, 11.2,
		11.0, 11.3, 11.5, 11.4, 11.6, 11.8,
	}

	result, err := MAMA(input, 0.5, 0.05)
	if err != nil {
		t.Fatalf("MAMA returned error: %v", err)
	}

	if len(result.MAMA) != len(input) {
		t.Errorf("expected MAMA length %d, got %d", len(input), len(result.MAMA))
	}
	if len(result.FAMA) != len(input) {
		t.Errorf("expected FAMA length %d, got %d", len(input), len(result.FAMA))
	}

	// 检查 NaN 约定：前 13 个元素应为 NaN
	for i := 0; i < 13; i++ {
		if !math.IsNaN(result.MAMA[i]) {
			t.Errorf("MAMA[%d]: expected NaN, got %f", i, result.MAMA[i])
		}
	}

	// 检查 lookback 后是否产生有效值
	validCount := 0
	for _, v := range result.MAMA {
		if !math.IsNaN(v) {
			validCount++
		}
	}
	if validCount == 0 {
		t.Error("MAMA produced no valid values")
	}
}

func TestMA_MAMA(t *testing.T) {
	input := []float64{
		10, 10.2, 10.1, 10.3, 10.5, 10.4, 10.6,
		10.8, 10.7, 10.9, 11.0, 10.8, 11.1, 11.2,
		11.0, 11.3, 11.5, 11.4, 11.6, 11.8,
	}

	// 通过 MA 调度器测试 MAMA 类型
	result, err := MA(input, 5, MAMAMA) // MAMA 忽略 period 参数
	if err != nil {
		t.Fatalf("MA(MAMA) returned error: %v", err)
	}

	if len(result) != len(input) {
		t.Errorf("expected length %d, got %d", len(input), len(result))
	}
}

// =============================================================================
// DX 测试
// =============================================================================

func TestDX_Basic(t *testing.T) {
	high := []float64{50, 52, 51, 53, 54, 55, 54, 53, 56, 55, 57, 56, 58, 57, 59}
	low := []float64{48, 50, 49, 51, 52, 53, 52, 51, 54, 53, 55, 54, 56, 55, 57}
	close := []float64{49, 51, 50, 52, 53, 54, 53, 52, 55, 54, 56, 55, 57, 56, 58}
	period := 3

	result, err := DX(high, low, close, period)
	if err != nil {
		t.Fatalf("DX returned error: %v", err)
	}

	if len(result) != len(close) {
		t.Errorf("expected length %d, got %d", len(close), len(result))
	}
}

// =============================================================================
// 错误/边界情况测试
// =============================================================================

func TestAVGPRICE_LengthMismatch(t *testing.T) {
	open := []float64{1, 2}
	high := []float64{1, 2, 3}
	_, err := AVGPRICE(open, high, high, high)
	if err == nil {
		t.Error("expected error for length mismatch")
	}
}

func TestMEDPRICE_InvalidInput(t *testing.T) {
	_, err := MEDPRICE(nil, []float64{1, 2})
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestMOM_InvalidPeriod(t *testing.T) {
	_, err := MOM([]float64{1, 2, 3}, 5)
	if err == nil {
		t.Error("expected error for period exceeding input length")
	}
}

func TestMIDPOINT_NaNLeading(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5}
	period := 3
	result, _ := MIDPOINT(input, period)
	for i := 0; i < period-1; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("MIDPOINT[%d]: expected NaN, got %f", i, result[i])
		}
	}
}

func TestSTDDEV_NaNLeading(t *testing.T) {
	input := []float64{1, 2, 3, 4, 5}
	period := 3
	result, _ := STDDEV(input, period, 1.0)
	for i := 0; i < period-1; i++ {
		if !math.IsNaN(result[i]) {
			t.Errorf("STDDEV[%d]: expected NaN, got %f", i, result[i])
		}
	}
}
