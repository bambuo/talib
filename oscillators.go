package talib

import (
	"math"
)

// APO 计算绝对价格振荡器。
//
// APO = EMA(输入, 快速周期) - EMA(输入, 慢速周期)
func APO(input []float64, fastPeriod, slowPeriod int, maType MAType) ([]float64, error) {
	if err := ValidateNumericInput(input, "APO:input"); err != nil {
		return nil, err
	}
	if fastPeriod >= slowPeriod {
		return nil, ErrFastPeriodGreaterOrEqualSlow
	}

	emaFast, err := MA(input, fastPeriod, maType)
	if err != nil {
		return nil, err
	}
	emaSlow, err := MA(input, slowPeriod, maType)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)
	for i := slowPeriod - 1; i < n; i++ {
		out[i] = emaFast[i] - emaSlow[i]
	}

	return out, nil
}

// APOLookback 返回 APO 的 lookback。
func APOLookback(fastPeriod, slowPeriod int) int {
	return slowPeriod - 1
}

// PPO 计算百分比价格振荡器。
//
// PPO = (EMA(快速) - EMA(慢速)) / EMA(慢速) * 100
func PPO(input []float64, fastPeriod, slowPeriod int, maType MAType) ([]float64, error) {
	if err := ValidateNumericInput(input, "PPO:input"); err != nil {
		return nil, err
	}
	if fastPeriod >= slowPeriod {
		return nil, ErrFastPeriodGreaterOrEqualSlow
	}

	emaFast, err := MA(input, fastPeriod, maType)
	if err != nil {
		return nil, err
	}
	emaSlow, err := MA(input, slowPeriod, maType)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)
	for i := slowPeriod - 1; i < n; i++ {
		if emaSlow[i] != 0 {
			out[i] = (emaFast[i] - emaSlow[i]) / emaSlow[i] * 100.0
		}
	}

	return out, nil
}

// PPOLookback 返回 PPO 的 lookback。
func PPOLookback(fastPeriod, slowPeriod int) int {
	return slowPeriod - 1
}

// MACDEXT 计算具有可控制 MA 类型的 MACD。
//
// 与 MACD 相同，但允许为每个组件选择 MA 类型。
func MACDEXT(input []float64, fastPeriod int, fastMAType MAType,
	slowPeriod int, slowMAType MAType,
	signalPeriod int, signalMAType MAType) (*MACDResult, error) {

	if err := ValidateNumericInput(input, "MACDEXT:input"); err != nil {
		return nil, err
	}
	if fastPeriod >= slowPeriod {
		return nil, ErrFastPeriodGreaterOrEqualSlow
	}

	n := len(input)

	emaFast, err := MA(input, fastPeriod, fastMAType)
	if err != nil {
		return nil, err
	}
	emaSlow, err := MA(input, slowPeriod, slowMAType)
	if err != nil {
		return nil, err
	}

	macd := MakeOutput(n)
	for i := slowPeriod - 1; i < n; i++ {
		macd[i] = emaFast[i] - emaSlow[i]
	}

	signal, err := MA(macd, signalPeriod, signalMAType)
	if err != nil {
		return nil, err
	}

	hist := MakeOutput(n)
	startIdx := slowPeriod + signalPeriod - 2
	for i := startIdx; i < n; i++ {
		if !IsNaN(macd[i]) && !IsNaN(signal[i]) {
			hist[i] = macd[i] - signal[i]
		}
	}

	return &MACDResult{
		MACD:      macd,
		Signal:    signal,
		Histogram: hist,
	}, nil
}

// MACDEXTLookback 返回 MACDEXT 的 lookback。
func MACDEXTLookback(fastPeriod, slowPeriod, signalPeriod int) int {
	return slowPeriod + signalPeriod - 2
}

// MACDFIX 使用固定的 12/26/9 周期计算 MACD。
//
// MACD 线    = EMA(输入, 12) - EMA(输入, 26)
// 信号线     = EMA(MACD, 9)
// 柱状图     = MACD - 信号
func MACDFIX(input []float64, signalPeriod int) ([]float64, error) {
	result, err := MACD(input, 12, 26, signalPeriod)
	if err != nil {
		return nil, err
	}
	return result.Histogram, nil
}

// MACDFIXLookback 返回 MACDFIX 的 lookback。
func MACDFIXLookback(signalPeriod int) int {
	return 26 + signalPeriod - 2
}

// T3 计算 T3 移动平均。
//
// T3 是带体积因子的三重指数移动平均。
//
//	e1 = EMA(输入, 周期)
//	e2 = EMA(e1, 周期)
//	e3 = EMA(e2, 周期)
//	e4 = EMA(e3, 周期)
//	e5 = EMA(e4, 周期)
//	e6 = EMA(e5, 周期)
//	T3 = c1*e6 + c2*e5 + c3*e4
//
// 其中 c1 = -v^3, c2 = 3v^2+3v^3, c3 = -6v^2-3v-3v^3
// 且 v = volumeFactor（默认 0.7）
func T3(input []float64, period int, volumeFactor float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "T3:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "T3:period"); err != nil {
		return nil, err
	}

	e1, _ := EMA(input, period)
	e2, _ := EMA(e1, period)
	e3, _ := EMA(e2, period)
	e4, _ := EMA(e3, period)
	e5, _ := EMA(e4, period)
	e6, _ := EMA(e5, period)

	v := volumeFactor
	v2 := v * v
	v3 := v2 * v

	c1 := -v3
	c2 := 3*v2 + 3*v3
	c3 := -6*v2 - 3*v - 3*v3

	n := len(input)
	out := MakeOutput(n)
	startIdx := 6*period - 6
	for i := startIdx; i < n; i++ {
		out[i] = c1*e6[i] + c2*e5[i] + c3*e4[i]
	}

	return out, nil
}

// T3Lookback 返回 T3 的 lookback。
func T3Lookback(period int) int {
	return 6*period - 6
}

// TRIX 计算三重平滑 EMA 的 1 日变动率。
//
// EMA1 = EMA(输入, 周期)
// EMA2 = EMA(EMA1, 周期)
// EMA3 = EMA(EMA2, 周期)
// TRIX = (EMA3[i] - EMA3[i-1]) / EMA3[i-1] * 100
func TRIX(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "TRIX:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "TRIX:period"); err != nil {
		return nil, err
	}

	e1, _ := EMA(input, period)
	e2, _ := EMA(e1, period)
	e3, _ := EMA(e2, period)

	n := len(input)
	out := MakeOutput(n)
	startIdx := 3*period - 3 + 1 // 3 EMAs + 1 for ROC
	for i := startIdx; i < n; i++ {
		if !math.IsNaN(e3[i-1]) && e3[i-1] != 0 {
			out[i] = (e3[i] - e3[i-1]) / e3[i-1] * 100.0
		}
	}

	return out, nil
}

// TRIXLookback 返回 TRIX 的 lookback。
func TRIXLookback(period int) int {
	return 3*period - 2
}
