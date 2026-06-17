package cdl

import (
	"testing"
)

// =====================================================================
// 共享测试辅助函数
// =====================================================================

// makeTestOHLC 创建用于测试的 OHLC 数据。
// n: 蜡烛数量（末尾补零）
// ohlc: 每个元素为 [o, h, l, c] 表示一根蜡烛
func makeTestOHLC(ohlc ...[4]float64) (open, high, low, close []float64) {
	n := len(ohlc)
	open = make([]float64, n)
	high = make([]float64, n)
	low = make([]float64, n)
	close = make([]float64, n)
	for i, c := range ohlc {
		open[i] = c[0]
		high[i] = c[1]
		low[i] = c[2]
		close[i] = c[3]
	}
	return
}

func assertNoError(t *testing.T, err error, name string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: unexpected error: %v", name, err)
	}
}

func assertError(t *testing.T, err error, name string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error but got nil", name)
	}
}

func assertLen(t *testing.T, got, want int, name string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: output length = %d, want %d", name, got, want)
	}
}

func assertSignalRange(t *testing.T, out []int, name string) {
	t.Helper()
	for i, v := range out {
		if v != Bullish && v != Bearish && v != Neutral {
			t.Errorf("%s: invalid signal %d at index %d", name, v, i)
		}
	}
}

// =====================================================================
// 输入验证测试
// =====================================================================

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name string
		fn   func(open, high, low, close []float64) ([]int, error)
		lb   func() int
	}{
		{"CDLDOJI", CDLDOJI, CDLDOJILookback},

		{"CDLENGULFING", CDLENGULFING, CDLENGULFINGLookback},
		{"CDLHARAMI", CDLHARAMI, CDLHARAMILookback},
		{"CDLHARAMICROSS", CDLHARAMICROSS, CDLHARAMICROSSLookback},
		{"CDLPIERCING", CDLPIERCING, CDLPIERCINGLookback},
		{"CDLDARKCLOUDCOVER", CDLDARKCLOUDCOVER, CDLDARKCLOUDCOVERLookback},
		{"CDLCOUNTERATTACK", CDLCOUNTERATTACK, CDLCOUNTERATTACKLookback},

		{"CDLMORNINGSTAR", CDLMORNINGSTAR, CDLMORNINGSTARLookback},
		{"CDLEVENINGSTAR", CDLEVENINGSTAR, CDLEVENINGSTARLookback},
		{"CDLABANDONEDBABY", CDLABANDONEDBABY, CDLABANDONEDBABYLookback},

		{"CDL3WHITESOLDIERS", CDL3WHITESOLDIERS, CDL3WHITESOLDIERSLookback},
		{"CDL3BLACKCROWS", CDL3BLACKCROWS, CDL3BLACKCROWSLookback},
		{"CDL3INSIDE", CDL3INSIDE, CDL3INSIDELookback},

		{"CDLBREAKAWAY", CDLBREAKAWAY, CDLBREAKAWAYLookback},
		{"CDLMATHOLD", CDLMATHOLD, CDLMATHOLDLookback},
		{"CDLRISEFALL3METHODS", CDLRISEFALL3METHODS, CDLRISEFALL3METHODSLookback},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_nil", func(t *testing.T) {
			_, err := tt.fn(nil, nil, nil, nil)
			assertError(t, err, tt.name)
		})
		t.Run(tt.name+"_empty", func(t *testing.T) {
			o := []float64{}
			_, err := tt.fn(o, o, o, o)
			assertError(t, err, tt.name)
		})
		t.Run(tt.name+"_mismatch", func(t *testing.T) {
			o := []float64{1, 2, 3}
			h := []float64{1, 2}
			_, err := tt.fn(o, h, o, o)
			assertError(t, err, tt.name)
		})
	}
}

// =====================================================================
// 单K线形态测试
// =====================================================================

func TestCDLDOJI(t *testing.T) {
	// 十字星：开盘价=收盘价，小实体
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 90, 105},   // 非十字星
		[4]float64{100, 102, 98, 100},   // 十字星（实体=0）
		[4]float64{100, 105, 95, 100.1}, // 近似十字星
		[4]float64{100, 120, 80, 105},   // 非十字星
	)
	out, err := CDLDOJI(open, high, low, close)
	assertNoError(t, err, "CDLDOJI")
	assertLen(t, len(out), 4, "CDLDOJI")
	if out[0] != Neutral {
		t.Errorf("CDLDOJI: candle 0 should not be doji, got %d", out[0])
	}
	if out[1] != Bullish {
		t.Errorf("CDLDOJI: candle 1 should be doji, got %d", out[1])
	}
	if out[3] != Neutral {
		t.Errorf("CDLDOJI: candle 3 should not be doji, got %d", out[3])
	}
	assertSignalRange(t, out, "CDLDOJI")
}

func TestCDLHAMMER(t *testing.T) {
	// 锤子线：小实体在顶部，长下影线，处于下降趋势
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102}, // 阴线（下降趋势）
		[4]float64{102, 103, 97, 101},  // 阴线（下降趋势）
		[4]float64{101, 102, 95, 100},  // 阴线（下降趋势）
		[4]float64{98, 99, 90, 97.5},   // 类锤子线
		[4]float64{97, 98, 94, 97.8},   // 看跌上下文持续
	)
	out, err := CDLHAMMER(open, high, low, close)
	assertNoError(t, err, "CDLHAMMER")
	assertLen(t, len(out), 5, "CDLHAMMER")
	assertSignalRange(t, out, "CDLHAMMER")
}

func TestCDLHANGINGMAN(t *testing.T) {
	// 吊人线：与锤子线相同形态，但处于上升趋势
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 104},  // 阳线（上升趋势）
		[4]float64{104, 108, 102, 106}, // 阳线（上升趋势）
		[4]float64{106, 109, 104, 108}, // 阳线（上升趋势）
		[4]float64{106, 107, 100, 105}, // 吊人线形态
	)
	out, err := CDLHANGINGMAN(open, high, low, close)
	assertNoError(t, err, "CDLHANGINGMAN")
	assertLen(t, len(out), 4, "CDLHANGINGMAN")
	assertSignalRange(t, out, "CDLHANGINGMAN")
}

func TestCDLMARUBOZU(t *testing.T) {
	// 光脚光头：长实体，无影线
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 100, 110}, // 白色光脚光头
		[4]float64{110, 110, 100, 100}, // 黑色光脚光头
		[4]float64{100, 108, 97, 105},  // 非光脚光头
	)
	out, err := CDLMARUBOZU(open, high, low, close)
	assertNoError(t, err, "CDLMARUBOZU")
	assertLen(t, len(out), 3, "CDLMARUBOZU")
	if out[0] != Bullish {
		t.Errorf("CDLMARUBOZU: candle 0 should be bullish marubozu, got %d", out[0])
	}
	if out[1] != Bearish {
		t.Errorf("CDLMARUBOZU: candle 1 should be bearish marubozu, got %d", out[1])
	}
	assertSignalRange(t, out, "CDLMARUBOZU")
}

func TestCDLSHOOTINGSTAR(t *testing.T) {
	// 流星线：小实体在底部，长上影线，处于上升趋势
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 104},  // 阳线（上升趋势）
		[4]float64{104, 108, 102, 106}, // 阳线（上升趋势）
		[4]float64{106, 109, 104, 108}, // 阳线（上升趋势）
		[4]float64{107, 115, 106, 108}, // 流星线
	)
	out, err := CDLSHOOTINGSTAR(open, high, low, close)
	assertNoError(t, err, "CDLSHOOTINGSTAR")
	assertLen(t, len(out), 4, "CDLSHOOTINGSTAR")
	assertSignalRange(t, out, "CDLSHOOTINGSTAR")
}

func TestCDLSPINNINGTOP(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 102, 95, 101},  // 陀螺线（小实体，有影线）
		[4]float64{100, 110, 100, 110}, // 光脚光头（非陀螺线）
	)
	out, err := CDLSPINNINGTOP(open, high, low, close)
	assertNoError(t, err, "CDLSPINNINGTOP")
	assertLen(t, len(out), 2, "CDLSPINNINGTOP")
	assertSignalRange(t, out, "CDLSPINNINGTOP")
}

func TestCDLSHORTLINE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 100, 110},  // 长实体
		[4]float64{100, 101, 99, 100.1}, // 非常短的实体
	)
	out, err := CDLSHORTLINE(open, high, low, close)
	assertNoError(t, err, "CDLSHORTLINE")
	assertLen(t, len(out), 2, "CDLSHORTLINE")
	assertSignalRange(t, out, "CDLSHORTLINE")
}

func TestCDLLONGLINE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 100, 110}, // 长阳线
		[4]float64{110, 110, 100, 100}, // 长阴线
		[4]float64{100, 102, 98, 101},  // 短线
	)
	out, err := CDLLONGLINE(open, high, low, close)
	assertNoError(t, err, "CDLLONGLINE")
	assertLen(t, len(out), 3, "CDLLONGLINE")
	if out[0] != Bullish {
		t.Errorf("CDLLONGLINE: candle 0 should be bullish, got %d", out[0])
	}
	if out[1] != Bearish {
		t.Errorf("CDLLONGLINE: candle 1 should be bearish, got %d", out[1])
	}
}

func TestDRAGONFLYDOJI(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 100.1, 90, 100}, // 蜻蜓十字星
	)
	out, err := CDLDRAGONFLYDOJI(open, high, low, close)
	assertNoError(t, err, "CDLDRAGONFLYDOJI")
	assertLen(t, len(out), 1, "CDLDRAGONFLYDOJI")
	assertSignalRange(t, out, "CDLDRAGONFLYDOJI")
}

func TestGRAVESTONEDOJI(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 99.9, 100}, // 墓碑十字星
	)
	out, err := CDLGRAVESTONEDOJI(open, high, low, close)
	assertNoError(t, err, "CDLGRAVESTONEDOJI")
	assertLen(t, len(out), 1, "CDLGRAVESTONEDOJI")
	assertSignalRange(t, out, "CDLGRAVESTONEDOJI")
}

func TestCDLLONGLEGGEDDOJI(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 90, 100}, // 长腿十字星
	)
	out, err := CDLLONGLEGGEDDOJI(open, high, low, close)
	assertNoError(t, err, "CDLLONGLEGGEDDOJI")
	assertLen(t, len(out), 1, "CDLLONGLEGGEDDOJI")
	assertSignalRange(t, out, "CDLLONGLEGGEDDOJI")
}

func TestCDLRICKSHAWMAN(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 92, 100}, // 人力车十字星
	)
	out, err := CDLRICKSHAWMAN(open, high, low, close)
	assertNoError(t, err, "CDLRICKSHAWMAN")
	assertLen(t, len(out), 1, "CDLRICKSHAWMAN")
	assertSignalRange(t, out, "CDLRICKSHAWMAN")
}

func TestCDLHIGHWAVE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 92, 100.5}, // 长脚线
	)
	out, err := CDLHIGHWAVE(open, high, low, close)
	assertNoError(t, err, "CDLHIGHWAVE")
	assertLen(t, len(out), 1, "CDLHIGHWAVE")
	assertSignalRange(t, out, "CDLHIGHWAVE")
}

func TestCDLCLOSINGMARUBOZU(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 100, 109.5}, // 近似收盘光脚光头
	)
	out, err := CDLCLOSINGMARUBOZU(open, high, low, close)
	assertNoError(t, err, "CDLCLOSINGMARUBOZU")
	assertLen(t, len(out), 1, "CDLCLOSINGMARUBOZU")
	assertSignalRange(t, out, "CDLCLOSINGMARUBOZU")
}

func TestCDLINVERTEDHAMMER(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102}, // downtrend
		[4]float64{102, 103, 98, 100},  // downtrend
		[4]float64{100, 101, 96, 98},   // downtrend
		[4]float64{97, 105, 96, 98},    // inverted hammer
	)
	out, err := CDLINVERTEDHAMMER(open, high, low, close)
	assertNoError(t, err, "CDLINVERTEDHAMMER")
	assertLen(t, len(out), 4, "CDLINVERTEDHAMMER")
	assertSignalRange(t, out, "CDLINVERTEDHAMMER")
}

func TestCDLTAKURI(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102}, // downtrend
		[4]float64{102, 103, 97, 101},  // downtrend
		[4]float64{101, 102, 96, 98},   // downtrend
		[4]float64{97, 97.1, 85, 97},   // takuri
	)
	out, err := CDLTAKURI(open, high, low, close)
	assertNoError(t, err, "CDLTAKURI")
	assertLen(t, len(out), 4, "CDLTAKURI")
	assertSignalRange(t, out, "CDLTAKURI")
}

// =====================================================================
// 双K线形态测试
// =====================================================================

func TestCDLENGULFING(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 102, 98, 101},  // small white (prev)
		[4]float64{98, 103, 95, 103},   // large white (engulfs prev) - bullish
		[4]float64{105, 106, 100, 102}, // small white
		[4]float64{106, 107, 99, 100},  // large black (engulfs prev) - bearish
	)
	out, err := CDLENGULFING(open, high, low, close)
	assertNoError(t, err, "CDLENGULFING")
	assertLen(t, len(out), 4, "CDLENGULFING")
	assertSignalRange(t, out, "CDLENGULFING")
}

func TestCDLHARAMI(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 95, 98},   // large black
		[4]float64{97, 99, 96.5, 98.5}, // small white inside (bullish harami)
	)
	out, err := CDLHARAMI(open, high, low, close)
	assertNoError(t, err, "CDLHARAMI")
	assertLen(t, len(out), 2, "CDLHARAMI")
	assertSignalRange(t, out, "CDLHARAMI")
}

func TestCDLDARKCLOUDCOVER(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 98, 108},  // long white
		[4]float64{112, 113, 103, 104}, // black, opens above, closes into body
	)
	out, err := CDLDARKCLOUDCOVER(open, high, low, close)
	assertNoError(t, err, "CDLDARKCLOUDCOVER")
	assertLen(t, len(out), 2, "CDLDARKCLOUDCOVER")
	assertSignalRange(t, out, "CDLDARKCLOUDCOVER")
}

func TestCDLPIERCING(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{108, 110, 98, 100}, // long black
		[4]float64{97, 105, 95, 104},  // white, opens below, closes into body (>50%)
	)
	out, err := CDLPIERCING(open, high, low, close)
	assertNoError(t, err, "CDLPIERCING")
	assertLen(t, len(out), 2, "CDLPIERCING")
	assertSignalRange(t, out, "CDLPIERCING")
}

func TestCDLKICKING(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 100, 90, 90}, // black marubozu
		[4]float64{85, 88, 85, 88},   // white marubozu, gap down (bullish kicking)
	)
	out, err := CDLKICKING(open, high, low, close)
	assertNoError(t, err, "CDLKICKING")
	assertLen(t, len(out), 2, "CDLKICKING")
	assertSignalRange(t, out, "CDLKICKING")
}

func TestCDLCOUNTERATTACK(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 100}, // white (bearish counterattack: white→black at same close)
		[4]float64{100, 102, 95, 100}, // black, same close
	)
	out, err := CDLCOUNTERATTACK(open, high, low, close)
	assertNoError(t, err, "CDLCOUNTERATTACK")
	assertLen(t, len(out), 2, "CDLCOUNTERATTACK")
	assertSignalRange(t, out, "CDLCOUNTERATTACK")
}

func TestCDLBELTHOLD(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102}, // downtrend
		[4]float64{102, 103, 98, 100},  // downtrend
		[4]float64{100, 101, 96, 98},   // downtrend
		[4]float64{98, 105, 98, 105},   // bullish belt hold
	)
	out, err := CDLBELTHOLD(open, high, low, close)
	assertNoError(t, err, "CDLBELTHOLD")
	assertLen(t, len(out), 4, "CDLBELTHOLD")
	assertSignalRange(t, out, "CDLBELTHOLD")
}

func TestCDLHOMINGPIGEON(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102},       // downtrend
		[4]float64{102, 105, 100, 101},       // black
		[4]float64{101.5, 102, 100.5, 101.2}, // black inside prev
	)
	out, err := CDLHOMINGPIGEON(open, high, low, close)
	assertNoError(t, err, "CDLHOMINGPIGEON")
	assertLen(t, len(out), 3, "CDLHOMINGPIGEON")
	assertSignalRange(t, out, "CDLHOMINGPIGEON")
}

func TestCDLMATCHINGLOW(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{102, 103, 98, 99}, // black
		[4]float64{100, 103, 98, 99}, // black, same close
	)
	out, err := CDLMATCHINGLOW(open, high, low, close)
	assertNoError(t, err, "CDLMATCHINGLOW")
	assertLen(t, len(out), 2, "CDLMATCHINGLOW")
	assertSignalRange(t, out, "CDLMATCHINGLOW")
}

func TestCDLSEPARATINGLINES(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 104}, // white
		[4]float64{100, 108, 99, 106}, // white, same open
	)
	out, err := CDLSEPARATINGLINES(open, high, low, close)
	assertNoError(t, err, "CDLSEPARATINGLINES")
	assertLen(t, len(out), 2, "CDLSEPARATINGLINES")
	assertSignalRange(t, out, "CDLSEPARATINGLINES")
}

func TestCDLSTICKSANDWICH(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 104},  // white
		[4]float64{104, 105, 101, 102}, // black
		[4]float64{101, 106, 100, 104}, // white, same close as first
	)
	out, err := CDLSTICKSANDWICH(open, high, low, close)
	assertNoError(t, err, "CDLSTICKSANDWICH")
	assertLen(t, len(out), 3, "CDLSTICKSANDWICH")
	assertSignalRange(t, out, "CDLSTICKSANDWICH")
}

func TestCDLONNECK(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{102, 103, 98, 99}, // black
		[4]float64{96, 100, 95, 99},  // white, opens below, closes at prev close
	)
	out, err := CDLONNECK(open, high, low, close)
	assertNoError(t, err, "CDLONNECK")
	assertLen(t, len(out), 2, "CDLONNECK")
	assertSignalRange(t, out, "CDLONNECK")
}

func TestCDLINNECK(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{102, 103, 98, 99},   // black
		[4]float64{96, 99.2, 95, 98.1}, // white, opens below, closes near prev low
	)
	out, err := CDLINNECK(open, high, low, close)
	assertNoError(t, err, "CDLINNECK")
	assertLen(t, len(out), 2, "CDLINNECK")
	assertSignalRange(t, out, "CDLINNECK")
}

func TestCDLTHRUSTING(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{102, 103, 98, 99},  // black
		[4]float64{96, 101, 95, 99.5}, // white, opens below, closes < midpoint
	)
	out, err := CDLTHRUSTING(open, high, low, close)
	assertNoError(t, err, "CDLTHRUSTING")
	assertLen(t, len(out), 2, "CDLTHRUSTING")
	assertSignalRange(t, out, "CDLTHRUSTING")
}

func TestCDLHARAMICROSS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 95, 98}, // large black
		[4]float64{97, 99, 96, 97},   // doji inside (bullish harami cross)
	)
	out, err := CDLHARAMICROSS(open, high, low, close)
	assertNoError(t, err, "CDLHARAMICROSS")
	assertLen(t, len(out), 2, "CDLHARAMICROSS")
	assertSignalRange(t, out, "CDLHARAMICROSS")
}

func TestCDLGAPSIDESIDEWHITE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 104},  // white
		[4]float64{108, 112, 107, 110}, // white, gap up
		[4]float64{109, 111, 108, 110}, // white, same close (up-gap side-by-side)
	)
	out, err := CDLGAPSIDESIDEWHITE(open, high, low, close)
	assertNoError(t, err, "CDLGAPSIDESIDEWHITE")
	assertLen(t, len(out), 3, "CDLGAPSIDESIDEWHITE")
	assertSignalRange(t, out, "CDLGAPSIDESIDEWHITE")
}

func TestCDLKICKINGBYLENGTH(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 100, 90, 90}, // black marubozu
		[4]float64{85, 88, 85, 88},   // white marubozu, gap down
	)
	out, err := CDLKICKINGBYLENGTH(open, high, low, close)
	assertNoError(t, err, "CDLKICKINGBYLENGTH")
	assertLen(t, len(out), 2, "CDLKICKINGBYLENGTH")
	assertSignalRange(t, out, "CDLKICKINGBYLENGTH")
}

func TestCDLDOJISTAR(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 98, 102},  // long white
		[4]float64{105, 106, 104, 105}, // doji, gap up (bearish doji star)
	)
	out, err := CDLDOJISTAR(open, high, low, close)
	assertNoError(t, err, "CDLDOJISTAR")
	assertLen(t, len(out), 2, "CDLDOJISTAR")
	assertSignalRange(t, out, "CDLDOJISTAR")
}

func TestCDLTASUKIGAP(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 105, 98, 104},  // white
		[4]float64{108, 112, 107, 110}, // white, gap up
		[4]float64{110, 111, 106, 107}, // black, closes into gap (bullish tasuki gap)
	)
	out, err := CDLTASUKIGAP(open, high, low, close)
	assertNoError(t, err, "CDLTASUKIGAP")
	assertLen(t, len(out), 3, "CDLTASUKIGAP")
	assertSignalRange(t, out, "CDLTASUKIGAP")
}

// =====================================================================
// 星形形态测试
// =====================================================================

func TestCDLMORNINGSTAR(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 98, 100}, // long black
		[4]float64{96, 97, 94, 95},    // small candle, gap down
		[4]float64{97, 104, 96, 103},  // long white, closes into black
	)
	out, err := CDLMORNINGSTAR(open, high, low, close)
	assertNoError(t, err, "CDLMORNINGSTAR")
	assertLen(t, len(out), 3, "CDLMORNINGSTAR")
	assertSignalRange(t, out, "CDLMORNINGSTAR")
}

func TestCDLEVENINGSTAR(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 107, 98, 105},  // long white
		[4]float64{107, 108, 105, 106}, // small candle, gap up
		[4]float64{108, 109, 100, 102}, // long black, closes into white
	)
	out, err := CDLEVENINGSTAR(open, high, low, close)
	assertNoError(t, err, "CDLEVENINGSTAR")
	assertLen(t, len(out), 3, "CDLEVENINGSTAR")
	assertSignalRange(t, out, "CDLEVENINGSTAR")
}

func TestCDLMORNINGDOJISTAR(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 98, 100}, // long black
		[4]float64{96, 97, 95, 96},    // doji, gap down
		[4]float64{97, 104, 96, 103},  // long white, closes into black
	)
	out, err := CDLMORNINGDOJISTAR(open, high, low, close)
	assertNoError(t, err, "CDLMORNINGDOJISTAR")
	assertLen(t, len(out), 3, "CDLMORNINGDOJISTAR")
	assertSignalRange(t, out, "CDLMORNINGDOJISTAR")
}

func TestCDLEVENINGDOJISTAR(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 107, 98, 105},  // long white
		[4]float64{107, 108, 106, 107}, // doji, gap up
		[4]float64{108, 109, 100, 102}, // long black, closes into white
	)
	out, err := CDLEVENINGDOJISTAR(open, high, low, close)
	assertNoError(t, err, "CDLEVENINGDOJISTAR")
	assertLen(t, len(out), 3, "CDLEVENINGDOJISTAR")
	assertSignalRange(t, out, "CDLEVENINGDOJISTAR")
}

func TestCDLABANDONEDBABY(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 98, 100}, // long black
		[4]float64{95, 96, 94, 95},    // doji, gap down with shadows not overlapping
		[4]float64{97, 104, 97, 103},  // long white, gap up from doji
	)
	out, err := CDLABANDONEDBABY(open, high, low, close)
	assertNoError(t, err, "CDLABANDONEDBABY")
	assertLen(t, len(out), 3, "CDLABANDONEDBABY")
	assertSignalRange(t, out, "CDLABANDONEDBABY")
}

func TestCDLTRISTAR(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 102, 98, 100}, // doji 1
		[4]float64{95, 96, 94, 95},    // doji 2, gap down
		[4]float64{97, 98, 96, 97},    // doji 3, gap up
	)
	out, err := CDLTRISTAR(open, high, low, close)
	assertNoError(t, err, "CDLTRISTAR")
	assertLen(t, len(out), 3, "CDLTRISTAR")
	assertSignalRange(t, out, "CDLTRISTAR")
}

// =====================================================================
// 三K线形态测试
// =====================================================================

func TestCDL2CROWS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 98, 106},  // long white
		[4]float64{108, 110, 105, 106}, // black, gaps up
		[4]float64{108, 109, 102, 103}, // black, opens inside, closes below
	)
	out, err := CDL2CROWS(open, high, low, close)
	assertNoError(t, err, "CDL2CROWS")
	assertLen(t, len(out), 3, "CDL2CROWS")
	assertSignalRange(t, out, "CDL2CROWS")
}

func TestCDL3BLACKCROWS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{98, 102, 96, 100},   // uptrend
		[4]float64{105, 106, 100, 102}, // white
		[4]float64{101, 102, 97, 98},   // long black
		[4]float64{98.5, 99, 94, 95},   // long black
		[4]float64{95.5, 96, 91, 92},   // long black
	)
	out, err := CDL3BLACKCROWS(open, high, low, close)
	assertNoError(t, err, "CDL3BLACKCROWS")
	assertLen(t, len(out), 5, "CDL3BLACKCROWS")
	assertSignalRange(t, out, "CDL3BLACKCROWS")
}

func TestCDL3INSIDE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 98, 100},  // black
		[4]float64{99, 102, 98.5, 101}, // white inside
		[4]float64{102, 105, 100, 104}, // white confirm above first
	)
	out, err := CDL3INSIDE(open, high, low, close)
	assertNoError(t, err, "CDL3INSIDE")
	assertLen(t, len(out), 3, "CDL3INSIDE")
	assertSignalRange(t, out, "CDL3INSIDE")
}

func TestCDL3OUTSIDE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{101, 103, 99, 100},  // small black
		[4]float64{98, 105, 97, 104},   // white engulfs
		[4]float64{105, 108, 103, 107}, // white confirm
	)
	out, err := CDL3OUTSIDE(open, high, low, close)
	assertNoError(t, err, "CDL3OUTSIDE")
	assertLen(t, len(out), 3, "CDL3OUTSIDE")
	assertSignalRange(t, out, "CDL3OUTSIDE")
}

func TestCDL3WHITESOLDIERS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102}, // downtrend
		[4]float64{102, 103, 98, 100},  // black
		[4]float64{98, 105, 97, 103},   // long white
		[4]float64{102, 108, 101, 106}, // long white
		[4]float64{105, 112, 104, 110}, // long white
	)
	out, err := CDL3WHITESOLDIERS(open, high, low, close)
	assertNoError(t, err, "CDL3WHITESOLDIERS")
	assertLen(t, len(out), 5, "CDL3WHITESOLDIERS")
	assertSignalRange(t, out, "CDL3WHITESOLDIERS")
}

func TestCDLIDENTICAL3CROWS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{104, 106, 100, 102}, // uptrend context
		[4]float64{103, 103, 98, 99},   // first black
		[4]float64{99.5, 99.5, 95, 96}, // second black
		[4]float64{96.5, 96.5, 92, 93}, // third black
	)
	out, err := CDLIDENTICAL3CROWS(open, high, low, close)
	assertNoError(t, err, "CDLIDENTICAL3CROWS")
	assertLen(t, len(out), 4, "CDLIDENTICAL3CROWS")
	assertSignalRange(t, out, "CDLIDENTICAL3CROWS")
}

func TestCDLADVANCEBLOCK(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 99, 106},  // long white
		[4]float64{106, 110, 105, 108}, // smaller white, bigger upper shadow
		[4]float64{108, 112, 107, 109}, // even smaller white
	)
	out, err := CDLADVANCEBLOCK(open, high, low, close)
	assertNoError(t, err, "CDLADVANCEBLOCK")
	assertLen(t, len(out), 3, "CDLADVANCEBLOCK")
	assertSignalRange(t, out, "CDLADVANCEBLOCK")
}

func TestCDLSTALLEDPATTERN(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 99, 106},  // long white
		[4]float64{106, 110, 105, 108}, // white
		[4]float64{108, 110, 107, 109}, // tiny white near prev close
	)
	out, err := CDLSTALLEDPATTERN(open, high, low, close)
	assertNoError(t, err, "CDLSTALLEDPATTERN")
	assertLen(t, len(out), 3, "CDLSTALLEDPATTERN")
	assertSignalRange(t, out, "CDLSTALLEDPATTERN")
}

func TestCDLUNIQUE3RIVER(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 98, 100}, // long black
		[4]float64{100, 101, 96, 97},  // black inside
		[4]float64{95, 98, 93, 96.5},  // small white, new low
	)
	out, err := CDLUNIQUE3RIVER(open, high, low, close)
	assertNoError(t, err, "CDLUNIQUE3RIVER")
	assertLen(t, len(out), 3, "CDLUNIQUE3RIVER")
	assertSignalRange(t, out, "CDLUNIQUE3RIVER")
}

func TestCDLUPSIDEGAP2CROWS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 98, 106},  // long white in uptrend
		[4]float64{110, 112, 108, 109}, // black, gap up
		[4]float64{111, 112, 105, 106}, // black, opens above, closes below
	)
	out, err := CDLUPSIDEGAP2CROWS(open, high, low, close)
	assertNoError(t, err, "CDLUPSIDEGAP2CROWS")
	assertLen(t, len(out), 3, "CDLUPSIDEGAP2CROWS")
	assertSignalRange(t, out, "CDLUPSIDEGAP2CROWS")
}

func TestCDL3LINESTRIKE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 100, 102}, // black 1
		[4]float64{102, 103, 98, 100},  // black 2
		[4]float64{100, 101, 96, 98},   // black 3
		[4]float64{97, 108, 96, 107},   // long white, closes above first open
	)
	out, err := CDL3LINESTRIKE(open, high, low, close)
	assertNoError(t, err, "CDL3LINESTRIKE")
	assertLen(t, len(out), 4, "CDL3LINESTRIKE")
	assertSignalRange(t, out, "CDL3LINESTRIKE")
}

func TestCDLHIKKAKE(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 90, 105}, // parent candle
		[4]float64{103, 107, 93, 104}, // inside bar
		[4]float64{100, 102, 88, 91},  // false breakdown (low < parent low)
		[4]float64{92, 98, 91, 97},    // confirmation (close higher)
	)
	out, err := CDLHIKKAKE(open, high, low, close)
	assertNoError(t, err, "CDLHIKKAKE")
	assertLen(t, len(out), 4, "CDLHIKKAKE")
	assertSignalRange(t, out, "CDLHIKKAKE")
}

func TestCDLHIKKAKEMOD(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 90, 105}, // parent candle
		[4]float64{103, 107, 93, 104}, // inside bar
		[4]float64{100, 102, 88, 91},  // false breakdown
		[4]float64{92, 98, 91, 97},    // confirm
	)
	out, err := CDLHIKKAKEMOD(open, high, low, close)
	assertNoError(t, err, "CDLHIKKAKEMOD")
	assertLen(t, len(out), 4, "CDLHIKKAKEMOD")
	assertSignalRange(t, out, "CDLHIKKAKEMOD")
}

func TestCDLINSIDE(t *testing.T) {
	// CDL3INSIDE 的别名 - 验证相同行为
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 98, 100},  // 阴线
		[4]float64{99, 102, 98.5, 101}, // 内部阳线
		[4]float64{102, 105, 100, 104}, // 阳线确认
	)
	out1, _ := CDL3INSIDE(open, high, low, close)
	out2, _ := CDLINSIDE(open, high, low, close)
	for i := range out1 {
		if out1[i] != out2[i] {
			t.Errorf("CDLINSIDE: mismatch with CDL3INSIDE at %d: %d vs %d", i, out1[i], out2[i])
		}
	}
}

func TestCDL3STARSINSOUTH(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 107, 100, 101}, // downtrend context
		[4]float64{103, 103, 98, 99},   // first black, long body/long shadow
		[4]float64{99.5, 99.5, 97, 98}, // second black, smaller
		[4]float64{98, 98.2, 96.5, 97}, // third black, even smaller
	)
	out, err := CDL3STARSINSOUTH(open, high, low, close)
	assertNoError(t, err, "CDL3STARSINSOUTH")
	assertLen(t, len(out), 4, "CDL3STARSINSOUTH")
	assertSignalRange(t, out, "CDL3STARSINSOUTH")
}

// =====================================================================
// 复杂形态测试
// =====================================================================

func TestCDLBREAKAWAY(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 108, 98, 100}, // long black
		[4]float64{100, 101, 96, 98},  // lower close
		[4]float64{98, 99, 94, 96},    // lower close
		[4]float64{96, 97, 92, 94},    // lower close
		[4]float64{95, 105, 94, 103},  // long white, above prev highs
	)
	out, err := CDLBREAKAWAY(open, high, low, close)
	assertNoError(t, err, "CDLBREAKAWAY")
	assertLen(t, len(out), 5, "CDLBREAKAWAY")
	assertSignalRange(t, out, "CDLBREAKAWAY")
}

func TestCDLMATHOLD(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 98, 108},  // long white
		[4]float64{108, 109, 105, 106}, // small black within
		[4]float64{106, 107, 103, 104}, // small black within
		[4]float64{104, 105, 101, 102}, // small black within
		[4]float64{103, 115, 102, 113}, // long white, above first
	)
	out, err := CDLMATHOLD(open, high, low, close)
	assertNoError(t, err, "CDLMATHOLD")
	assertLen(t, len(out), 5, "CDLMATHOLD")
	assertSignalRange(t, out, "CDLMATHOLD")
}

func TestCDLLADDERBOTTOM(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{105, 106, 100, 102}, // black 1
		[4]float64{102, 103, 98, 100},  // black 2
		[4]float64{100, 101, 96, 98},   // black 3
		[4]float64{98, 100, 94, 96},    // black 4 with upper shadow
		[4]float64{97, 103, 97, 102},   // white, gap up
	)
	out, err := CDLLADDERBOTTOM(open, high, low, close)
	assertNoError(t, err, "CDLLADDERBOTTOM")
	assertLen(t, len(out), 5, "CDLLADDERBOTTOM")
	assertSignalRange(t, out, "CDLLADDERBOTTOM")
}

func TestCDLRISEFALL3METHODS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 110, 98, 108},  // long white
		[4]float64{108, 109, 105, 106}, // small black
		[4]float64{106, 107, 104, 105}, // small black
		[4]float64{105, 106, 103, 104}, // small black
		[4]float64{104, 115, 103, 113}, // long white new high
	)
	out, err := CDLRISEFALL3METHODS(open, high, low, close)
	assertNoError(t, err, "CDLRISEFALL3METHODS")
	assertLen(t, len(out), 5, "CDLRISEFALL3METHODS")
	assertSignalRange(t, out, "CDLRISEFALL3METHODS")
}

func TestCDLXSIDEGAP3METHODS(t *testing.T) {
	open, high, low, close := makeTestOHLC(
		[4]float64{100, 108, 98, 106},  // long white
		[4]float64{110, 114, 109, 112}, // white, gap up
		[4]float64{112, 113, 108, 109}, // black filling gap partially
	)
	out, err := CDLXSIDEGAP3METHODS(open, high, low, close)
	assertNoError(t, err, "CDLXSIDEGAP3METHODS")
	assertLen(t, len(out), 3, "CDLXSIDEGAP3METHODS")
	assertSignalRange(t, out, "CDLXSIDEGAP3METHODS")
}

// =====================================================================
// 回溯期测试
// =====================================================================

func TestAllLookbacks(t *testing.T) {
	tests := []struct {
		name string
		fn   func() int
		min  int
	}{
		{"CDLDOJI", CDLDOJILookback, 1},
		{"CDLDRAGONFLYDOJI", CDLDRAGONFLYDOJILookback, 1},
		{"CDLGRAVESTONEDOJI", CDLGRAVESTONEDOJILookback, 1},
		{"CDLMARUBOZU", CDLMARUBOZULookback, 1},
		{"CDLHAMMER", CDLHAMMERLookback, 4},
		{"CDLHANGINGMAN", CDLHANGINGMANLookback, 4},
		{"CDLSHOOTINGSTAR", CDLSHOOTINGSTARLookback, 4},
		{"CDLENGULFING", CDLENGULFINGLookback, 2},
		{"CDLHARAMI", CDLHARAMILookback, 2},
		{"CDLPIERCING", CDLPIERCINGLookback, 2},
		{"CDLDARKCLOUDCOVER", CDLDARKCLOUDCOVERLookback, 2},
		{"CDLMORNINGSTAR", CDLMORNINGSTARLookback, 3},
		{"CDLEVENINGSTAR", CDLEVENINGSTARLookback, 3},
		{"CDL3WHITESOLDIERS", CDL3WHITESOLDIERSLookback, 4},
		{"CDL3BLACKCROWS", CDL3BLACKCROWSLookback, 4},
		{"CDL3INSIDE", CDL3INSIDELookback, 3},
		{"CDL3OUTSIDE", CDL3OUTSIDELookback, 3},
		{"CDLBREAKAWAY", CDLBREAKAWAYLookback, 5},
		{"CDLMATHOLD", CDLMATHOLDLookback, 5},
		{"CDLRISEFALL3METHODS", CDLRISEFALL3METHODSLookback, 5},
		{"CDLINSIDE", CDLINSIDELookback, 3},
	}

	for _, tt := range tests {
		lb := tt.fn()
		if lb < tt.min {
			t.Errorf("%s Lookback = %d, want >= %d", tt.name, lb, tt.min)
		}
		if lb > 10 {
			t.Errorf("%s Lookback = %d, unreasonably large", tt.name, lb)
		}
	}
}

// =====================================================================
// 边界情况测试
// =====================================================================

func TestAllZeroes(t *testing.T) {
	// 所有蜡烛未变化（零实体，无范围）
	open := []float64{100, 100, 100}
	high := []float64{100, 100, 100}
	low := []float64{100, 100, 100}
	close := []float64{100, 100, 100}
	// 不应引发 panic
	out, err := CDLDOJI(open, high, low, close)
	assertNoError(t, err, "CDLDOJI zeroes")
	assertLen(t, len(out), 3, "CDLDOJI zeroes")
	assertSignalRange(t, out, "CDLDOJI zeroes")
}

func TestLargeNumbers(t *testing.T) {
	// 使用非常大的价格值进行测试（无溢出）
	open := []float64{1e6, 1.01e6, 1.02e6}
	high := []float64{1.05e6, 1.06e6, 1.07e6}
	low := []float64{0.99e6, 1.0e6, 1.01e6}
	close := []float64{1.03e6, 1.04e6, 1.05e6}
	out, err := CDLSPINNINGTOP(open, high, low, close)
	assertNoError(t, err, "large numbers")
	assertLen(t, len(out), 3, "large numbers")
	assertSignalRange(t, out, "large numbers")
}

func TestAllSignalsValid(t *testing.T) {
	// 综合信号有效性测试
	tests := []struct {
		name string
		fn   func(open, high, low, close []float64) ([]int, error)
	}{
		{"CDLDOJI", CDLDOJI},
		{"CDLDRAGONFLYDOJI", CDLDRAGONFLYDOJI},
		{"CDLGRAVESTONEDOJI", CDLGRAVESTONEDOJI},
		{"CDLLONGLEGGEDDOJI", CDLLONGLEGGEDDOJI},
		{"CDLMARUBOZU", CDLMARUBOZU},
		{"CDLCLOSINGMARUBOZU", CDLCLOSINGMARUBOZU},
		{"CDLLONGLINE", CDLLONGLINE},
		{"CDLSHORTLINE", CDLSHORTLINE},
		{"CDLSPINNINGTOP", CDLSPINNINGTOP},
		{"CDLHIGHWAVE", CDLHIGHWAVE},
		{"CDLRICKSHAWMAN", CDLRICKSHAWMAN},
		{"CDLHAMMER", CDLHAMMER},
		{"CDLHANGINGMAN", CDLHANGINGMAN},
		{"CDLINVERTEDHAMMER", CDLINVERTEDHAMMER},
		{"CDLSHOOTINGSTAR", CDLSHOOTINGSTAR},
		{"CDLTAKURI", CDLTAKURI},
		{"CDLENGULFING", CDLENGULFING},
		{"CDLHARAMI", CDLHARAMI},
		{"CDLHARAMICROSS", CDLHARAMICROSS},
		{"CDLDARKCLOUDCOVER", CDLDARKCLOUDCOVER},
		{"CDLPIERCING", CDLPIERCING},
		{"CDLCOUNTERATTACK", CDLCOUNTERATTACK},
		{"CDLKICKING", CDLKICKING},
		{"CDLKICKINGBYLENGTH", CDLKICKINGBYLENGTH},
		{"CDLGAPSIDESIDEWHITE", CDLGAPSIDESIDEWHITE},
		{"CDLDOJISTAR", CDLDOJISTAR},
		{"CDLBELTHOLD", CDLBELTHOLD},
		{"CDLHOMINGPIGEON", CDLHOMINGPIGEON},
		{"CDLMATCHINGLOW", CDLMATCHINGLOW},
		{"CDLSEPARATINGLINES", CDLSEPARATINGLINES},
		{"CDLSTICKSANDWICH", CDLSTICKSANDWICH},
		{"CDLTASUKIGAP", CDLTASUKIGAP},
		{"CDLONNECK", CDLONNECK},
		{"CDLINNECK", CDLINNECK},
		{"CDLTHRUSTING", CDLTHRUSTING},
		{"CDLMORNINGSTAR", CDLMORNINGSTAR},
		{"CDLEVENINGSTAR", CDLEVENINGSTAR},
		{"CDLMORNINGDOJISTAR", CDLMORNINGDOJISTAR},
		{"CDLEVENINGDOJISTAR", CDLEVENINGDOJISTAR},
		{"CDLABANDONEDBABY", CDLABANDONEDBABY},
		{"CDLTRISTAR", CDLTRISTAR},
		{"CDL2CROWS", CDL2CROWS},
		{"CDL3BLACKCROWS", CDL3BLACKCROWS},
		{"CDL3INSIDE", CDL3INSIDE},
		{"CDL3OUTSIDE", CDL3OUTSIDE},
		{"CDL3WHITESOLDIERS", CDL3WHITESOLDIERS},
		{"CDL3STARSINSOUTH", CDL3STARSINSOUTH},
		{"CDLIDENTICAL3CROWS", CDLIDENTICAL3CROWS},
		{"CDLADVANCEBLOCK", CDLADVANCEBLOCK},
		{"CDLSTALLEDPATTERN", CDLSTALLEDPATTERN},
		{"CDLUNIQUE3RIVER", CDLUNIQUE3RIVER},
		{"CDLUPSIDEGAP2CROWS", CDLUPSIDEGAP2CROWS},
		{"CDLHIKKAKE", CDLHIKKAKE},
		{"CDLHIKKAKEMOD", CDLHIKKAKEMOD},
		{"CDL3LINESTRIKE", CDL3LINESTRIKE},
		{"CDLINSIDE", CDLINSIDE},
		{"CDLBREAKAWAY", CDLBREAKAWAY},
		{"CDLMATHOLD", CDLMATHOLD},
		{"CDLLADDERBOTTOM", CDLLADDERBOTTOM},
		{"CDLRISEFALL3METHODS", CDLRISEFALL3METHODS},
		{"CDLXSIDEGAP3METHODS", CDLXSIDEGAP3METHODS},
	}

	// 生成 20 根随机蜡烛
	n := 20
	open := make([]float64, n)
	high := make([]float64, n)
	low := make([]float64, n)
	close := make([]float64, n)
	prices := []float64{100, 102, 99, 105, 98, 103, 97, 106, 95, 104,
		96, 107, 94, 108, 93, 109, 92, 110, 91, 111}
	for i := 0; i < n; i++ {
		base := prices[i]
		open[i] = base
		high[i] = base + float64(i%5)
		low[i] = base - float64(i%4)
		close[i] = base + float64((i%3)-1)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.fn(open, high, low, close)
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.name, err)
				return
			}
			if len(out) != n {
				t.Errorf("%s: output length = %d, want %d", tt.name, len(out), n)
			}
			assertSignalRange(t, out, tt.name)
		})
	}
}
