package talib

// RSI 使用 Wilder 平滑方法计算相对强弱指标。
//
// RSI = 100 - 100 / (1 + RS)
//
// 其中 RS = 平均涨幅 / 平均跌幅，在给定周期内。
//
// 首次平均值计算为初始周期内的 SMA，
// 后续值使用 Wilder 的 EMA（k = 1/周期）。
//
// 这与 TA-Lib 的 C 语言实现完全一致。
func RSI(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "RSI:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "RSI:period"); err != nil {
		return nil, err
	}

	out := MakeOutput(len(input))
	n := len(input)

	// 如果输入长度 <= 周期则无法计算（仅 NaN 输出）
	if n <= period {
		return out, nil
	}

	// 第一遍：计算涨跌幅
	gains := make([]float64, n)
	losses := make([]float64, n)

	for i := 1; i < n; i++ {
		diff := input[i] - input[i-1]
		if diff > 0 {
			gains[i] = diff
		} else {
			losses[i] = -diff
		}
	}

	// 初始平均值：从索引 1 开始的第一个 (周期+1) 个值的 SMA
	//（因为 gains[0] = losses[0] = 0）
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// 索引 period 处的第一个 RSI 值
	if avgLoss == 0 {
		out[period] = 100.0
	} else {
		rs := avgGain / avgLoss
		out[period] = 100.0 - 100.0/(1.0+rs)
	}

	// Wilder 平滑：k = 1/周期
	oneOverPeriod := 1.0 / float64(period)

	for i := period + 1; i < n; i++ {
		avgGain = (avgGain*float64(period-1) + gains[i]) * oneOverPeriod
		avgLoss = (avgLoss*float64(period-1) + losses[i]) * oneOverPeriod

		if avgLoss == 0 {
			out[i] = 100.0
		} else {
			rs := avgGain / avgLoss
			out[i] = 100.0 - 100.0/(1.0+rs)
		}
	}

	return out, nil
}

// RSILookback 返回 RSI 输出中前导 NaN 值的数量。
func RSILookback(period int) int {
	return period
}
