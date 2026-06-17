package talib

// StochResult 包含随机振荡器的 %K 和 %D 线。
type StochResult struct {
	K []float64 // %K（快随机线）
	D []float64 // %D（慢随机线，%K 的 SMA）
}

// STOCH 计算随机振荡器。
//
// %K = (收盘价 - 最低低) / (最高高 - 最低低) * 100
// %D = SMA(%K, 信号周期)
//
// 使用慢速随机约定，其中 %K 在计算 %D 之前被 slowKPeriod 平滑（原始 %K 的 SMA）。
func STOCH(high, low, close []float64, fastKPeriod, slowKPeriod, slowDPeriod int, maType MAType) (*StochResult, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	n := len(close)

	// 计算原始 %K
	rawK := MakeOutput(n)
	for i := fastKPeriod - 1; i < n; i++ {
		highest := high[i]
		lowest := low[i]
		for j := i - fastKPeriod + 1; j <= i; j++ {
			if high[j] > highest {
				highest = high[j]
			}
			if low[j] < lowest {
				lowest = low[j]
			}
		}
		denom := highest - lowest
		if denom == 0 {
			rawK[i] = 0
		} else {
			rawK[i] = (close[i] - lowest) / denom * 100.0
		}
	}

	// 用 slowKPeriod SMA 平滑 %K
	k, err := SMA(rawK, slowKPeriod)
	if err != nil {
		return nil, err
	}

	// %D = 平滑后 %K 的 SMA
	d, err := MA(k, slowDPeriod, maType)
	if err != nil {
		return nil, err
	}

	return &StochResult{K: k, D: d}, nil
}

// STOCHLookback 返回随机振荡器的总 lookback。
func STOCHLookback(fastKPeriod, slowKPeriod, slowDPeriod int) int {
	return fastKPeriod + slowKPeriod + slowDPeriod - 3
}

// STOCHF 计算快随机振荡器。
//
// %K = (收盘价 - 最低低) / (最高高 - 最低低) * 100（不平滑）
// %D = SMA(%K, 信号周期)
func STOCHF(high, low, close []float64, fastKPeriod, fastDPeriod int, maType MAType) (*StochResult, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	n := len(close)

	// 计算原始 %K（不平滑）
	k := MakeOutput(n)
	for i := fastKPeriod - 1; i < n; i++ {
		highest := high[i]
		lowest := low[i]
		for j := i - fastKPeriod + 1; j <= i; j++ {
			if high[j] > highest {
				highest = high[j]
			}
			if low[j] < lowest {
				lowest = low[j]
			}
		}
		denom := highest - lowest
		if denom == 0 {
			k[i] = 0
		} else {
			k[i] = (close[i] - lowest) / denom * 100.0
		}
	}

	// %D = %K 的 MA
	d, err := MA(k, fastDPeriod, maType)
	if err != nil {
		return nil, err
	}

	return &StochResult{K: k, D: d}, nil
}

// STOCHFLookback 返回快随机振荡器的总 lookback。
func STOCHFLookback(fastKPeriod, fastDPeriod int) int {
	return fastKPeriod + fastDPeriod - 2
}

// STOCHRSI 计算随机 RSI。
//
// 首先计算 RSI(输入, 周期)，然后对 RSI 值应用随机公式，
// 对 %K 和 %D 均使用 Wilder 平滑。
func STOCHRSI(input []float64, period, fastKPeriod, fastDPeriod int, maType MAType) (*StochResult, error) {
	if err := ValidateNumericInput(input, "STOCHRSI:input"); err != nil {
		return nil, err
	}

	// 先计算 RSI
	rsi, err := RSI(input, period)
	if err != nil {
		return nil, err
	}
	n := len(rsi)

	// 对 RSI 值应用随机（与 STOCH 相同但作用于 RSI）
	k := MakeOutput(n)
	startIdx := period + fastKPeriod - 1
	for i := startIdx; i < n; i++ {
		highest := rsi[i]
		lowest := rsi[i]
		for j := i - fastKPeriod + 1; j <= i; j++ {
			if IsNaN(rsi[j]) {
				continue
			}
			if rsi[j] > highest || IsNaN(highest) {
				highest = rsi[j]
			}
			if rsi[j] < lowest || IsNaN(lowest) {
				lowest = rsi[j]
			}
		}
		denom := highest - lowest
		if denom == 0 {
			k[i] = 0
		} else {
			k[i] = (rsi[i] - lowest) / denom * 100.0
		}
	}

	d, err := MA(k, fastDPeriod, maType)
	if err != nil {
		return nil, err
	}

	return &StochResult{K: k, D: d}, nil
}

// STOCHRSILookback 返回随机 RSI 的总 lookback。
func STOCHRSILookback(period, fastKPeriod, fastDPeriod int) int {
	return period + fastKPeriod + fastDPeriod - 2
}
