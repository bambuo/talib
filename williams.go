package talib

// WILLIAMS_R 计算威廉 %R 动能振荡器。
//
// %R = ((最高高 - 收盘价) / (最高高 - 最低低)) * (-100)
//
// 数值范围从 -100（超卖）到 0（超买）。
func WILLIAMS_R(high, low, close []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(close), "WILLIAMS_R:period"); err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		highest := high[i]
		lowest := low[i]
		for j := i - period + 1; j <= i; j++ {
			if high[j] > highest {
				highest = high[j]
			}
			if low[j] < lowest {
				lowest = low[j]
			}
		}

		denom := highest - lowest
		if denom == 0 {
			out[i] = 0
		} else {
			out[i] = ((highest - close[i]) / denom) * (-100.0)
		}
	}

	return out, nil
}

// WILLIAMSRLookback 返回 WILLIAMS_R 输出中前导 NaN 值的数量。
func WILLIAMSRLookback(period int) int {
	return period - 1
}
