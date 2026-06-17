// Package talib 提供了常见技术分析指标的纯 Go 实现，
// 功能上等同于 TA-Lib 的 C 语言实现。
//
// 所有函数都是纯函数（无共享状态），接受 []float64 输入，
// 返回 []float64 输出，前导值用 NaN 填充以匹配
// TA-Lib 的输出长度约定。
//
// 精度目标：所有函数的输出与 TA-Lib C 的偏差 < 1e-10。
package talib
