package talib

import (
	"math"
)

// MIDPOINT 计算滑动窗口内的中点。
//
//	MIDPOINT[i] = (最高高 + 最低低) / 2
//
// 其中最高高和最低低是在窗口 [i-周期+1, i] 内的值。
// 前 (周期-1) 个元素为 NaN。
func MIDPOINT(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "MIDPOINT:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "MIDPOINT:period"); err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		highest := input[i-period+1]
		lowest := input[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if input[j] > highest {
				highest = input[j]
			}
			if input[j] < lowest {
				lowest = input[j]
			}
		}
		out[i] = (highest + lowest) / 2.0
	}

	return out, nil
}

// MIDPOINTLookback 返回 MIDPOINT 的 lookback。
func MIDPOINTLookback(period int) int {
	return period - 1
}

// MIDPRICE 计算滑动窗口内的中点价。
//
//	MIDPRICE[i] = (avg(周期高) + avg(周期低)) / 2
//
// 其中周期高是每个子区间内最高高点，
// 周期低是每个子区间内最低低点。
// 前 (周期-1) 个元素为 NaN。
func MIDPRICE(high, low []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(high, "MIDPRICE:high"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(low, "MIDPRICE:low"); err != nil {
		return nil, err
	}
	if len(high) != len(low) {
		return nil, ErrInputLengthMismatch("MIDPRICE: high, low")
	}
	if err := ValidatePeriod(period, len(high), "MIDPRICE:period"); err != nil {
		return nil, err
	}

	n := len(high)
	out := MakeOutput(n)

	for i := period - 1; i < n; i++ {
		// 实际上 MIDPRICE 计算窗口内最高高的平均值
		// 和窗口内最低低的平均值。
		// 在 TA-Lib 中，MIDPRICE = (最高(high, period) + 最低(low, period)) / 2
		highest := high[i-period+1]
		lowest := low[i-period+1]
		for j := i - period + 2; j <= i; j++ {
			if high[j] > highest {
				highest = high[j]
			}
			if low[j] < lowest {
				lowest = low[j]
			}
		}
		out[i] = (highest + lowest) / 2.0
	}

	return out, nil
}

// MIDPRICELookback 返回 MIDPRICE 的 lookback。
func MIDPRICELookback(period int) int {
	return period - 1
}

// SAR 计算抛物线 SAR（停损和反转）。
//
// 抛物线 SAR 是一个趋势跟踪指标，提供
// 跟踪止损水平。它使用加速度因子（AF），
// 随着趋势延伸而增加，导致 SAR 向
// 价格收敛。
//
// 参数：
//   - acceleration：初始加速度因子（通常为 0.02）
//   - maximum：最大加速度因子（通常为 0.20）
//
// 前导元素可能为 NaN，直到第一个有效 SAR 值。
func SAR(high, low []float64, acceleration, maximum float64) ([]float64, error) {
	if err := ValidateNumericInput(high, "SAR:high"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(low, "SAR:low"); err != nil {
		return nil, err
	}
	if len(high) != len(low) {
		return nil, ErrInputLengthMismatch("SAR: high, low")
	}

	n := len(high)
	out := MakeOutput(n)

	if n < 2 {
		return out, nil
	}

	// 确定初始趋势方向
	// 从上升趋势开始，使用前两根K线的最低价
	// 作为初始 SAR 值。
	isLong := true // true = 上升趋势, false = 下降趋势

	// 初始 SAR = 前两个最低价中的最低值
	sar := low[0]
	if low[1] < sar {
		sar = low[1]
	}

	// 极值点：趋势开始以来的最高高（对于多头）
	// 或趋势开始以来的最低低（对于空头）
	ep := high[0]
	if high[1] > ep {
		ep = high[1]
	}

	af := acceleration

	out[0] = math.NaN()
	out[1] = sar

	for i := 2; i < n; i++ {
		prevSAR := sar

		if isLong {
			// 上升趋势中
			// SAR = 前一个SAR + AF * (EP - 前一个SAR)
			sar = prevSAR + af*(ep-prevSAR)

			// SAR 必须低于前两根K线的最低价
			low1 := low[i-1]
			low2 := low[i-2]
			if sar > low1 {
				sar = low1
			}
			if sar > low2 {
				sar = low2
			}

			// 检查趋势反转：价格跌破 SAR
			if low[i] < sar {
				// 反转为下降趋势
				isLong = false
				sar = ep // SAR becomes the previous EP
				ep = low[i]
				af = acceleration
				// 在反转中，此K线的 SAR 是前一个 EP
				sar = ep
				ep = low[i]
				out[i] = sar
				continue
			}

			// 更新极值点
			if high[i] > ep {
				ep = high[i]
				af += acceleration
				if af > maximum {
					af = maximum
				}
			}
		} else {
			// 下降趋势中
			sar = prevSAR - af*(prevSAR-ep)

			// SAR 必须高于前两根K线的最高价
			high1 := high[i-1]
			high2 := high[i-2]
			if sar < high1 {
				sar = high1
			}
			if sar < high2 {
				sar = high2
			}

			// 检查趋势反转：价格突破 SAR 上行
			if high[i] > sar {
				// 反转为上升趋势
				isLong = true
				sar = ep // SAR becomes the previous EP
				ep = high[i]
				af = acceleration
				out[i] = sar
				continue
			}

			// Update Extreme Point
			if low[i] < ep {
				ep = low[i]
				af += acceleration
				if af > maximum {
					af = maximum
				}
			}
		}

		out[i] = sar
	}

	return out, nil
}

// SARLookback 返回 SAR 的 lookback。
func SARLookback() int {
	return 0
}

// SAREXT 计算具有可配置起始/最大加速度因子和反转偏移量的抛物线 SAR。
//
// 参数：
//   - startAF：初始加速度因子（例如 0.02）
//   - maxAF：最大加速度因子（例如 0.20）
//   - offsetOnReverse：反转趋势时应用的偏移量（0 = 标准 SAR）
//
// 这是 SAR 的扩展版本，具有与 TA-Lib 的 SAREXT 函数匹配的额外参数。
func SAREXT(high, low []float64, startAF, maxAF, offsetOnReverse float64) ([]float64, error) {
	if err := ValidateNumericInput(high, "SAREXT:high"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(low, "SAREXT:low"); err != nil {
		return nil, err
	}
	if len(high) != len(low) {
		return nil, ErrInputLengthMismatch("SAREXT: high, low")
	}

	n := len(high)
	out := MakeOutput(n)

	if n < 2 {
		return out, nil
	}

	isLong := true

	sar := low[0]
	if low[1] < sar {
		sar = low[1]
	}

	ep := high[0]
	if high[1] > ep {
		ep = high[1]
	}

	af := startAF

	out[0] = math.NaN()
	out[1] = sar

	for i := 2; i < n; i++ {
		prevSAR := sar

		if isLong {
			sar = prevSAR + af*(ep-prevSAR)

			// SAR must be below the lowest low of the previous two bars
			low1 := low[i-1]
			low2 := low[i-2]
			if sar > low1 {
				sar = low1
			}
			if sar > low2 {
				sar = low2
			}

			if low[i] < sar {
				isLong = false
				ep = low[i]
				af = startAF
				// 反转时，SAR = EP + offsetOnReverse
				sar = ep + offsetOnReverse
				out[i] = sar
				continue
			}

			if high[i] > ep {
				ep = high[i]
				af += startAF
				if af > maxAF {
					af = maxAF
				}
			}
		} else {
			sar = prevSAR - af*(prevSAR-ep)

			high1 := high[i-1]
			high2 := high[i-2]
			if sar < high1 {
				sar = high1
			}
			if sar < high2 {
				sar = high2
			}

			if high[i] > sar {
				isLong = true
				ep = high[i]
				af = startAF
				sar = ep - offsetOnReverse
				out[i] = sar
				continue
			}

			if low[i] < ep {
				ep = low[i]
				af += startAF
				if af > maxAF {
					af = maxAF
				}
			}
		}

		out[i] = sar
	}

	return out, nil
}

// SAREXTLookback 返回 SAREXT 的 lookback。
func SAREXTLookback() int {
	return 0
}

// MAVP 计算变周期移动平均。
//
// 对于索引 i 处的每个输出元素，周期取自 periods[i]
// （限制在 [minPeriod, maxPeriod] 范围内），然后在适当的
// 窗口上使用指定的 MA 类型计算移动平均。
//
// 前导元素为 NaN，直到有足够的数据可用。
func MAVP(input, periods []float64, minPeriod, maxPeriod int, maType MAType) ([]float64, error) {
	if err := ValidateNumericInput(input, "MAVP:input"); err != nil {
		return nil, err
	}
	if err := ValidateNumericInput(periods, "MAVP:periods"); err != nil {
		return nil, err
	}
	if len(input) != len(periods) {
		return nil, ErrInputLengthMismatch("MAVP: input, periods")
	}

	n := len(input)
	out := MakeOutput(n)

	for i := 0; i < n; i++ {
		p := int(periods[i])
		if p < minPeriod {
			p = minPeriod
		}
		if p > maxPeriod {
			p = maxPeriod
		}
		if p < 1 {
			continue
		}

		// 在窗口 [i-p+1, i] 上计算 MA
		start := i - p + 1
		if start < 0 {
			continue
		}

		window := input[start : i+1]
		maVal, err := MA(window, p, maType)
		if err != nil {
			return nil, err
		}
		// MA 返回一个切片；最后一个元素是当前值
		out[i] = maVal[len(maVal)-1]
	}

	return out, nil
}

// MAVPLookback 返回 MAVP 的 lookback。
func MAVPLookback(minPeriod int) int {
	return minPeriod - 1
}
