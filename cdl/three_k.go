package cdl

import "math"

// =====================================================================
// Three+ Candle Patterns — 三K线及以上形态
// =====================================================================

// CDL2CROWS 识别两只乌鸦形态。
// 上升趋势中一根阳线后跟两根阴线。
// 第一根阴线向上跳空，第二根阴线吞没第一根。
// 看跌反转。
// 最小回溯：3
func CDL2CROWS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 蜡烛 1：长阳线
		o0, c0, h0, l0 := open[i-2], close[i-2], high[i-2], low[i-2]
		if !whiteBody(o0, c0) || !isLongBody(o0, c0, h0, l0) {
			continue
		}
		// 蜡烛 2：阴线，向上跳空
		o1, c1, l1 := open[i-1], close[i-1], low[i-1]
		if !blackBody(o1, c1) || !isGapUp(h0, l1) {
			continue
		}
		// 蜡烛 3：阴线，开盘在第二根内部，收盘在第二根下方，进入第一根实体
		o2, c2 := open[i], close[i]
		if !blackBody(o2, c2) {
			continue
		}
		if o2 < o1 && o2 > c1 && c2 < c1 && c2 < o0 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDL2CROWSLookback 返回 CDL2CROWS 的最小回溯期。
func CDL2CROWSLookback() int { return 3 }

// CDL3BLACKCROWS 识别三只乌鸦形态。
// 连续三根长阴线，每根开盘在前一根实体之内，
// 收盘接近其最低价，收盘价依次走低。
// 看跌反转。
// 最小回溯：4
func CDL3BLACKCROWS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 需要前期的上升趋势
		upCount := 0
		for j := i - 4; j < i-2; j++ {
			if j >= 0 && isBullish(open[j], close[j]) {
				upCount++
			}
		}
		if upCount < 1 {
			continue
		}
		// 连续三根长阴线
		ok := true
		for k := 0; k < 3; k++ {
			idx := i - 2 + k
			o, c, h, l := open[idx], close[idx], high[idx], low[idx]
			if !blackBody(o, c) || !isLongBody(o, c, h, l) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// 每根开盘在前一根实体之内，每根收盘接近其最低价
		o0, c0 := open[i-2], close[i-2]
		o1, c1 := open[i-1], close[i-1]
		o2, c2, l2 := open[i], close[i], low[i]
		if o1 < c0 && o1 > o0 && o2 < c1 && o2 > o1 &&
			c0 > c1 && c1 > c2 &&
			lowerShadow(low[i-2], open[i-2], close[i-2]) <= realBody(open[i-2], close[i-2])*0.1 &&
			lowerShadow(low[i-1], open[i-1], close[i-1]) <= realBody(open[i-1], close[i-1])*0.1 &&
			lowerShadow(l2, o2, c2) <= realBody(o2, c2)*0.1 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDL3BLACKCROWSLookback 返回 CDL3BLACKCROWS 的最小回溯期。
func CDL3BLACKCROWSLookback() int { return 4 }

// CDL3INSIDE 识别三内部上涨/下跌形态。
// 看涨：第一根阴线，第二根阳线在第一根内部，第三根阳线确认（高于第一根）。
// 看跌：第一根阳线，第二根阴线在第一根内部，第三根阴线确认（低于第一根）。
// 最小回溯：3
func CDL3INSIDE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		o0, c0 := open[i-2], close[i-2]
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		// 三外部上涨（看涨）：阴线，阳线孕线，阳线确认
		if blackBody(o0, c0) && whiteBody(o1, c1) && whiteBody(o2, c2) {
			// 第二根实体在第一根内部
			if math.Min(o1, c1) > math.Min(o0, c0) && math.Max(o1, c1) < math.Max(o0, c0) {
				// 第三根收盘价高于第一根
				if c2 > c0 {
					out[i] = Bullish
				}
			}
		}
		// 三内部下跌（看跌）：阳线，阴线孕线，阴线确认
		if whiteBody(o0, c0) && blackBody(o1, c1) && blackBody(o2, c2) {
			if math.Min(o1, c1) > math.Min(o0, c0) && math.Max(o1, c1) < math.Max(o0, c0) {
				if c2 < c0 {
					out[i] = Bearish
				}
			}
		}
	}
	return out, nil
}

// CDL3INSIDELookback 返回 CDL3INSIDE 的最小回溯期。
func CDL3INSIDELookback() int { return 3 }

// CDLINSIDE 识别内部上涨/下跌（三内部上涨/下跌的别名）。
// 详见 CDL3INSIDE。
func CDLINSIDE(open, high, low, close []float64) ([]int, error) {
	return CDL3INSIDE(open, high, low, close)
}

// CDLINSIDELookback 返回 CDLINSIDE 的最小回溯期。
func CDLINSIDELookback() int { return 3 }

// CDL3OUTSIDE 识别三外部上涨/下跌形态。
// 类似于三内部形态，但第二根蜡烛吞没第一根（与孕线相反）。
// 最小回溯：3
func CDL3OUTSIDE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		o0, c0 := open[i-2], close[i-2]
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		// 三外部上涨（看涨）：阴线，阳线吞没，阳线确认
		if blackBody(o0, c0) && whiteBody(o1, c1) && whiteBody(o2, c2) {
			if math.Min(o1, c1) < math.Min(o0, c0) && math.Max(o1, c1) > math.Max(o0, c0) {
				if c2 > c1 {
					out[i] = Bullish
				}
			}
		}
		// 三外部下跌（看跌）：阳线，阴线吞没，阴线确认
		if whiteBody(o0, c0) && blackBody(o1, c1) && blackBody(o2, c2) {
			if math.Min(o1, c1) < math.Min(o0, c0) && math.Max(o1, c1) > math.Max(o0, c0) {
				if c2 < c1 {
					out[i] = Bearish
				}
			}
		}
	}
	return out, nil
}

// CDL3OUTSIDELookback 返回 CDL3OUTSIDE 的最小回溯期。
func CDL3OUTSIDELookback() int { return 3 }

// CDL3WHITESOLDIERS 识别三白兵形态。
// 连续三根长阳线，每根开盘在前一根实体之内，
// 收盘接近其最高价，收盘价依次走高。
// 看涨反转形态。
// 最小回溯：4
func CDL3WHITESOLDIERS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 需要前期的下降趋势
		downCount := 0
		for j := i - 4; j < i-2; j++ {
			if j >= 0 && isBearish(open[j], close[j]) {
				downCount++
			}
		}
		if downCount < 1 {
			continue
		}
		// 连续三根长阳线
		ok := true
		for k := 0; k < 3; k++ {
			idx := i - 2 + k
			o, c, h, l := open[idx], close[idx], high[idx], low[idx]
			if !whiteBody(o, c) || !isLongBody(o, c, h, l) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// 每根开盘在前一根实体之内，每根收盘接近其最高价
		o0, c0 := open[i-2], close[i-2]
		o1, c1 := open[i-1], close[i-1]
		o2, c2, h2 := open[i], close[i], high[i]
		if o1 > o0 && o1 < c0 && o2 > o1 && o2 < c1 &&
			c0 < c1 && c1 < c2 &&
			upperShadow(high[i-2], open[i-2], close[i-2]) <= realBody(open[i-2], close[i-2])*0.1 &&
			upperShadow(high[i-1], open[i-1], close[i-1]) <= realBody(open[i-1], close[i-1])*0.1 &&
			upperShadow(h2, o2, c2) <= realBody(o2, c2)*0.1 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDL3WHITESOLDIERSLookback 返回 CDL3WHITESOLDIERS 的最小回溯期。
func CDL3WHITESOLDIERSLookback() int { return 4 }

// CDL3STARSINSOUTH 识别南方三星形态。
// 三根阴线，实体和下影线逐渐缩小，
// 暗示看跌动能减弱。看涨反转。
// 最小回溯：4
func CDL3STARSINSOUTH(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 三根阴线
		if !blackBody(open[i-2], close[i-2]) || !blackBody(open[i-1], close[i-1]) || !blackBody(open[i], close[i]) {
			continue
		}
		// 实体逐渐变小
		b0 := realBody(open[i-2], close[i-2])
		b1 := realBody(open[i-1], close[i-1])
		b2 := realBody(open[i], close[i])
		if b0 < b1 || b1 < b2 {
			continue
		}
		// 下影线逐渐变小
		ls0 := lowerShadow(low[i-2], open[i-2], close[i-2])
		ls1 := lowerShadow(low[i-1], open[i-1], close[i-1])
		ls2 := lowerShadow(low[i], open[i], close[i])
		if ls0 < ls1 || ls1 < ls2 {
			continue
		}
		// 收盘依次降低
		if close[i-2] > close[i-1] && close[i-1] > close[i] {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDL3STARSINSOUTHLookback 返回 CDL3STARSINSOUTH 的最小回溯期。
func CDL3STARSINSOUTHLookback() int { return 4 }

// CDLIDENTICAL3CROWS 识别相同三只乌鸦形态。
// 三根等长的阴线，收盘依次降低。
// 看跌反转。
// 最小回溯：4
func CDLIDENTICAL3CROWS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 三根阴线
		if !blackBody(open[i-2], close[i-2]) || !blackBody(open[i-1], close[i-1]) || !blackBody(open[i], close[i]) {
			continue
		}
		// 实体大致相等
		b0 := realBody(open[i-2], close[i-2])
		b1 := realBody(open[i-1], close[i-1])
		b2 := realBody(open[i], close[i])
		if b0 == 0 {
			continue
		}
		if math.Abs(b0-b1)/b0 > 0.1 || math.Abs(b1-b2)/b1 > 0.1 {
			continue
		}
		// 收盘必须依次降低
		if close[i-2] > close[i-1] && close[i-1] > close[i] {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLIDENTICAL3CROWSLookback 返回 CDLIDENTICAL3CROWS 的最小回溯期。
func CDLIDENTICAL3CROWSLookback() int { return 4 }

// CDLADVANCEBLOCK 识别 Advance Block（前进阻挡）形态。
// 三根阳线，实体逐渐缩小和/或上影线逐渐增长。
// 看跌反转形态。
// 最小回溯：4
func CDLADVANCEBLOCK(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 上升趋势中的三根阳线
		if !whiteBody(open[i-2], close[i-2]) || !whiteBody(open[i-1], close[i-1]) || !whiteBody(open[i], close[i]) {
			continue
		}
		// 实体逐渐变小或上影线逐渐变大
		b0 := realBody(open[i-2], close[i-2])
		b1 := realBody(open[i-1], close[i-1])
		b2 := realBody(open[i], close[i])
		us0 := upperShadow(high[i-2], open[i-2], close[i-2])
		us1 := upperShadow(high[i-1], open[i-1], close[i-1])
		us2 := upperShadow(high[i], open[i], close[i])
		bodyShrinking := b0 > b1 && b1 > b2
		shadowGrowing := us0 < us1 && us1 < us2
		// 至少需要其中一个信号
		if bodyShrinking || shadowGrowing {
			if close[i-2] < close[i-1] && close[i-1] < close[i] {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLADVANCEBLOCKLookback 返回 CDLADVANCEBLOCK 的最小回溯期。
func CDLADVANCEBLOCKLookback() int { return 4 }

// CDLSTALLEDPATTERN 识别停顿形态。
// 类似于 Advance Block，但第三根蜡烛实体显著变小
// 且开盘于第二根收盘价附近。
// 看跌反转。
// 最小回溯：3
func CDLSTALLEDPATTERN(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		if !whiteBody(open[i-2], close[i-2]) || !whiteBody(open[i-1], close[i-1]) || !whiteBody(open[i], close[i]) {
			continue
		}
		b1 := realBody(open[i-1], close[i-1])
		b2 := realBody(open[i], close[i])
		// 第三根实体显著变小
		if b2 > b1*0.5 {
			continue
		}
		// 开盘价在第二根收盘价附近
		if math.Abs(open[i]-close[i-1]) < realBody(open[i], close[i])*0.5 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLSTALLEDPATTERNLookback 返回 CDLSTALLEDPATTERN 的最小回溯期。
func CDLSTALLEDPATTERNLookback() int { return 3 }

// CDLUNIQUE3RIVER 识别 Unique Three River（独特三河）形态。
// 三根蜡烛的看涨反转形态：
// 1. 下降趋势中的长阴线
// 2. 在第一根内部的阴线孕线
// 3. 具有更低最低价的小阳线
// 最小回溯：4
func CDLUNIQUE3RIVER(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 蜡烛 1：长阴线
		o0, c0, h0, l0 := open[i-2], close[i-2], high[i-2], low[i-2]
		if !blackBody(o0, c0) || !isLongBody(o0, c0, h0, l0) {
			continue
		}
		// 蜡烛 2：阴线，在第一根内部
		o1, c1 := open[i-1], close[i-1]
		if !blackBody(o1, c1) {
			continue
		}
		if math.Min(o1, c1) <= math.Min(o0, c0) || math.Max(o1, c1) >= math.Max(o0, c0) {
			continue
		}
		// 蜡烛 3：小阳线，新低，收盘在第二根实体之内
		o2, c2, l2 := open[i], close[i], low[i]
		if !whiteBody(o2, c2) {
			continue
		}
		if l2 >= low[i-1] || c2 <= c1 || c2 >= o1 {
			continue
		}
		out[i] = Bullish
	}
	return out, nil
}

// CDLUNIQUE3RIVERLookback 返回 CDLUNIQUE3RIVER 的最小回溯期。
func CDLUNIQUE3RIVERLookback() int { return 4 }

// CDLUPSIDEGAP2CROWS 识别向上跳空两只乌鸦形态。
// 一根阳线，向上跳空，两根阴线。第二根阴线吞没第一根阴线。
// 看跌反转。
// 最小回溯：3
func CDLUPSIDEGAP2CROWS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 蜡烛 1：上升趋势中的长阳线
		o0, c0, h0, l0 := open[i-2], close[i-2], high[i-2], low[i-2]
		if !whiteBody(o0, c0) || !isLongBody(o0, c0, h0, l0) {
			continue
		}
		// 蜡烛 2：阴线，向上跳空
		o1, c1, l1 := open[i-1], close[i-1], low[i-1]
		if !blackBody(o1, c1) || !isGapUp(h0, l1) {
			continue
		}
		// 蜡烛 3：阴线，开盘高于第二根，收盘低于第二根
		o2, c2 := open[i], close[i]
		if !blackBody(o2, c2) {
			continue
		}
		if o2 > o1 && c2 < c1 && c2 < o0 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLUPSIDEGAP2CROWSLookback 返回 CDLUPSIDEGAP2CROWS 的最小回溯期。
func CDLUPSIDEGAP2CROWSLookback() int { return 3 }

// CDLHIKKAKE 识别 Hikkake（引挂）形态。
// 两根蜡烛形成内部 bar（孕线），后跟第三根蜡烛假突破。
// 看涨：内部 bar 的低点被跌破但回升。
// 看跌：内部 bar 的高点被突破但反转。
// 最小回溯：4
func CDLHIKKAKE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 蜡烛 1（母蜡烛）：设置蜡烛
		h0, l0 := high[i-3], low[i-3]
		// 蜡烛 2：内部 bar（在蜡烛 1 范围内）
		h1, l1 := high[i-2], low[i-2]
		if h1 > h0 || l1 < l0 {
			continue
		}
		// 蜡烛 3：跌破设置低点（看涨陷阱）或突破设置高点（看跌陷阱）
		h2, l2 := high[i-1], low[i-1]
		// 看涨 Hikkake：假跌破
		if l2 < l0 && close[i-1] > l0 {
			// 确认：当前 bar 收盘更高
			if close[i] > close[i-1] {
				out[i] = Bullish
			}
		}
		// 看跌 Hikkake：假突破
		if h2 > h0 && close[i-1] < h0 {
			if close[i] < close[i-1] {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLHIKKAKELookback 返回 CDLHIKKAKE 的最小回溯期。
func CDLHIKKAKELookback() int { return 4 }

// CDLHIKKAKEMOD 识别修改版 Hikkake 形态。
// 类似于 Hikkake，但具有不同的确认机制。
// 最小回溯：4
func CDLHIKKAKEMOD(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		h0, l0 := high[i-3], low[i-3]
		h1, l1 := high[i-2], low[i-2]
		// 内部 bar
		if h1 > h0 || l1 < l0 {
			continue
		}
		// 寻找低点跌破母蜡烛后收盘回上方
		if low[i-1] < l0 && close[i-1] > l0 {
			// 确认：当前收盘价高于第三根蜡烛（与 Hikkake 相同）
			if close[i] > high[i-1] {
				out[i] = Bullish
			} else if close[i] < low[i-1] {
				out[i] = Bearish
			}
		}
		if high[i-1] > h0 && close[i-1] < h0 {
			if close[i] < low[i-1] {
				out[i] = Bearish
			} else if close[i] > high[i-1] {
				out[i] = Bullish
			}
		}
	}
	return out, nil
}

// CDLHIKKAKEMODLookback 返回 CDLHIKKAKEMOD 的最小回溯期。
func CDLHIKKAKEMODLookback() int { return 4 }

// CDL3LINESTRIKE 识别三线打击形态。
// 四根蜡烛形态：
// 看涨：三根阴线收盘依次走低，然后一根长阳线
// 开盘在第一根阴线开盘价下方，收盘在其上方。
// 看跌：三根阳线收盘依次走高，然后一根长阴线
// 开盘在第一根阳线开盘价上方，收盘在其下方。
// 最小回溯：4
func CDL3LINESTRIKE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		// 看涨三线打击：3 根阴线后跟长阳线
		if blackBody(open[i-3], close[i-3]) && blackBody(open[i-2], close[i-2]) && blackBody(open[i-1], close[i-1]) {
			if close[i-3] > close[i-2] && close[i-2] > close[i-1] {
				o, c := open[i], close[i]
				if whiteBody(o, c) && c > open[i-3] {
					out[i] = Bullish
				}
			}
		}
		// 看跌三线打击：3 根阳线后跟长阴线
		if whiteBody(open[i-3], close[i-3]) && whiteBody(open[i-2], close[i-2]) && whiteBody(open[i-1], close[i-1]) {
			if close[i-3] < close[i-2] && close[i-2] < close[i-1] {
				o, c := open[i], close[i]
				if blackBody(o, c) && c < open[i-3] {
					out[i] = Bearish
				}
			}
		}
	}
	return out, nil
}

// CDL3LINESTRIKELookback 返回 CDL3LINESTRIKE 的最小回溯期。
func CDL3LINESTRIKELookback() int { return 4 }
