package talib

// ULTOSC 计算终极振荡器。
//
// 使用三个时间周期（默认：7, 14, 28）。
//
// 买入压力 = 收盘价 - Min(最低价, 前收盘价)
// 真实波幅 = Max(最高价, 前收盘价) - Min(最低价, 前收盘价)
//
// 平均值N = SMA(买入压力, 周期N) / SMA(真实波幅, 周期N)
//
// UO = 100 * (4*平均7 + 2*平均14 + 平均28) / 7
func ULTOSC(high, low, close []float64, period1, period2, period3 int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	maxPeriod := period1
	if period2 > maxPeriod {
		maxPeriod = period2
	}
	if period3 > maxPeriod {
		maxPeriod = period3
	}
	if err := ValidatePeriod(maxPeriod, len(close), "ULTOSC:maxPeriod"); err != nil {
		return nil, err
	}

	n := len(close)

	// 计算买入压力和真实波幅
	bp := make([]float64, n)
	tr := make([]float64, n)
	for i := 1; i < n; i++ {
		prevClose := close[i-1]
		trHigh := high[i]
		if prevClose > trHigh {
			trHigh = prevClose
		}
		trLow := low[i]
		if prevClose < trLow {
			trLow = prevClose
		}
		bp[i] = close[i] - trLow
		tr[i] = trHigh - trLow
	}

	out := MakeOutput(n)

	startIdx := maxPeriod
	for i := startIdx; i < n; i++ {
		avg1 := windowAverage(bp, tr, i, period1)
		avg2 := windowAverage(bp, tr, i, period2)
		avg3 := windowAverage(bp, tr, i, period3)

		out[i] = 100.0 * (4.0*avg1 + 2.0*avg2 + avg3) / 7.0
	}

	return out, nil
}

// ULTOSCLookback 返回 ULTOSC 的 lookback。
func ULTOSCLookback(period1, period2, period3 int) int {
	return period3
}

func windowAverage(bp, tr []float64, idx, period int) float64 {
	var sumBP, sumTR float64
	for j := idx - period + 1; j <= idx; j++ {
		sumBP += bp[j]
		sumTR += tr[j]
	}
	if sumTR == 0 {
		return 0
	}
	return sumBP / sumTR
}
