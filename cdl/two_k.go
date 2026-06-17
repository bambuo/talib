package cdl

import "math"

// =====================================================================
// Two-Candle Patterns — 双K线形态
// =====================================================================

// CDLENGULFING 识别吞没形态。
// 看涨：一根小阴线后跟一根大阳线，完全吞没前一根阴线实体。
// 看跌：一根小阳线后跟一根大阴线，完全吞没前一根阳线实体。
// 最小回溯：2
func CDLENGULFING(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		o, c := open[i], close[i]
		prevBody := realBody(prevO, prevC)
		body := realBody(o, c)
		if prevBody == 0 || body == 0 {
			continue
		}
		// 看涨吞没：前一根阴线，当前阳线，实体吞没前一根
		if blackBody(prevO, prevC) && whiteBody(o, c) &&
			math.Min(prevO, prevC) >= math.Min(o, c) && math.Max(prevO, prevC) <= math.Max(o, c) {
			out[i] = Bullish
		}
		// 看跌吞没：前一根阳线，当前阴线，实体吞没前一根
		if whiteBody(prevO, prevC) && blackBody(o, c) &&
			math.Min(prevO, prevC) >= math.Min(o, c) && math.Max(prevO, prevC) <= math.Max(o, c) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLENGULFINGLookback 返回 CDLENGULFING 的最小回溯期。
func CDLENGULFINGLookback() int { return 2 }

// CDLHARAMI 识别孕线形态。
// 看涨：一根大阴线后跟一根小阳线，其实体完全在前一根实体之内。
// 看跌：一根大阳线后跟一根小阴线，其实体完全在前一根实体之内。
// 最小回溯：2
func CDLHARAMI(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		o, c := open[i], close[i]
		prevBody := realBody(prevO, prevC)
		body := realBody(o, c)
		if prevBody == 0 {
			continue
		}
		// 当前实体必须在上一根实体范围之内
		currMin := math.Min(o, c)
		currMax := math.Max(o, c)
		prevMin := math.Min(prevO, prevC)
		prevMax := math.Max(prevO, prevC)
		if currMin <= prevMin || currMax >= prevMax {
			continue
		}
		// 实体大小检查：当前实体显著小于前一根
		if body >= prevBody*0.6 {
			continue
		}
		// 看涨孕线：前根阴线，当前阳线
		if blackBody(prevO, prevC) && whiteBody(o, c) {
			out[i] = Bullish
		}
		// 看跌孕线：前根阳线，当前阴线
		if whiteBody(prevO, prevC) && blackBody(o, c) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLHARAMILookback 返回 CDLHARAMI 的最小回溯期。
func CDLHARAMILookback() int { return 2 }

// CDLHARAMICROSS 识别十字孕线形态。
// 一种第二根蜡烛为十字星的孕线形态。
// 最小回溯：2
func CDLHARAMICROSS(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		o, c, h, l := open[i], close[i], high[i], low[i]
		if !isDoji(o, c, h, l) {
			continue
		}
		// 十字星实体必须在上一根实体之内
		prevMin := math.Min(prevO, prevC)
		prevMax := math.Max(prevO, prevC)
		if o < prevMin || o > prevMax {
			continue
		}
		if blackBody(prevO, prevC) {
			out[i] = Bullish
		}
		if whiteBody(prevO, prevC) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLHARAMICROSSLookback 返回 CDLHARAMICROSS 的最小回溯期。
func CDLHARAMICROSSLookback() int { return 2 }

// CDLDARKCLOUDCOVER 识别乌云盖顶形态。
// 一根阳线后跟一根阴线，开盘价高于前一根最高价，
// 但收盘价深入前一根阳线实体（≥50% 的穿透）。
// 为看跌反转形态。
// 最小回溯：2
func CDLDARKCLOUDCOVER(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevH := high[i-1]
		o, c := open[i], close[i]
		// 第一根蜡烛必须为强势阳线
		if !whiteBody(prevO, prevC) || realBody(prevO, prevC) < candleRange(prevH, low[i-1])*0.5 {
			continue
		}
		// 第二根蜡烛：阴线，开盘价高于前一根最高价
		if !blackBody(o, c) || o <= prevH {
			continue
		}
		// 收盘价进入前一根阳线实体（低于中点）
		midpoint := (prevO + prevC) / 2
		if c >= midpoint {
			continue
		}
		out[i] = Bearish
	}
	return out, nil
}

// CDLDARKCLOUDCOVERLookback 返回 CDLDARKCLOUDCOVER 的最小回溯期。
func CDLDARKCLOUDCOVERLookback() int { return 2 }

// CDLPIERCING 识别刺透形态。
// 一根阴线后跟一根阳线，开盘价低于前一根最低价，
// 但收盘价深入前一根阴线实体（>50% 的穿透）。
// 为看涨反转形态。
// 最小回溯：2
func CDLPIERCING(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevL := low[i-1]
		o, c := open[i], close[i]
		// 第一根蜡烛必须为强势阴线
		if !blackBody(prevO, prevC) || realBody(prevO, prevC) < candleRange(high[i-1], prevL)*0.5 {
			continue
		}
		// 第二根蜡烛：阳线，开盘价低于前一根最低价
		if !whiteBody(o, c) || o >= prevL {
			continue
		}
		// 收盘价进入前一根阴线实体（高于中点）
		midpoint := (prevO + prevC) / 2
		if c <= midpoint {
			continue
		}
		out[i] = Bullish
	}
	return out, nil
}

// CDLPIERCINGLookback 返回 CDLPIERCING 的最小回溯期。
func CDLPIERCINGLookback() int { return 2 }

// CDLCOUNTERATTACK 识别反击线形态。
// 两根颜色相反的蜡烛，具有相同（或相近）的收盘价。
// 看涨：阴线后阳线，收盘价相等。看跌：阳线后阴线，收盘价相等。
// 最小回溯：2
func CDLCOUNTERATTACK(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		o, c := open[i], close[i]
		// 收盘价必须非常接近
		if math.Abs(c-prevC) > math.Abs(c-o)*0.1 {
			continue
		}
		if blackBody(prevO, prevC) && whiteBody(o, c) {
			out[i] = Bullish
		}
		if whiteBody(prevO, prevC) && blackBody(o, c) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLCOUNTERATTACKLookback 返回 CDLCOUNTERATTACK 的最小回溯期。
func CDLCOUNTERATTACKLookback() int { return 2 }

// CDLKICKING 识别踢腿形态。
// 一种颜色的光脚光头后跟一个跳空及相反颜色的光脚光头。
// 看涨：黑色光脚光头，向下跳空，白色光脚光头。
// 看跌：白色光脚光头，向上跳空，黑色光脚光头。
// 最小回溯：2
func CDLKICKING(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevH, prevL := high[i-1], low[i-1]
		o, c := open[i], close[i]
		h, l := high[i], low[i]
		if !isMarubozu(prevO, prevC, prevH, prevL) || !isMarubozu(o, c, h, l) {
			continue
		}
		// 看涨踢腿：黑色光脚光头，向下跳空，白色光脚光头
		if blackBody(prevO, prevC) && whiteBody(o, c) && isGapDown(prevL, h) {
			out[i] = Bullish
		}
		// 看跌踢腿：白色光脚光头，向上跳空，黑色光脚光头
		if whiteBody(prevO, prevC) && blackBody(o, c) && isGapUp(prevH, l) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLKICKINGLookback 返回 CDLKICKING 的最小回溯期。
func CDLKICKINGLookback() int { return 2 }

// CDLKICKINGBYLENGTH 识别按长度确定方向的踢腿形态。
// 类似于踢腿形态，但信号方向由两根蜡烛中较长者决定。
// 最小回溯：2
func CDLKICKINGBYLENGTH(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevH, prevL := high[i-1], low[i-1]
		o, c := open[i], close[i]
		h, l := high[i], low[i]
		if !isMarubozu(prevO, prevC, prevH, prevL) || !isMarubozu(o, c, h, l) {
			continue
		}
		// 相反颜色且存在跳空
		if blackBody(prevO, prevC) && whiteBody(o, c) && isGapDown(prevL, h) {
			out[i] = Bullish
		}
		if whiteBody(prevO, prevC) && blackBody(o, c) && isGapUp(prevH, l) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLKICKINGBYLENGTHLookback 返回 CDLKICKINGBYLENGTH 的最小回溯期。
func CDLKICKINGBYLENGTHLookback() int { return 2 }

// CDLGAPSIDESIDEWHITE 识别向上/向下跳空并列阳线。
// 跳空后两根收盘价大致相等的阳线。
// 向上跳空：看涨持续。向下跳空：看跌持续。
// 最小回溯：3
func CDLGAPSIDESIDEWHITE(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		// 前两根蜡烛必须为阳线
		if !whiteBody(open[i-1], close[i-1]) || !whiteBody(open[i], close[i]) {
			continue
		}
		// 收盘价大致相等
		if math.Abs(close[i]-close[i-1]) > realBody(open[i], close[i])*0.3 {
			continue
		}
		// 检查从蜡烛 i-2 到 i-1 的跳空方向
		prevL, prevH := low[i-2], high[i-2]
		currL, currH := low[i-1], high[i-1]
		if isGapUp(prevH, currL) {
			out[i] = Bullish
		} else if isGapDown(prevL, currH) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLGAPSIDESIDEWHITELookback 返回 CDLGAPSIDESIDEWHITE 的最小回溯期。
func CDLGAPSIDESIDEWHITELookback() int { return 3 }

// CDLDOJISTAR 识别十字星形态。
// 一根与前一根蜡烛存在跳空的十字星。前一根蜡烛应有长实体。
// 看涨：阴线后向下跳空。看跌：阳线后向上跳空。
// 最小回溯：2
func CDLDOJISTAR(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		if !isDoji(open[i], close[i], high[i], low[i]) {
			continue
		}
		prevO, prevC := open[i-1], close[i-1]
		prevH, prevL := high[i-1], low[i-1]
		// 前一根蜡烛必须有长实体
		if !isLongBody(prevO, prevC, prevH, prevL) {
			continue
		}
		// 阴线后向下跳空 = 看涨
		if blackBody(prevO, prevC) && isGapDown(prevL, high[i]) {
			out[i] = Bullish
		}
		// 阳线后向上跳空 = 看跌
		if whiteBody(prevO, prevC) && isGapUp(prevH, low[i]) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLDOJISTARLookback 返回 CDLDOJISTAR 的最小回溯期。
func CDLDOJISTARLookback() int { return 2 }

// CDLBELTHOLD 识别捉腰带线。
// 看涨：一根长白色光脚光头，开盘于或接近最低价，收盘于接近最高价，
// 通常出现在下降趋势中。
// 看跌：一根长黑色光脚光头，开盘于或接近最高价，收盘于接近最低价，出现在上升趋势中。
// 最小回溯：4
func CDLBELTHOLD(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 3; i < n; i++ {
		o, c, h, l := open[i], close[i], high[i], low[i]
		body := realBody(o, c)
		rng := candleRange(h, l)
		if rng == 0 || body < rng*0.7 {
			continue
		}
		// 看涨捉腰带线：阳线，开盘接近最低价
		if whiteBody(o, c) && lowerShadow(l, o, c) <= body*0.1 {
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
		// 看跌捉腰带线：阴线，开盘接近最高价
		if blackBody(o, c) && upperShadow(h, o, c) <= body*0.1 {
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
	}
	return out, nil
}

// CDLBELTHOLDLookback 返回 CDLBELTHOLD 的最小回溯期。
func CDLBELTHOLDLookback() int { return 4 }

// CDLHOMINGPIGEON 识别家鸽形态。
// 两根阴线，第二根完全在第一根的范围内。
// 下降趋势中的看涨反转。
// 最小回溯：3
func CDLHOMINGPIGEON(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevH, prevL := high[i-1], low[i-1]
		o, c, h, l := open[i], close[i], high[i], low[i]
		// 两根蜡烛都必须为阴线
		if !blackBody(prevO, prevC) || !blackBody(o, c) {
			continue
		}
		// 第二根蜡烛必须在第一根蜡烛范围内
		if o >= prevO || c <= prevC || h > prevH || l < prevL {
			continue
		}
		// 需要下降趋势上下文
		downCount := 0
		for j := i - 2; j < i; j++ {
			if isBearish(open[j], close[j]) {
				downCount++
			}
		}
		if downCount >= 1 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLHOMINGPIGEONLookback 返回 CDLHOMINGPIGEON 的最小回溯期。
func CDLHOMINGPIGEONLookback() int { return 3 }

// CDLMATCHINGLOW 识别低档并列阴线。
// 两根收盘价相同（或非常接近）的阴线。
// 通常为看涨反转信号。
// 最小回溯：2
func CDLMATCHINGLOW(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		o, c := open[i], close[i]
		if !blackBody(prevO, prevC) || !blackBody(o, c) {
			continue
		}
		// 收盘价近似相等
		if math.Abs(prevC-c) <= math.Abs(prevO-prevC)*0.1 {
			out[i] = Bullish
		}
	}
	return out, nil
}

// CDLMATCHINGLOWLookback 返回 CDLMATCHINGLOW 的最小回溯期。
func CDLMATCHINGLOWLookback() int { return 2 }

// CDLSEPARATINGLINES 识别分离线形态。
// 两根颜色相同的蜡烛，开盘价相同。
// 看涨：两根阳线。看跌：两根阴线。
// 最小回溯：2
func CDLSEPARATINGLINES(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		o, c := open[i], close[i]
		// 开盘价必须接近相等
		if math.Abs(prevO-o) > math.Abs(prevC-prevO)*0.1 {
			continue
		}
		if whiteBody(prevO, prevC) && whiteBody(o, c) {
			out[i] = Bullish
		}
		if blackBody(prevO, prevC) && blackBody(o, c) {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLSEPARATINGLINESLookback 返回 CDLSEPARATINGLINES 的最小回溯期。
func CDLSEPARATINGLINESLookback() int { return 2 }

// CDLSTICKSANDWICH 识别夹棍形态。
// 两根阳线夹着一根阴线，且收盘价相等。
// 看涨反转形态。
// 最小回溯：3
func CDLSTICKSANDWICH(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		o0, c0 := open[i-2], close[i-2]
		o1, c1 := open[i-1], close[i-1]
		o2, c2 := open[i], close[i]
		// 阳线、阴线、阳线
		if !whiteBody(o0, c0) || !blackBody(o1, c1) || !whiteBody(o2, c2) {
			continue
		}
		// 两根阳线的收盘价相等
		if math.Abs(c0-c2) > math.Abs(c0-o0)*0.1 {
			continue
		}
		out[i] = Bullish
	}
	return out, nil
}

// CDLSTICKSANDWICHLookback 返回 CDLSTICKSANDWICH 的最小回溯期。
func CDLSTICKSANDWICHLookback() int { return 3 }

// CDLTASUKIGAP 识别 Tasuki 跳空形态。
// 趋势方向上的跳空，后跟一根收盘价进入跳空区域的蜡烛。
// 看涨（Tasuki 向上跳空）：阳线，向上跳空，然后阴线回补跳空。
// 看跌（Tasuki 向下跳空）：阴线，向下跳空，然后阳线回补跳空。
// 最小回溯：3
func CDLTASUKIGAP(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 2; i < n; i++ {
		prevO, prevC := open[i-2], close[i-2]
		prevH, prevL := high[i-2], low[i-2]
		midO, midC := open[i-1], close[i-1]
		midH, midL := high[i-1], low[i-1]
		o, c := open[i], close[i]
		// 向上 Tasuki 跳空（看涨持续）
		if whiteBody(prevO, prevC) && isGapUp(prevH, midL) && blackBody(o, c) {
			// 第三根蜡烛开盘在第二根蜡烛实体内，收盘进入跳空区域
			if o > midC && o < midO && c > prevC && c < midL {
				out[i] = Bullish
			}
		}
		// 向下 Tasuki 跳空（看跌持续）
		if blackBody(prevO, prevC) && isGapDown(prevL, midH) && whiteBody(o, c) {
			if o < midC && o > midO && c < prevC && c > midH {
				out[i] = Bearish
			}
		}
	}
	return out, nil
}

// CDLTASUKIGAPLookback 返回 CDLTASUKIGAP 的最小回溯期。
func CDLTASUKIGAPLookback() int { return 3 }

// CDLONNECK 识别颈上线形态。
// 一根阴线后跟一根阳线，开盘价低于最低价，收盘价在或接近前一根收盘价。
// 看跌持续形态。
// 最小回溯：2
func CDLONNECK(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevL := low[i-1]
		o, c := open[i], close[i]
		if !blackBody(prevO, prevC) || !whiteBody(o, c) {
			continue
		}
		// 开盘价低于前一根最低价，收盘价在/接近前一根收盘价
		if o < prevL && math.Abs(c-prevC) <= realBody(o, c)*0.1 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLONNECKLookback 返回 CDLONNECK 的最小回溯期。
func CDLONNECKLookback() int { return 2 }

// CDLINNECK 识别颈内线形态。
// 一根阴线后跟一根阳线，开盘价低于最低价，收盘价略微进入前一根阴线实体。
// 看跌持续形态。
// 最小回溯：2
func CDLINNECK(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevL := low[i-1]
		o, c := open[i], close[i]
		if !blackBody(prevO, prevC) || !whiteBody(o, c) {
			continue
		}
		// 开盘价低于前一根最低价，收盘价仅略高于前一根收盘价
		if o < prevL && math.Abs(c-prevL) <= realBody(o, c)*0.2 {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLINNECKLookback 返回 CDLINNECK 的最小回溯期。
func CDLINNECKLookback() int { return 2 }

// CDLTHRUSTING 识别插入线形态。
// 一根阴线后跟一根阳线，开盘价低于最低价，收盘价进入前一根阴线实体但不足 50%。
// 看跌持续形态。
// 最小回溯：2
func CDLTHRUSTING(open, high, low, close []float64) ([]int, error) {
	if err := validateOHLC(open, high, low, close); err != nil {
		return nil, err
	}
	n := len(open)
	out := makeOutput(n)
	for i := 1; i < n; i++ {
		prevO, prevC := open[i-1], close[i-1]
		prevL := low[i-1]
		o, c := open[i], close[i]
		if !blackBody(prevO, prevC) || !whiteBody(o, c) {
			continue
		}
		// 开盘价低于前一根最低价，收盘价进入阴线实体但低于中点
		midpoint := (prevO + prevC) / 2
		if o < prevL && c > prevL && c < midpoint {
			out[i] = Bearish
		}
	}
	return out, nil
}

// CDLTHRUSTINGLookback returns the minimum lookback for CDLTHRUSTING.
func CDLTHRUSTINGLookback() int { return 2 }
