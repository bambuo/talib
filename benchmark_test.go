package talib

import (
	"math"
	"testing"
)

// =============================================================================
// 辅助函数：生成合成测试数据
// =============================================================================

// generateBenchData 创建用于 benchmark 的合成价格序列。
func generateBenchData(n int) []float64 {
	data := make([]float64, n)
	price := 100.0
	for i := 0; i < n; i++ {
		price += (math.Sin(float64(i)*0.1) + math.Cos(float64(i)*0.05)) * 0.5
		data[i] = price
	}
	return data
}

// generateOHLCData 创建用于 benchmark 的合成 OHLC 数据。
func generateOHLCData(n int) (high, low, close []float64) {
	high = make([]float64, n)
	low = make([]float64, n)
	close = make([]float64, n)
	price := 100.0
	for i := 0; i < n; i++ {
		price += (math.Sin(float64(i)*0.1) + math.Cos(float64(i)*0.05)) * 0.5
		close[i] = price
		high[i] = price + math.Abs(math.Sin(float64(i)*0.3))*2
		low[i] = price - math.Abs(math.Cos(float64(i)*0.3))*2
	}
	return
}

// benchmark 数据集大小常量
const (
	benchSmallSize = 500   // 实时 K 线处理规模
	benchLargeSize = 10000 // 回测处理规模
)

// =============================================================================
// P0：缠论核心依赖 — 每根 K 线实时调用
// =============================================================================

// --- SMA ---

func BenchmarkSMA_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SMA(data, 14)
	}
}

func BenchmarkSMA_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SMA(data, 14)
	}
}

// --- EMA ---

func BenchmarkEMA_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EMA(data, 14)
	}
}

func BenchmarkEMA_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = EMA(data, 14)
	}
}

// --- MACD ---

func BenchmarkMACD_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MACD(data, 12, 26, 9)
	}
}

func BenchmarkMACD_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MACD(data, 12, 26, 9)
	}
}

// --- RSI ---

func BenchmarkRSI_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RSI(data, 14)
	}
}

func BenchmarkRSI_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RSI(data, 14)
	}
}

// --- BBands (SMA) ---

func BenchmarkBBands_SMA_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BBands(data, 20, 2, 2, MASMA)
	}
}

func BenchmarkBBands_SMA_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BBands(data, 20, 2, 2, MASMA)
	}
}

// BBands (EMA)
func BenchmarkBBands_EMA_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BBands(data, 20, 2, 2, MAEMA)
	}
}

func BenchmarkBBands_EMA_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = BBands(data, 20, 2, 2, MAEMA)
	}
}

// =============================================================================
// P1：计算复杂 — 潜在线性能瓶颈
// =============================================================================

// --- WMA ---

func BenchmarkWMA_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = WMA(data, 14)
	}
}

func BenchmarkWMA_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = WMA(data, 14)
	}
}

// --- ADX ---

func BenchmarkADX_Small(b *testing.B) {
	high, low, close := generateOHLCData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ADX(high, low, close, 14)
	}
}

func BenchmarkADX_Large(b *testing.B) {
	high, low, close := generateOHLCData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ADX(high, low, close, 14)
	}
}

// --- ATR ---

func BenchmarkATR_Small(b *testing.B) {
	high, low, close := generateOHLCData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ATR(high, low, close, 14)
	}
}

func BenchmarkATR_Large(b *testing.B) {
	high, low, close := generateOHLCData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ATR(high, low, close, 14)
	}
}

// --- Stochastic ---

func BenchmarkSTOCH_Small(b *testing.B) {
	high, low, close := generateOHLCData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = STOCH(high, low, close, 5, 3, 3, MASMA)
	}
}

func BenchmarkSTOCH_Large(b *testing.B) {
	high, low, close := generateOHLCData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = STOCH(high, low, close, 5, 3, 3, MASMA)
	}
}

// --- KAMA ---

func BenchmarkKAMA_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = KAMA(data, 10)
	}
}

func BenchmarkKAMA_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = KAMA(data, 10)
	}
}

// --- HT_TRENDLINE ---

func BenchmarkHT_TRENDLINE_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_TRENDLINE(data)
	}
}

func BenchmarkHT_TRENDLINE_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_TRENDLINE(data)
	}
}

// --- HT_DCPERIOD ---

func BenchmarkHT_DCPERIOD_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_DCPERIOD(data)
	}
}

func BenchmarkHT_DCPERIOD_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_DCPERIOD(data)
	}
}

// --- HT_DCPHASE ---

func BenchmarkHT_DCPHASE_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_DCPHASE(data)
	}
}

func BenchmarkHT_DCPHASE_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_DCPHASE(data)
	}
}

// --- HT_PHASOR ---

func BenchmarkHT_PHASOR_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_PHASOR(data)
	}
}

func BenchmarkHT_PHASOR_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_PHASOR(data)
	}
}

// --- HT_SINE ---

func BenchmarkHT_SINE_Small(b *testing.B) {
	data := generateBenchData(benchSmallSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_SINE(data)
	}
}

func BenchmarkHT_SINE_Large(b *testing.B) {
	data := generateBenchData(benchLargeSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HT_SINE(data)
	}
}
