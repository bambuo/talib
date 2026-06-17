package cdl

// =====================================================================
// Star Patterns — 星形形态
// =====================================================================

// CDLMORNINGSTAR 识别早晨之星。
// 三根蜡烛的看涨反转形态：
// 1. 一根长阴线
// 2. 一根小实体蜡烛（任意颜色），在第一根下方存在跳空
// 3. 一根长阳线，收盘深入第一根阴线实体
// 第二根蜡烛实体不应与第一根重叠。
// 最小回溯：3
func CDLMORNINGSTAR(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 蜡烛 1（前二）：长阴线
		o0, c0, h0, l0 := open[i-2], close[i-2], high[i-2], low[i-2]
		if !blackBody(o0, c0) || !isLongBody(o0, c0, h0, l0) {
			continue
		}
		// 蜡烛 2（前一）：小实体，在蜡烛 1 下方跳空
		o1, c1, h1, l1 := open[i-1], close[i-1], high[i-1], low[i-1]
		// 蜡烛 2 必须从蜡烛 1 向下跳空且为短实体
		if !isGapDown(l0, h1) || !isShortBody(o1, c1, h1, l1) {
			continue
		}
		// 蜡烛 3（当前）：长阳线，收盘进入第一根蜡烛实体
		o2, c2, h2, l2 := open[i], close[i], high[i], low[i]
		if !whiteBody(o2, c2) || !isLongBody(o2, c2, h2, l2) {
			continue
		}
		// 收盘价高于第一根阴线实体的中点
		midpoint := (o0 + c0) / 2
		if c2 > midpoint {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLMORNINGSTARLookback 返回 CDLMORNINGSTAR 的最小回溯期。
func CDLMORNINGSTARLookback() int { return 3 }

// CDLEVENINGSTAR 识别黄昏之星。
// 三根蜡烛的看跌反转形态（早晨之星的镜像）：
// 1. 一根长阳线
// 2. 一根小实体蜡烛（任意颜色），在第一根上方存在跳空
// 3. 一根长阴线，收盘深入第一根阳线实体
// 最小回溯：3
func CDLEVENINGSTAR(open, high, low, close []float64) ([]int, error) {
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
		// 蜡烛 2：小实体，在蜡烛 1 上方跳空
		o1, c1, h1, l1 := open[i-1], close[i-1], high[i-1], low[i-1]
		if !isGapUp(h0, l1) || !isShortBody(o1, c1, h1, l1) {
			continue
		}
		// 蜡烛 3：长阴线，收盘进入第一根蜡烛实体
		o2, c2, h2, l2 := open[i], close[i], high[i], low[i]
		if !blackBody(o2, c2) || !isLongBody(o2, c2, h2, l2) {
			continue
		}
		midpoint := (o0 + c0) / 2
		if c2 < midpoint {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLEVENINGSTARLookback 返回 CDLEVENINGSTAR 的最小回溯期。
func CDLEVENINGSTARLookback() int { return 3 }

// CDLMORNINGDOJISTAR 识别早晨十字星。
// 与早晨之星相同，但中间蜡烛为十字星。
// 最小回溯：3
func CDLMORNINGDOJISTAR(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 蜡烛 1：长阴线
		o0, c0, h0, l0 := open[i-2], close[i-2], high[i-2], low[i-2]
		if !blackBody(o0, c0) || !isLongBody(o0, c0, h0, l0) {
			continue
		}
		// 蜡烛 2：必须为十字星并向下跳空
		o1, c1, h1, l1 := open[i-1], close[i-1], high[i-1], low[i-1]
		if !isDoji(o1, c1, h1, l1) {
			continue
		}
		// 向下跳空：十字星在第一根蜡烛下方
		if !bodyGapDown(o0, c0, o1, c1) {
			continue
		}
		// 蜡烛 3：长阳线，收盘进入第一根蜡烛实体
		o2, c2, h2, l2 := open[i], close[i], high[i], low[i]
		if !whiteBody(o2, c2) || !isLongBody(o2, c2, h2, l2) {
			continue
		}
		midpoint := (o0 + c0) / 2
		if c2 > midpoint {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLMORNINGDOJISTARLookback 返回 CDLMORNINGDOJISTAR 的最小回溯期。
func CDLMORNINGDOJISTARLookback() int { return 3 }

// CDLEVENINGDOJISTAR 识别黄昏十字星。
// 与黄昏之星相同，但中间蜡烛为十字星。
// 最小回溯：3
func CDLEVENINGDOJISTAR(open, high, low, close []float64) ([]int, error) {
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
		// 蜡烛 2：必须为十字星并向上跳空
		o1, c1, h1, l1 := open[i-1], close[i-1], high[i-1], low[i-1]
		if !isDoji(o1, c1, h1, l1) {
			continue
		}
		if !bodyGapUp(o0, c0, o1, c1) {
			continue
		}
		// 蜡烛 3：长阴线，收盘进入第一根蜡烛实体
		o2, c2, h2, l2 := open[i], close[i], high[i], low[i]
		if !blackBody(o2, c2) || !isLongBody(o2, c2, h2, l2) {
			continue
		}
		midpoint := (o0 + c0) / 2
		if c2 < midpoint {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLEVENINGDOJISTARLookback 返回 CDLEVENINGDOJISTAR 的最小回溯期。
func CDLEVENINGDOJISTARLookback() int { return 3 }

// CDLABANDONEDBABY 识别弃婴形态。
// 一种类似于早晨/黄昏之星的三根蜡烛形态，但中间的十字星
// 在两侧均留下跳空（影线不重叠）。
// 看涨：阴线，向下跳空十字星，向上跳空阳线。
// 看跌：阳线，向上跳空十字星，向下跳空阴线。
// 最小回溯：3
func CDLABANDONEDBABY(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 蜡烛 2（中间）：必须为十字星
		if !isDoji(open[i-1], close[i-1], high[i-1], low[i-1]) {
			continue
		}
		// 看涨弃婴形态
		if blackBody(open[i-2], close[i-2]) && isLongBody(open[i-2], close[i-2], high[i-2], low[i-2]) &&
			whiteBody(open[i], close[i]) && isLongBody(open[i], close[i], high[i], low[i]) {
			// 十字星从第一根向下跳空（影线不重叠）
			if high[i-1] < low[i-2] && low[i-1] > high[i] {
				continue // 错误 - 让我重新思考
			}
			// 中间十字星的上影线在第一根蜡烛的下影线下方
			// 第三根蜡烛的下影线在中间十字星的上影线上方
			if high[i-1] < low[i-2] && low[i] > high[i-1] {
				out[i] = Bullish
			}
		}
		// 看跌弃婴形态
		if whiteBody(open[i-2], close[i-2]) && isLongBody(open[i-2], close[i-2], high[i-2], low[i-2]) &&
			blackBody(open[i], close[i]) && isLongBody(open[i], close[i], high[i], low[i]) {
			// 中间十字星向上跳空（上影线在第三根下方，下影线在第一根上方）
			if low[i-1] > high[i-2] && high[i] < low[i-1] {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLABANDONEDBABYLookback 返回 CDLABANDONEDBABY 的最小回溯期。
func CDLABANDONEDBABYLookback() int { return 3 }

// CDLTRISTAR 识别三颗星形态。
// 连续三根十字星。
// 看涨：第三根十字星在第二根上方，处于下降趋势。
// 看跌：第三根十字星在第二根下方，处于上升趋势。
// 最小回溯：3
func CDLTRISTAR(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 三根都必须为十字星
		if !isDoji(open[i-2], close[i-2], high[i-2], low[i-2]) ||
			!isDoji(open[i-1], close[i-1], high[i-1], low[i-1]) ||
			!isDoji(open[i], close[i], high[i], low[i]) {
			continue
		}
		// 中间十字星必须与第一根存在跳空
		if isGapDown(low[i-2], high[i-1]) {
			// 第一根：上方，第二根：下方（向下跳空），第三根：从第二根向上跳空
			if isGapUp(high[i-1], low[i]) {
				out[i] = Bullish
			}
		}
		if isGapUp(high[i-2], low[i-1]) {
			if isGapDown(low[i-1], high[i]) {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLTRISTARLookback 返回 CDLTRISTAR 的最小回溯期。
func CDLTRISTARLookback() int { return 3 }
