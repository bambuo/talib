package cdl

// =====================================================================
// Single Candle Patterns — 单K线形态
// =====================================================================

// CDLDOJI 识别十字星形态。
// 十字星的实体相对于蜡烛范围非常小（≤10%）。
// 最小回溯：1
func CDLDOJI(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		if isDoji(open[i], close[i], high[i], low[i]) {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLDOJILookback 返回 CDLDOJI 的最小回溯期。
func CDLDOJILookback() int { return 1 }

// CDLDRAGONFLYDOJI 识别蜻蜓十字星。
// 一种具有长下影线且无上影线的十字星，位于交易时段的高位。
// 实体必须在或接近高位（上影线 ≤ 实体，长下影线）。
// 最小回溯：1
func CDLDRAGONFLYDOJI(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if !isDoji(o, c, h, l) {
			continue
		}
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 上影线非常短，下影线显著长于实体
		if us <= body && ls > body*2 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLDRAGONFLYDOJILookback 返回 CDLDRAGONFLYDOJI 的最小回溯期。
func CDLDRAGONFLYDOJILookback() int { return 1 }

// CDLGRAVESTONEDOJI 识别墓碑十字星。
// 一种具有长上影线且无下影线的十字星，位于交易时段的低位。
// 最小回溯：1
func CDLGRAVESTONEDOJI(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if !isDoji(o, c, h, l) {
			continue
		}
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 下影线非常短，上影线显著长于实体
		if ls <= body && us > body*2 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLGRAVESTONEDOJILookback 返回 CDLGRAVESTONEDOJI 的最小回溯期。
func CDLGRAVESTONEDOJILookback() int { return 1 }

// CDLLONGLEGGEDDOJI 识别长腿十字星。
// 一种具有大致等长的长上影线和下影线的十字星。
// 最小回溯：1
func CDLLONGLEGGEDDOJI(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if !isDoji(o, c, h, l) {
			continue
		}
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		rng := candleRange(h, l)
		// 两个影线至少为实体的 3 倍，影线大致相等，且实体在中部附近
		if us > body*2 && ls > body*2 && rng > body*3 {
			// 影线之间大致成比例
			upperRatio := us / (us + ls)
			if upperRatio > 0.4 && upperRatio < 0.6 {
				out[i] = Bullish
			}
		}
	}
	return out, nil
}

// CDLLONGLEGGEDDOJILookback 返回 CDLLONGLEGGEDDOJI 的最小回溯期。
func CDLLONGLEGGEDDOJILookback() int { return 1 }

// CDLMARUBOZU 识别光脚光头（Marubozu）。
// 一种没有或只有很短影线的长蜡烛（实体几乎填满整个范围）。
// 看涨光脚光头开盘于最低价，收盘于最高价。
// 看跌光脚光头开盘于最高价，收盘于最低价。
// 最小回溯：1
func CDLMARUBOZU(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		if isMarubozu(open[i], close[i], high[i], low[i]) {
			if isBullish(open[i], close[i]) {
				out[i] = Bullish
			} else if isBearish(open[i], close[i]) {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLMARUBOZULookback 返回 CDLMARUBOZU 的最小回溯期。
func CDLMARUBOZULookback() int { return 1 }

// CDLCLOSINGMARUBOZU 识别收盘光脚光头。
// 一种没有上影线（看涨）或没有下影线（看跌）的长蜡烛。
// 与光脚光头不同，允许有一个小影线。
// 最小回溯：1
func CDLCLOSINGMARUBOZU(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		body := realBody(o, c)
		rng := candleRange(h, l)
		if rng == 0 || body < rng*0.6 {
			continue
		}
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		if isBullish(o, c) {
			// 看涨收盘光脚光头：开盘接近最低价，收盘接近最高价
			if ls <= body*0.1 && us <= body*0.1 {
				out[i] = Bullish
			}
		} else if isBearish(o, c) {
			// 看跌收盘光脚光头：开盘接近最高价，收盘接近最低价
			if us <= body*0.1 && ls <= body*0.1 {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLCLOSINGMARUBOZULookback 返回 CDLCLOSINGMARUBOZU 的最小回溯期。
func CDLCLOSINGMARUBOZULookback() int { return 1 }

// CDLLONGLINE 识别长实体线。
// 一种实体相对于范围较长的蜡烛。上下文相关：
// 在上升趋势中，白色长实体线为持续信号；在下降趋势中，黑色长实体线为持续信号。
// 最小回溯：1
func CDLLONGLINE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		if isLongBody(open[i], close[i], high[i], low[i]) {
			if isBullish(open[i], close[i]) {
				out[i] = Bullish
			} else if isBearish(open[i], close[i]) {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLLONGLINELookback 返回 CDLLONGLINE 的最小回溯期。
func CDLLONGLINELookback() int { return 1 }

// CDLSHORTLINE 识别短实体线。
// 一种实体相对于范围非常短的蜡烛。
// 最小回溯：1
func CDLSHORTLINE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		body := realBody(open[i], close[i])
		rng := candleRange(high[i], low[i])
		if rng > 0 && body < rng*0.2 {
			if isBullish(open[i], close[i]) {
				out[i] = Bullish
			} else if isBearish(open[i], close[i]) {
				out[i] = Bearish
			} else {
				out[i] = Bullish
			}
		}
	}
	return out, nil
}

// CDLSHORTLINELookback 返回 CDLSHORTLINE 的最小回溯期。
func CDLSHORTLINELookback() int { return 1 }

// CDLSPINNINGTOP 识别陀螺线。
// 一种小实体且影线大致等长的蜡烛。
// 实体 < 范围的 20%，两影线均 > 实体。
// 最小回溯：1
func CDLSPINNINGTOP(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if isDoji(o, c, h, l) {
			continue
		}
		body := realBody(o, c)
		rng := candleRange(h, l)
		if rng == 0 {
			continue
		}
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 小实体，两侧有影线，两者均长于实体
		if body < rng*0.2 && us > body && ls > body {
			if isBullish(o, c) {
				out[i] = Bullish
			} else {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLSPINNINGTOPLookback 返回 CDLSPINNINGTOP 的最小回溯期。
func CDLSPINNINGTOPLookback() int { return 1 }

// CDLHIGHWAVE 识别长脚线（High Wave）。
// 一种具有非常长的上影线和下影线以及小实体的蜡烛。
// 最小回溯：1
func CDLHIGHWAVE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if isDoji(o, c, h, l) {
			continue
		}
		body := realBody(o, c)
		rng := candleRange(h, l)
		if rng == 0 {
			continue
		}
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 实体非常小，两个影线非常长
		if body < rng*0.15 && us > body*3 && ls > body*3 {
			if isBullish(o, c) {
				out[i] = Bullish
			} else {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLHIGHWAVELookback 返回 CDLHIGHWAVE 的最小回溯期。
func CDLHIGHWAVELookback() int { return 1 }

// CDLRICKSHAWMAN 识别人力车十字星。
// 一种实体位于蜡烛范围中部的十字星。
// 最小回溯：1
func CDLRICKSHAWMAN(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 0; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if !isDoji(o, c, h, l) {
			continue
		}
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		body := realBody(o, c)
		// 实体大致在中间：影线近似相等
		totalShadow := us + ls
		if totalShadow > body*3 && us > body && ls > body {
			ratio := us / totalShadow
			if ratio > 0.3 && ratio < 0.7 {
				out[i] = Neutral
			}
		}
	}
	return out, nil
}

// CDLRICKSHAWMANLookback 返回 CDLRICKSHAWMAN 的最小回溯期。
func CDLRICKSHAWMANLookback() int { return 1 }

// CDLTAKURI 识别 Takuri 线（具有非常长下影线的蜻蜓十字星）。
// 类似于蜻蜓十字星，但下影线更加极端。
// 在下降趋势中，这是一个看涨反转信号。
// 最小回溯：4
func CDLTAKURI(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		if !isDoji(o, c, h, l) {
			continue
		}
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 非常长的下影线，几乎没有上影线
		if us <= body && ls > body*3 {
			// 下降趋势上下文：前面的蜡烛大多为阴线
			downCount := 0
			for j := i - 3; j < i; j++ {
				if isBearish(open[j], close[j]) {
					downCount++
				}
			}
			if downCount >= 2 {
				out[i] = Bullish
			}
		}
	}
	return out, nil
}

// CDLTAKURILookback 返回 CDLTAKURI 的最小回溯期。
func CDLTAKURILookback() int { return 4 }

// CDLHAMMER 识别锤子线。
// 一种小实体靠近范围顶部、具有长下影线（至少为实体的 2 倍）、
// 且几乎没有上影线的蜡烛。锤子线出现在下降趋势中，是看涨反转信号。
// 最小回溯：4
func CDLHAMMER(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 锤子线：小实体，长下影线（≥ 实体的 2 倍），极短上影线
		if us > body*0.5 || ls < body*2 || body == 0 {
			continue
		}
		// 需要下降趋势上下文
		downCount := 0
		for j := i - 3; j < i; j++ {
			if isBearish(open[j], close[j]) {
				downCount++
			}
		}
		if downCount >= 2 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLHAMMERLookback 返回 CDLHAMMER 的最小回溯期。
func CDLHAMMERLookback() int { return 4 }

// CDLHANGINGMAN 识别吊人线。
// 与锤子线相同的蜡烛形态，但出现在上升趋势中（看跌反转）。
// 最小回溯：4
func CDLHANGINGMAN(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		if us > body*0.5 || ls < body*2 || body == 0 {
			continue
		}
		// 需要上升趋势上下文
		upCount := 0
		for j := i - 3; j < i; j++ {
			if isBullish(open[j], close[j]) {
				upCount++
			}
		}
		if upCount >= 2 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLHANGINGMANLookback 返回 CDLHANGINGMAN 的最小回溯期。
func CDLHANGINGMANLookback() int { return 4 }

// CDLINVERTEDHAMMER 识别倒锤子线。
// 一种小实体靠近底部、具有长上影线（≥ 实体的 2 倍）的蜡烛。
// 出现在下降趋势中，为看涨反转信号。
// 最小回溯：4
func CDLINVERTEDHAMMER(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		// 长上影线，几乎没有下影线
		if ls > body*0.5 || us < body*2 || body == 0 {
			continue
		}
		downCount := 0
		for j := i - 3; j < i; j++ {
			if isBearish(open[j], close[j]) {
				downCount++
			}
		}
		if downCount >= 2 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLINVERTEDHAMMERLookback 返回 CDLINVERTEDHAMMER 的最小回溯期。
func CDLINVERTEDHAMMERLookback() int { return 4 }

// CDLSHOOTINGSTAR 识别流星线。
// 与倒锤子线相同的蜡烛形态，但出现在上升趋势中（看跌反转）。
// 最小回溯：4
func CDLSHOOTINGSTAR(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		body := realBody(o, c)
		us := upperShadow(h, o, c)
		ls := lowerShadow(l, o, c)
		if ls > body*0.5 || us < body*2 || body == 0 {
			continue
		}
		upCount := 0
		for j := i - 3; j < i; j++ {
			if isBullish(open[j], close[j]) {
				upCount++
			}
		}
		if upCount >= 2 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLSHOOTINGSTARLookback 返回 CDLSHOOTINGSTAR 的最小回溯期。
func CDLSHOOTINGSTARLookback() int { return 4 }
