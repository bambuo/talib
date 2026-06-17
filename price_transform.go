package talib

// AVGPRICE 计算平均价格。
//
//	AVGPRICE = (开盘价 + 最高价 + 最低价 + 收盘价) / 4
func AVGPRICE(open, high, low, close []float64) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(open, "AVGPRICE:open"); err != nil {
		return nil, err
	}
	if len(open) != len(close) {
		return nil, ErrInputLengthMismatch("AVGPRICE: open, close")
	}

	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (open[i] + high[i] + low[i] + close[i]) / 4.0
	}
	return out, nil
}

// MEDPRICE 计算中间价。
//
//	MEDPRICE = (最高价 + 最低价) / 2
func MEDPRICE(high, low []float64) ([]float64, error) {
	if err := ValidateNumericInput(high, "MEDPRICE:high"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(low, "MEDPRICE:low"); err != nil {
		return nil, err
	}
	if len(high) != len(low) {
		return nil, ErrInputLengthMismatch("MEDPRICE: high, low")
	}

	n := len(high)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i]) / 2.0
	}
	return out, nil
}

// TYPPRICE 计算典型价格。
//
//	TYPPRICE = (最高价 + 最低价 + 收盘价) / 3
func TYPPRICE(high, low, close []float64) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}

	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i] + close[i]) / 3.0
	}
	return out, nil
}

// WCLPRICE 计算加权收盘价。
//
//	WCLPRICE = (最高价 + 最低价 + 2*收盘价) / 4
func WCLPRICE(high, low, close []float64) ([]float64, error) {
	if err := ValidateOHLCInput(high, low, close); err != nil {
		return nil, err
	}

	n := len(close)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = (high[i] + low[i] + 2*close[i]) / 4.0
	}
	return out, nil
}
