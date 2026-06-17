package talib

// AroonResult 包含 Aroon 上升和 Aroon 下降指标。
type AroonResult struct {
	Up   []float64 // Aroon 上升：(周期 - 最高价以来的K线数) / 周期 * 100
	Down []float64 // Aroon 下降：(周期 - 最低价以来的K线数) / 周期 * 100
}

// AROON 计算 Aroon 指标。
//
// Aroon 上升 = (周期 - 最高价以来的K线数) / 周期 * 100
// Aroon 下降 = (周期 - 最低价以来的K线数) / 周期 * 100
func AROON(high, low []float64, period int) (*AroonResult, error) {
	if err := ValidateNumericInput(high, "AROON:high"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(low, "AROON:low"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(high), "AROON:period"); err != nil {
		return nil, err
	}
	if len(high) != len(low) {
		return nil, ErrInputLengthMismatch("AROON: high and low")
	}

	n := len(high)
	up := MakeOutput(n)
	down := MakeOutput(n)

	for i := period; i < n; i++ {
		// 在窗口内查找最高高点和最低低点
		highest := high[i-period]
		lowest := low[i-period]
		highestIdx := i
		lowestIdx := i

		for j := i - period; j <= i; j++ {
			if high[j] >= highest { // 最新的出现获胜 (>=)
				highest = high[j]
				highestIdx = j
			}
			if low[j] <= lowest {
				lowest = low[j]
				lowestIdx = j
			}
		}

		barsSinceHigh := i - highestIdx
		barsSinceLow := i - lowestIdx

		up[i] = float64(period-barsSinceHigh) / float64(period) * 100.0
		down[i] = float64(period-barsSinceLow) / float64(period) * 100.0
	}

	return &AroonResult{Up: up, Down: down}, nil
}

// AROONLookback 返回 AROON 的 lookback。
func AROONLookback(period int) int {
	return period
}

// AROONOSC 计算 Aroon 振荡器。
//
// AroonOsc = Aroon 上升 - Aroon 下降
func AROONOSC(high, low []float64, period int) ([]float64, error) {
	aroon, err := AROON(high, low, period)
	if err != nil {
		return nil, err
	}

	n := len(high)
	out := MakeOutput(n)
	for i := period; i < n; i++ {
		out[i] = aroon.Up[i] - aroon.Down[i]
	}

	return out, nil
}

// AROONOSCLookback 返回 AROONOSC 的 lookback。
func AROONOSCLookback(period int) int {
	return period
}
