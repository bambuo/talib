package talib

// MACDResult 包含 MACD 线、信号线和柱状图。
type MACDResult struct {
	MACD      []float64 // MACD 线：EMA(快速) - EMA(慢速)
	Signal    []float64 // 信号线：EMA(MACD, 信号周期)
	Histogram []float64 // 柱状图：MACD - 信号
}

// MACD 计算移动平均收敛/发散指标。
//
// MACD 线    = EMA(输入, 快速周期) - EMA(输入, 慢速周期)
// 信号线     = EMA(MACD 线, 信号周期)
// 柱状图     = MACD 线 - 信号线
//
// 前 (慢速周期+信号周期-2) 个元素为 NaN，
// 匹配 TA-Lib 约定。
func MACD(input []float64, fastPeriod, slowPeriod, signalPeriod int) (*MACDResult, error) {
	if err := ValidateNumericInput(input, "MACD:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(fastPeriod, len(input), "MACD:fastPeriod"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(slowPeriod, len(input), "MACD:slowPeriod"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(signalPeriod, len(input), "MACD:signalPeriod"); err != nil {
		return nil, err
	}
	if fastPeriod >= slowPeriod {
		return nil, ErrFastPeriodGreaterOrEqualSlow
	}

	n := len(input)

	// 计算 EMA(快速) 和 EMA(慢速)
	emaFast, err := EMA(input, fastPeriod)
	if err != nil {
		return nil, err
	}

	emaSlow, err := EMA(input, slowPeriod)
	if err != nil {
		return nil, err
	}

	// MACD 线 = EMA(快速) - EMA(慢速)
	macd := MakeOutput(n)
	for i := slowPeriod - 1; i < n; i++ {
		macd[i] = emaFast[i] - emaSlow[i]
	}

	// 信号线 = EMA(MACD, 信号周期)
	signal, err := EMA(macd, signalPeriod)
	if err != nil {
		return nil, err
	}

	// 柱状图 = MACD - 信号
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

// MACDLookback 返回 MACD 输出中前导 NaN 值的总数。
func MACDLookback(fastPeriod, slowPeriod, signalPeriod int) int {
	return slowPeriod + signalPeriod - 2
}
