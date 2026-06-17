package talib

// ADResult 包含 AD 线。
type ADResult struct {
	AD []float64
}

// AD 计算 Chaikin 累积/分配线。
//
// 资金流量乘数 = ((收盘价 - 最低价) - (最高价 - 收盘价)) / (最高价 - 最低价)
// AD 成交量 = MFL * 成交量
// AD 线 = AD 成交量的累积和
func AD(high, low, close, volume []float64) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(volume, "AD:volume"); err != nil {
		return nil, err
	}
	n := len(close)
	if len(volume) != n {
		return nil, ErrInputLengthMismatch("AD: close and volume")
	}

	out := make([]float64, n)

	var cumSum float64
	for i := 0; i < n; i++ {
		denom := high[i] - low[i]
		var mfl float64
		if denom == 0 {
			mfl = 0
		} else {
			mfl = ((close[i] - low[i]) - (high[i] - close[i])) / denom
		}
		cumSum += mfl * volume[i]
		out[i] = cumSum
	}

	return out, nil
}

// ADLookback 返回 AD 的 lookback。
func ADLookback() int {
	return 0
}

// ADOSC 计算 Chaikin 累积/分配振荡器。
//
// ADOSC = EMA(AD, 快速周期) - EMA(AD, 慢速周期)
func ADOSC(high, low, close, volume []float64, fastPeriod, slowPeriod int) ([]float64, error) {
	if fastPeriod >= slowPeriod {
		return nil, ErrFastPeriodGreaterOrEqualSlow
	}

	ad, err := AD(high, low, close, volume)
	if err != nil {
		return nil, err
	}

	// 计算 AD 线上的 EMA(快速) - EMA(慢速)
	emaFast, err := EMA(ad, fastPeriod)
	if err != nil {
		return nil, err
	}
	emaSlow, err := EMA(ad, slowPeriod)
	if err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)
	for i := slowPeriod - 1; i < n; i++ {
		out[i] = emaFast[i] - emaSlow[i]
	}

	return out, nil
}

// ADOSCLookback 返回 ADOSC 的 lookback。
func ADOSCLookback(fastPeriod, slowPeriod int) int {
	return slowPeriod - 1
}

// ADXR 计算平均方向运动评级。
//
// ADXR = (ADX[i] + ADX[i-周期]) / 2
func ADXR(high, low, close []float64, period int) ([]float64, error) {
	adx, err := ADX(high, low, close, period)
	if err != nil {
		return nil, err
	}

	n := len(close)
	out := MakeOutput(n)
	startIdx := 2*period + period
	for i := startIdx; i < n; i++ {
		if !IsNaN(adx[i]) && !IsNaN(adx[i-period]) {
			out[i] = (adx[i] + adx[i-period]) / 2.0
		}
	}

	return out, nil
}

// ADXRLookback 返回 ADXR 的 lookback。
func ADXRLookback(period int) int {
	return 3 * period
}
