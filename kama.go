package talib

import (
	"math"
)

// KAMA 计算考夫曼自适应移动平均。
//
// KAMA 根据市场噪声（波动性）使用效率比率动态调整其平滑常数：
//
//	ER = |变化量| / 波动性
//
// 其中：
//   - 变化量 = |价格[i] - 价格[i-周期]|
//   - 波动性 = Σ|价格[j] - 价格[j-1]| 对于 j 在 [i-周期+1, i] 范围内
//   - SC = [ER × (最快SC - 最慢SC) + 最慢SC]²
//   - 最快SC = 2/(快周期+1)
//   - 最慢SC = 2/(慢周期+1)
//   - KAMA[i] = KAMA[i-1] + SC × (价格[i] - KAMA[i-1])
//
// 典型参数：周期=10, 快周期=2, 慢周期=30。
// 前 (周期) 个元素为 NaN。
func KAMA(input []float64, period int) ([]float64, error) {
	if err := ValidateNumericInput(input, "KAMA:input"); err != nil {
		return nil, err
	}
	if err := ValidatePeriod(period, len(input), "KAMA:period"); err != nil {
		return nil, err
	}

	// 根据 TA-Lib 约定的默认快/慢周期
	const fastPeriod = 2
	const slowPeriod = 30

	fastestSC := 2.0 / float64(fastPeriod+1)
	slowestSC := 2.0 / float64(slowPeriod+1)
	diffSC := fastestSC - slowestSC

	n := len(input)
	out := MakeOutput(n)

	if n <= period {
		return out, nil
	}

	// 第一个有效 KAMA 在索引 period 处以 SMA 播种
	var sum float64
	for i := 0; i <= period; i++ {
		sum += input[i]
	}
	out[period] = sum / float64(period+1)

	// 预先计算绝对差值以提高效率
	absDiffs := make([]float64, n)
	for i := 1; i < n; i++ {
		absDiffs[i] = math.Abs(input[i] - input[i-1])
	}

	for i := period + 1; i < n; i++ {
		// 变化量 = |价格 - 周期的价格|
		change := math.Abs(input[i] - input[i-period])

		// 波动性 = 窗口内绝对差值之和
		var volatility float64
		for j := i - period + 1; j <= i; j++ {
			volatility += absDiffs[j]
		}

		var sc float64
		if volatility == 0 {
			sc = slowestSC
		} else {
			er := change / volatility
			sc = er*diffSC + slowestSC
			sc = sc * sc
		}

		out[i] = out[i-1] + sc*(input[i]-out[i-1])
	}

	return out, nil
}

// KAMALookback 返回 KAMA 的 lookback。
func KAMALookback(period int) int {
	return period
}
