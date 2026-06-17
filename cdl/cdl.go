// Package cdl 提供了日本蜡烛图形态识别的纯 Go 实现，
// 功能上等同于 TA-Lib 的 Pattern Recognition 函数族（CDL* 系列）。
//
// 所有函数接受 []float64 类型的 OHLC 数据，返回 []int，其中：
//
//	+100 = 检测到看涨形态
//	-100 = 检测到看跌形态
//	   0 = 无信号
//
// 前导值补零以匹配 TA-Lib 输出长度约定。
// 每个函数都有对应的 Lookback() 辅助函数，返回第一个有效输出前所需的最小蜡烛数量。
//
// 精度目标：所有形态的信号与 TA-Lib C 实现一致。
package cdl

import (
	"errors"
	"math"
)

// Signal 定义了形态识别的输出信号值。
const (
	Bullish = 100  // 看涨信号
	Bearish = -100 // 看跌信号
	Neutral = 0    // 无信号
)

// Common 定义了常见错误。
var (
	ErrNilInput       = errors.New("cdl: input is nil")
	ErrEmptyInput     = errors.New("cdl: input is empty")
	ErrLengthMismatch = errors.New("cdl: OHLC slices must have equal length")
)

// validateOHLC 检查 OHLC 切片是否有效且长度相等。
func validateOHLC(open, high, low, close []float64) error {
	if open == nil || high == nil || low == nil || close == nil {
		return ErrNilInput
	}
	n := len(open)
	if n == 0 || len(high) == 0 || len(low) == 0 || len(close) == 0 {
		return ErrEmptyInput
	}
	if len(high) != n || len(low) != n || len(close) != n {
		return ErrLengthMismatch
	}
	return nil
}

// makeOutput 创建零值初始化的输出切片。
func makeOutput(length int) []int {
	return make([]int, length)
}

// ---- 蜡烛图解剖辅助函数 ----

// realBody 返回绝对实体大小：|收盘价 - 开盘价|。
func realBody(o, c float64) float64 {
	return math.Abs(c - o)
}

// candleRange 返回完整的蜡烛范围：最高价 - 最低价。
func candleRange(h, l float64) float64 {
	return h - l
}

// upperShadow 返回上影线长度。
func upperShadow(h, o, c float64) float64 {
	return h - math.Max(o, c)
}

// lowerShadow 返回下影线长度。
func lowerShadow(l, o, c float64) float64 {
	return math.Min(o, c) - l
}

// isBullish 如果收盘价高于开盘价则返回 true。
func isBullish(o, c float64) bool {
	return c > o
}

// isBearish 如果收盘价低于开盘价则返回 true。
func isBearish(o, c float64) bool {
	return c < o
}

// ---- 蜡烛图分类辅助函数 ----

// isDoji 检查实体相对于范围是否可忽略。
// 使用实体 <= 范围的 10% 作为阈值（与 TA-Lib 惯例一致）。
func isDoji(o, c, h, l float64) bool {
	rng := candleRange(h, l)
	if rng == 0 {
		return true
	}
	return realBody(o, c) <= rng*0.1
}

// isLongBody 检查实体是否显著（> 范围的 50%）。
func isLongBody(o, c, h, l float64) bool {
	rng := candleRange(h, l)
	if rng == 0 {
		return false
	}
	return realBody(o, c) > rng*0.5
}

// isShortBody 检查实体是否较小（< 范围的 20%）。
func isShortBody(o, c, h, l float64) bool {
	rng := candleRange(h, l)
	if rng == 0 {
		return true
	}
	return realBody(o, c) < rng*0.2
}

// isMarubozu 检查蜡烛的影线是否可忽略。
// 实体 >= 范围的 80% 且每个影线 <= 范围的 10%。
func isMarubozu(o, c, h, l float64) bool {
	rng := candleRange(h, l)
	if rng == 0 {
		return true
	}
	body := realBody(o, c)
	if body < rng*0.8 {
		return false
	}
	us := upperShadow(h, o, c)
	ls := lowerShadow(l, o, c)
	return us <= rng*0.1 && ls <= rng*0.1
}

// ---- 上下文辅助函数（基于周期） ----

// avgRealBody 计算一个蜡烛窗口内的平均实体大小。
func avgRealBody(open, close []float64, start, end int) float64 {
	if start >= end {
		return 0
	}
	var sum float64
	for i := start; i < end; i++ {
		sum += realBody(open[i], close[i])
	}
	return sum / float64(end-start)
}

// ---- 跳空检测辅助函数 ----

// isGapUp 检查当前是否开盘在前一根蜡烛范围之上（向上跳空）。
func isGapUp(prevHigh, currLow float64) bool {
	return currLow > prevHigh
}

// isGapDown 检查当前是否开盘在前一根蜡烛范围之下（向下跳空）。
func isGapDown(prevLow, currHigh float64) bool {
	return currHigh < prevLow
}

// bodyGapUp 检查实体（开盘-收盘范围）是否向上分离。
func bodyGapUp(prevO, prevC, currO, currC float64) bool {
	return math.Min(currO, currC) > math.Max(prevO, prevC)
}

// bodyGapDown 检查实体（开盘-收盘范围）是否向下分离。
func bodyGapDown(prevO, prevC, currO, currC float64) bool {
	return math.Max(currO, currC) < math.Min(prevO, prevC)
}

// whiteBody 如果蜡烛有白色（看涨）实体则返回 true。
func whiteBody(o, c float64) bool {
	return c > o
}

// blackBody 如果蜡烛有黑色（看跌）实体则返回 true。
func blackBody(o, c float64) bool {
	return c < o
}
