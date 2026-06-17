package talib

import (
	"math"
)

// MAMAResult 包含 MAMA 和 FAMA 线。
type MAMAResult struct {
	MAMA []float64 // MESA 自适应移动平均（快线）
	FAMA []float64 // 跟随自适应移动平均（慢线）
}

// MAMA 计算 MESA 自适应移动平均。
//
// MAMA 使用基于 Hilbert 变换的零差鉴别器来估计
// 主导周期，然后相应地调整平滑常数 alpha。
// 这创建了自动适应市场周期的移动平均。
//
// 参数：
//   - fastLimit：平滑常数上限（通常为 0.5）
//   - slowLimit：平滑常数下限（通常为 0.05）
//
// 输出包括：
//   - MAMA：快自适应线
//   - FAMA：更慢更平滑的跟随线
//
// 由 John Ehlers 在《Rocket Science for Traders》中发表。
func MAMA(input []float64, fastLimit, slowLimit float64) (*MAMAResult, error) {
	if err := ValidateNumericInput(input, "MAMA:input"); err != nil {
		return nil, err
	}

	n := len(input)
	mama := MakeOutput(n)
	fama := MakeOutput(n)

	if n < 7 {
		return &MAMAResult{MAMA: mama, FAMA: fama}, nil
	}

	// HT 滤波器系数
	const alpha1 = 0.0962
	const alpha2 = 0.5769

	// 状态变量
	var (
		detrender     [7]float64 // circular buffer for detrender
		q1Buf         [7]float64 // circular buffer for Q1
		prevDetrender float64
		prevQ1        float64
		smoothPrev    float64
		mamaPrev      float64
		famaPrev      float64
		phasePrev     float64
		periodPrev    float64 = 0.0 // Initialize as 0
		first                 = true
	)

	// 从索引 6 开始计算（HT 至少需要 7 个数据点）
	for i := 6; i < n; i++ {
		// 第1步：平滑价格
		var smooth float64
		if i >= 3 {
			smooth = (4.0*input[i] + 3.0*input[i-1] + 2.0*input[i-2] + input[i-3]) / 10.0
		} else {
			smooth = input[i]
		}
		smoothPrev = smooth

		// 第2步：去趋势
		detrender[i%7] = alpha1*smooth + alpha2*input[maxInt(i-2, 0)] -
			alpha2*input[maxInt(i-4, 0)] - alpha1*input[maxInt(i-6, 0)]

		// 在计算 Q1 之前需要 6 个去趋势值
		if i < 12 {
			continue
		}

		// 第3步：Q1
		q1Idx := i % 7
		q1Buf[q1Idx] = alpha1*detrender[i%7] + alpha2*detrender[(i-2)%7] -
			alpha2*detrender[(i-4)%7] - alpha1*detrender[(i-6)%7]

		// 第4步：I1 = 延迟的去趋势
		i1 := detrender[(i-3)%7]

		// 第5步：计算相位
		deltaPhase := 0.0
		if !first {
			q1 := q1Buf[q1Idx]
			// 计算相位角差
			// 使用零差鉴别器方法
			re := i1*prevDetrender + q1*prevQ1
			im := i1*prevQ1 - q1*prevDetrender

			if math.Abs(re) > 1e-10 || math.Abs(im) > 1e-10 {
				phase := math.Atan2(im, re) * 180.0 / math.Pi
				deltaPhase = phasePrev - phase

				// 处理相位环绕
				if deltaPhase < -180 {
					deltaPhase += 360
				} else if deltaPhase > 180 {
					deltaPhase -= 360
				}

				phasePrev = phase
			}
		} else {
			first = false
			// 初始化相位
			phasePrev = 0
		}

		// 存储上一次迭代的值
		prevDetrender = detrender[i%7]
		prevQ1 = q1Buf[q1Idx]

		if i < 13 {
			continue
		}

		// 第6步：计算瞬时周期
		deltaPhase = math.Abs(deltaPhase)
		if deltaPhase < 1.0 {
			deltaPhase = 1.0
		}

		period := 360.0 / deltaPhase

		// 限制周期范围
		if period > 1.5*periodPrev && periodPrev > 0 {
			period = 1.5 * periodPrev
		}
		if period < 0.67*periodPrev && periodPrev > 0 {
			period = 0.67 * periodPrev
		}
		if period < 6.0 {
			period = 6.0
		}
		if period > 50.0 {
			period = 50.0
		}

		// 平滑周期
		if periodPrev == 0 {
			periodPrev = period
		} else {
			period = 0.2*period + 0.8*periodPrev
		}
		periodPrev = period

		// 第7步：根据周期计算 alpha
		alpha := fastLimit / deltaPhase
		if alpha < slowLimit {
			alpha = slowLimit
		}
		if alpha > fastLimit {
			alpha = fastLimit
		}

		// 第8步：计算 MAMA 和 FAMA
		if i < 13 {
			mamaPrev = smoothPrev
			famaPrev = smoothPrev
		}

		mamaVal := alpha*smoothPrev + (1.0-alpha)*mamaPrev
		famaVal := 0.5*alpha*mamaVal + (1.0-0.5*alpha)*famaPrev

		mama[i] = mamaVal
		fama[i] = famaVal

		mamaPrev = mamaVal
		famaPrev = famaVal
	}

	return &MAMAResult{
		MAMA: mama,
		FAMA: fama,
	}, nil
}

// MAMALookback 返回 MAMA 的 lookback。
func MAMALookback() int {
	return 13
}

// maxInt 返回两个整数中的最大值。
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
