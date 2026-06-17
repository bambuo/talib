package cdl

// =====================================================================
// Complex Multi-Candle Patterns — 复杂多K线形态
// =====================================================================

// CDLBREAKAWAY 识别 Breakaway（脱离）形态。
// 一个五根蜡烛的形态，标志着主要趋势反转：
// 看涨：一根长阴线，三根依次走低的蜡烛，
//
//	然后一根长阳线收盘在前三根最高价之上。
//
// 看跌：一根长阳线，三根依次走高的蜡烛，
//
//	然后一根长阴线收盘在前三根最低价之下。
//
// 最小回溯：5
func CDLBREAKAWAY(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 4; i < n; i++ {
		o0, c0, h0 := open[i-4], close[i-4], high[i-4]
		c1, c2, c3 := close[i-3], close[i-2], close[i-1]
		o4, c4, l4 := open[i], close[i], low[i]
		// 看涨 Breakaway
		if blackBody(o0, c0) && isLongBody(o0, c0, h0, low[i-4]) {
			// 三个依次降低的收盘价
			if c0 > c1 && c1 > c2 && c2 > c3 {
				// 第五根：长阳线，收盘在前几根最高价之上
				if whiteBody(o4, c4) && isLongBody(o4, c4, high[i], l4) {
					maxHigh := mathMax3(high[i-3], high[i-2], high[i-1])
					if c4 > maxHigh {
						out[i] = Bullish
					}
				}
			}
		}
		// 看跌 Breakaway
		if whiteBody(o0, c0) && isLongBody(o0, c0, h0, low[i-4]) {
			if c0 < c1 && c1 < c2 && c2 < c3 {
				if blackBody(o4, c4) && isLongBody(o4, c4, high[i], l4) {
					minLow := mathMin3(low[i-3], low[i-2], low[i-1])
					if c4 < minLow {
						out[i] = Bearish
					}
				}
			}
		}
	}
	return out, nil
}

// CDLBREAKAWAYLookback 返回 CDLBREAKAWAY 的最小回溯期。
func CDLBREAKAWAYLookback() int { return 5 }

// CDLMATHOLD 识别 Mat Hold（铺垫）形态。
// 一个五根蜡烛的看涨持续形态：
// 1. 长阳线
// 2-4. 三根走低的小阴线，保持在第一根蜡烛范围内
// 5. 长阳线收盘在前三根最高价之上
// 最小回溯：5
func CDLMATHOLD(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 4; i < n; i++ {
		// 蜡烛 1：长阳线
		o0, c0, h0, l0 := open[i-4], close[i-4], high[i-4], low[i-4]
		if !whiteBody(o0, c0) || !isLongBody(o0, c0, h0, l0) {
			continue
		}
		// 蜡烛 2-4：三根小阴线，在第一根阳线范围内，形成更低低点
		ok := true
		for k := 1; k <= 3; k++ {
			idx := i - 4 + k
			o, c, h, l := open[idx], close[idx], high[idx], low[idx]
			if !blackBody(o, c) || h > h0 || l < l0 || !isShortBody(o, c, h, l) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// 蜡烛 5：长阳线，收盘在前三根最高价之上
		o4, c4, h4, l4 := open[i], close[i], high[i], low[i]
		if !whiteBody(o4, c4) || !isLongBody(o4, c4, h4, l4) {
			continue
		}
		if c4 > h0 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLMATHOLDLookback 返回 CDLMATHOLD 的最小回溯期。
func CDLMATHOLDLookback() int { return 5 }

// CDLLADDERBOTTOM 识别 Ladder Bottom（梯底）形态。
// 一个五根蜡烛的看涨反转形态：
// 1-3. 三根收盘依次走低的阴线
// 4. 一根带有上影线的阴线
// 5. 在蜡烛 4 上方跳空的阳线
// 最小回溯：5
func CDLLADDERBOTTOM(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 4; i < n; i++ {
		// 蜡烛 1-3：三根连续的阴线，收盘降低
		if !blackBody(open[i-4], close[i-4]) || !blackBody(open[i-3], close[i-3]) || !blackBody(open[i-2], close[i-2]) {
			continue
		}
		if close[i-4] > close[i-3] && close[i-3] > close[i-2] {
			// 蜡烛 4：带有上影线的阴线
			o3, c3, h3 := open[i-1], close[i-1], high[i-1]
			if !blackBody(o3, c3) {
				continue
			}
			us3 := upperShadow(h3, o3, c3)
			if us3 <= realBody(o3, c3)*0.1 {
				continue
			}
			// 蜡烛 5：阳线，在蜡烛 4 上方跳空
			o4, c4, l4 := open[i], close[i], low[i]
			if whiteBody(o4, c4) && isGapUp(h3, l4) {
				out[i] = Bullish
			}
		}
	}
	return out, nil
}

// CDLLADDERBOTTOMLookback 返回 CDLLADDERBOTTOM 的最小回溯期。
func CDLLADDERBOTTOMLookback() int { return 5 }

// CDLRISEFALL3METHODS 识别上升/下降三法形态。
// 上升三法（看涨持续）：
//  1. 长阳线，2-4. 三根小阴线在第一根内回撤，5. 长阳线创新高
//
// 下降三法（看跌持续）：
//  1. 长阴线，2-4. 三根小阳线在第一根内回撤，5. 长阴线创新低
//
// 最小回溯：5
func CDLRISEFALL3METHODS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 4; i < n; i++ {
		// 上升三法（看涨）
		o0, c0, h0, l0 := open[i-4], close[i-4], high[i-4], low[i-4]
		if whiteBody(o0, c0) && isLongBody(o0, c0, h0, l0) {
			ok := true
			for k := 1; k <= 3; k++ {
				idx := i - 4 + k
				o, c, h, l := open[idx], close[idx], high[idx], low[idx]
				if !blackBody(o, c) || h > h0 || !isShortBody(o, c, h, l) {
					ok = false
					break
				}
			}
			if ok {
				o4, c4 := open[i], close[i]
				if whiteBody(o4, c4) && c4 > c0 {
					out[i] = Bullish
				}
			}
		}
		// 下降三法（看跌）
		if blackBody(o0, c0) && isLongBody(o0, c0, h0, l0) {
			ok := true
			for k := 1; k <= 3; k++ {
				idx := i - 4 + k
				o, c, h, l := open[idx], close[idx], high[idx], low[idx]
				if !whiteBody(o, c) || l < l0 || !isShortBody(o, c, h, l) {
					ok = false
					break
				}
			}
			if ok {
				o4, c4 := open[i], close[i]
				if blackBody(o4, c4) && c4 < c0 {
					out[i] = Bearish
				}
			}
		}
	}
	return out, nil
}

// CDLRISEFALL3METHODSLookback 返回 CDLRISEFALL3METHODS 的最小回溯期。
func CDLRISEFALL3METHODSLookback() int { return 5 }

// CDLXSIDEGAP3METHODS 识别向上/向下跳空三法形态。
// 向上跳空三法（看涨持续）：
//
//	1-2. 两根长阳线，中间有跳空
//	3. 一根部分回补跳空的阴线
//
// 向下跳空三法（看跌持续）：
//
//	1-2. 两根长阴线，中间有跳空
//	3. 一根部分回补跳空的阳线
//
// 最小回溯：3
func CDLXSIDEGAP3METHODS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		o0, c0, h0, l0 := open[i-2], close[i-2], high[i-2], low[i-2]
		o1, c1, h1, l1 := open[i-1], close[i-1], high[i-1], low[i-1]
		o2, c2 := open[i], close[i]
		// 向上跳空三法
		if whiteBody(o0, c0) && whiteBody(o1, c1) && blackBody(o2, c2) {
			if isGapUp(h0, l1) && c2 < l1 && c2 > h0 {
				out[i] = Bullish
			}
		}
		// 向下跳空三法
		if blackBody(o0, c0) && blackBody(o1, c1) && whiteBody(o2, c2) {
			if isGapDown(l0, h1) && c2 > h1 && c2 < l0 {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLXSIDEGAP3METHODSLookback 返回 CDLXSIDEGAP3METHODS 的最小回溯期。
func CDLXSIDEGAP3METHODSLookback() int { return 3 }

// ---- 内部辅助函数 ----

func mathMax3(a, b, c float64) float64 {
	if a > b {
		if a > c {
			return a
		}
		return c
	}
	if b > c {
		return b
	}
	return c
}

func mathMin3(a, b, c float64) float64 {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
