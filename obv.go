package talib

// OBV 计算能量潮指标。
//
// OBV[0] = volume[0] 如果 close[0] > close[-1] 否则 -volume[0]（第一个周期为 0）
// OBV[i] = OBV[i-1] + volume[i]   如果 close[i] > close[i-1]
// OBV[i] = OBV[i-1]                如果 close[i] == close[i-1]
// OBV[i] = OBV[i-1] - volume[i]   如果 close[i] < close[i-1]
func OBV(close, volume []float64) ([]float64, error) {
	if err := ValidateNumericInput(close, "OBV:close"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(volume, "OBV:volume"); err != nil {
		return nil, err
	}
	if len(close) != len(volume) {
		return nil, ErrInputLengthMismatch("OBV: close and volume")
	}

	n := len(close)
	out := make([]float64, n)

	// OBV[0] = 0（没有前一个收盘价可比较）
	out[0] = 0

	for i := 1; i < n; i++ {
		if close[i] > close[i-1] {
			out[i] = out[i-1] + volume[i]
		} else if close[i] < close[i-1] {
			out[i] = out[i-1] - volume[i]
		} else {
			out[i] = out[i-1]
		}
	}

	return out, nil
}

// OBVLookback 返回 OBV 输出中前导 NaN 值的数量。
func OBVLookback() int {
	return 0
}
