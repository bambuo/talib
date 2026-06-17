package talib

// MFI 计算资金流量指数。
//
// 典型价格 = (最高价 + 最低价 + 收盘价) / 3
// 资金流量 = 典型价格 * 成交量
//
// 正资金流量 = 典型价格 > 前一典型价格时的资金流量
// 负资金流量 = 典型价格 < 前一典型价格时的资金流量
//
// 资金比率 = 正MF总和(周期) / 负MF总和(周期)
// MFI = 100 - 100 / (1 + 资金比率)
func MFI(high, low, close, volume []float64, period int) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(volume, "MFI:volume"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(close), "MFI:period"); err != nil {
		return nil, err
	}

	n := len(close)
	if len(volume) != n {
		return nil, ErrInputLengthMismatch("MFI: close and volume")
	}

	// 计算典型价格和资金流量
	tp := make([]float64, n)
	mf := make([]float64, n)
	for i := 0; i < n; i++ {
		tp[i] = (high[i] + low[i] + close[i]) / 3.0
		mf[i] = tp[i] * volume[i]
	}

	out := MakeOutput(n)

	// 滑动窗口 MFI
	for i := period; i < n; i++ {
		var posMF, negMF float64
		for j := i - period + 1; j <= i; j++ {
			if tp[j] > tp[j-1] {
				posMF += mf[j]
			} else if tp[j] < tp[j-1] {
				negMF += mf[j]
			}
		}

		if negMF == 0 {
			out[i] = 100.0
		} else if posMF == 0 {
			out[i] = 0.0
		} else {
			mr := posMF / negMF
			out[i] = 100.0 - 100.0/(1.0+mr)
		}
	}

	return out, nil
}

// MFILookback 返回 MFI 输出中前导 NaN 值的数量。
func MFILookback(period int) int {
	return period
}
