package talib

import (
	"math"
)

// HTResult 包含所有 HT 指标使用的 Hilbert 变换内部状态。
type HTResult struct {
	Detrender  []float64 // 去趋势价格
	Q1         []float64 // 正交分量
	I1         []float64 // 同相分量
	JQ         []float64 // 平滑后的正交（Q1 的 EMA）
	JI         []float64 // 平滑后的同相（I1 的 EMA）
	JPhase     []float64 // 相位角
	Quadrature []float64 // 正交（相位差）
}

// HT_TRENDLINE 计算 Hilbert 变换瞬时趋势线。
//
// 趋势线计算为同相分量的加权版本，
// 使用平滑后 I1 的简单移动平均。
func HT_TRENDLINE(input []float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "HT_TRENDLINE:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	// 趋势线 = SMA(JI, 3) 带平滑
	startIdx := 13 // 7 用于去趋势预热 + 6 更多用于平滑
	for i := startIdx; i < n; i++ {
		if IsNaN(ht.JI[i]) {
			continue
		}
		if i >= 2 {
			out[i] = (ht.JI[i] + ht.JI[i-1] + ht.JI[i-2]) / 3.0
		}
	}

	return out, nil
}

// HT_TRENDLINE_Lookback 返回 HT_TRENDLINE 的 lookback。
func HT_TRENDLINELookback() int {
	return 14
}

// HT_TRENDMODE 计算 Hilbert 变换趋势/周期模式。
//
// 模式线 = 仅趋势周期分量
// 返回 +1 表示趋势模式，0 表示周期模式。
func HT_TRENDMODE(input []float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "HT_TRENDMODE:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	startIdx := 13
	for i := startIdx; i < n; i++ {
		if IsNaN(ht.Quadrature[i]) {
			continue
		}

		// 趋势模式：对正交分量应用 SMA
		var trend float64
		if IsNaN(ht.Detrender[i]) {
			continue
		}
		rawTrend := ht.Quadrature[i]
		// 使用 3 周期 SMA 平滑
		if i >= 2 && !IsNaN(out[i-1]) {
			trend = (rawTrend + ht.Quadrature[i-1] + ht.Quadrature[i-2]) / 3.0
		} else {
			trend = rawTrend
		}

		// 模式信号：+1 表示正趋势，0 表示其他
		if trend > 0 {
			out[i] = 1.0
		} else {
			out[i] = 0.0
		}
	}

	return out, nil
}

// HT_TRENDMODELookback 返回 HT_TRENDMODE 的 lookback。
func HT_TRENDMODELookback() int {
	return 14
}

// HT_DCPERIOD 计算 Hilbert 变换主导周期。
//
// 使用相位差确定主导周期。
func HT_DCPERIOD(input []float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "HT_DCPERIOD:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	startIdx := 13
	for i := startIdx; i < n; i++ {
		if IsNaN(ht.Quadrature[i]) || IsNaN(ht.JPhase[i]) {
			continue
		}

		deltaPhase := ht.Quadrature[i]
		period := ht.JPhase[i]

		// 稳定周期计算
		const maxPeriod = 50.0
		const minPeriod = 6.0

		if deltaPhase > 0.1 {
			period = 6.2832 / deltaPhase
			if period > maxPeriod {
				period = maxPeriod
			}
			if period < minPeriod {
				period = minPeriod
			}
		}

		out[i] = period
	}

	return out, nil
}

// HT_DCPERIODLookback 返回 HT_DCPERIOD 的 lookback。
func HT_DCPERIODLookback() int {
	return 14
}

// HT_DCPHASE 计算 Hilbert 变换主导周期相位。
//
// 相位从平滑后的 I1 和 Q1 分量推导。
func HT_DCPHASE(input []float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "HT_DCPHASE:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	startIdx := 13
	for i := startIdx; i < n; i++ {
		if IsNaN(ht.JI[i]) || IsNaN(ht.JQ[i]) {
			continue
		}
		// 相位 = atan2(JI, JQ) 或 atan2(虚数, 实数)，取决于约定
		phase := math.Atan2(ht.JQ[i], ht.JI[i]) * 180.0 / math.Pi
		out[i] = phase
	}

	return out, nil
}

// HT_DCPHASELookback 返回 HT_DCPHASE 的 lookback。
func HT_DCPHASELookback() int {
	return 14
}

// HTPhasorResult 包含同相和正交输出。
type HTPhasorResult struct {
	InPhase    []float64 // 平滑后的同相分量（JI）
	Quadrature []float64 // 平滑后的正交分量（JQ）
}

// HT_PHASOR 计算 Hilbert 变换相量分量。
func HT_PHASOR(input []float64) (*HTPhasorResult, error) {
	if err := ValidateNumericInput(input, "HT_PHASOR:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	return &HTPhasorResult{
		InPhase:    ht.JI,
		Quadrature: ht.JQ,
	}, nil
}

// HT_PHASORLookback 返回 HT_PHASOR 的 lookback。
func HT_PHASORLookback() int {
	return 14
}

// HTSineResult 包含正弦和超前正弦输出。
type HTSineResult struct {
	Sine     []float64 // DC 相位的正弦值
	LeadSine []float64 // 超前正弦（相位 + 45° 的正弦值）
}

// HT_SINE 计算 Hilbert 变换正弦波。
func HT_SINE(input []float64) (*HTSineResult, error) {
	if err := ValidateNumericInput(input, "HT_SINE:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	n := len(input)
	sine := MakeOutput(n)
	leadSine := MakeOutput(n)

	startIdx := 13
	for i := startIdx; i < n; i++ {
		if IsNaN(ht.JPhase[i]) {
			continue
		}

		phaseRad := ht.JPhase[i] * math.Pi / 180.0
		sine[i] = math.Sin(phaseRad)
		leadSine[i] = math.Sin(phaseRad + math.Pi/4.0)
	}

	return &HTSineResult{
		Sine:     sine,
		LeadSine: leadSine,
	}, nil
}

// HT_SINELookback 返回 HT_SINE 的 lookback。
func HT_SINELookback() int {
	return 14
}

// HT_DCCOMPONENT 计算 Hilbert 变换主导周期分量。
//
// DC 分量是 Hilbert 变换平滑后的"趋势"，
// 在此定义为实部（同相）平滑分量的 6 周期 SMA。
func HT_DCCOMPONENT(input []float64) ([]float64, error) {
	if err := ValidateNumericInput(input, "HT_DCCOMPONENT:input"); err != nil {
		return nil, err
	}

	ht, err := computeHT(input)
	if err != nil {
		return nil, err
	}

	n := len(input)
	out := MakeOutput(n)

	startIdx := 13
	for i := startIdx; i < n; i++ {
		if i < 6 || IsNaN(ht.JI[i]) {
			continue
		}
		// 实部（同相）分量的 6 周期 SMA
		var sum float64
		for j := i - 5; j <= i; j++ {
			sum += ht.JI[j]
		}
		out[i] = sum / 6.0
	}

	return out, nil
}

// HT_DCCOMPONENTLookback 返回 HT_DCCOMPONENT 的 lookback。
func HT_DCCOMPONENTLookback() int {
	return 14
}

// computeHT 执行所有 HT 指标使用的核心 Hilbert 变换计算。
//
// 这实现了来自 TA-Lib 的 ht_utility.c 的数字 Hilbert 变换滤波器。
// 算法遵循 John Ehlers 的方法：
//
// 1. 使用特定的 IIR 滤波器对价格数据进行去趋势
// 2. 使用 HT 滤波器系数计算正交分量 Q1
// 3. 同相分量 I1 是延迟的价格
// 4. 用 EMA 平滑 Q1 和 I1
// 5. 从平滑分量计算相位和周期
func computeHT(input []float64) (*HTResult, error) {
	n := len(input)
	if n < 7 {
		return nil, ErrInputLengthMismatch("HT: input too short, need at least 7")
	}

	// 初始化输出数组
	detrender := MakeOutput(n)
	q1 := MakeOutput(n)
	i1 := MakeOutput(n)
	jQ := MakeOutput(n)
	jI := MakeOutput(n)
	jPhase := MakeOutput(n)
	quadrature := MakeOutput(n)

	// 预热：从索引 6 开始处理
	wkPhase := 0.0
	var wkI1, wkQ1 float64

	// EMA 平滑系数
	alpha1 := 0.0962
	alpha2 := 0.5769

	for i := 6; i < n; i++ {
		// 第1步：去趋势
		// 去趋势 = 0.0962*input[i] + 0.5769*input[i-2] - 0.5769*input[i-4] - 0.0962*input[i-6]
		detrender[i] = alpha1*input[i] + alpha2*input[i-2] -
			alpha2*input[i-4] - alpha1*input[i-6]

		// 第2步：Q1
		// Q1 = 0.0962*Detrender[i] + 0.5769*Detrender[i-2] - 0.5769*Detrender[i-4] - 0.0962*Detrender[i-6]
		if !IsNaN(detrender[i]) && i >= 12 &&
			!IsNaN(detrender[i-2]) && !IsNaN(detrender[i-4]) &&
			!IsNaN(detrender[i-6]) {
			q1[i] = alpha1*detrender[i] + alpha2*detrender[i-2] -
				alpha2*detrender[i-4] - alpha1*detrender[i-6]
		}

		// 第3步：I1 = Detrender[i-3]
		if i >= 9 && !IsNaN(detrender[i-3]) {
			i1[i] = detrender[i-3]
		}

		// 平滑（I1 使用 alpha=0.2 的 EMA，Q1 使用 alpha=0.33 的 EMA）
		wkI1 = i1[i]
		wkQ1 = q1[i]

		if IsNaN(jI[i-1]) {
			jI[i] = wkI1
			jQ[i] = wkQ1
		} else {
			// EMA 平滑
			if !IsNaN(wkI1) {
				jI[i] = 0.2*wkI1 + 0.8*jI[i-1]
			} else {
				jI[i] = jI[i-1]
			}
			if !IsNaN(wkQ1) {
				jQ[i] = 0.33*wkQ1 + 0.67*jQ[i-1]
			} else {
				jQ[i] = jQ[i-1]
			}
		}

		// 相位 = atan2(JQ, JI)
		if !IsNaN(jI[i]) && !IsNaN(jQ[i]) {
			jPhase[i] = math.Atan2(jQ[i], jI[i]) * 180.0 / math.Pi
			if jPhase[i] < 0 {
				jPhase[i] += 360.0
			}
		}

		// 相位差 = 正交分量（相位差）
		if !IsNaN(jPhase[i]) {
			deltaPhase := jPhase[i] - wkPhase
			wkPhase = jPhase[i]

			// 处理相位环绕
			if deltaPhase < -180 {
				deltaPhase += 360
			} else if deltaPhase > 180 {
				deltaPhase -= 360
			}
			quadrature[i] = deltaPhase * math.Pi / 180.0
		}
	}

	return &HTResult{
		Detrender:  detrender,
		Q1:         q1,
		I1:         i1,
		JQ:         jQ,
		JI:         jI,
		JPhase:     jPhase,
		Quadrature: quadrature,
	}, nil
}
