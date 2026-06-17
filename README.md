# TA-Lib 纯 Go 实现

## 项目背景

[TA-Lib](https://ta-lib.org/)（Technical Analysis Library）是金融市场技术分析领域的事实标准库，提供 150+ 个技术分析指标，广泛应用于量化交易、算法交易和金融数据分析。原始 TA-Lib 使用 C 语言编写，有 C/C++、Python、Java、.NET 等语言的封装，但缺少官方的 Go 语言绑定。

本项目是完全使用**纯 Go 语言**重新实现的 TA-Lib 核心功能，零 CGo 依赖，便于 Go 生态系统的集成。

## 参考

- 官方 TA-Lib: https://ta-lib.org/
- TA-Lib C 源码: https://github.com/ta-lib/ta-lib
- TA-Lib Python 封装: https://github.com/ta-lib/ta-lib-python
- TA-Lib 函数列表: https://ta-lib.org/functions/

## 安装

```bash
go get github.com/bambuo/talib
```

## 使用约定

所有函数遵循以下 Go 惯例：

- **纯函数**：无共享状态，所有输入通过参数传递
- **输入输出**：接受 `[]float64`，返回 `[]float64`
- **前导 NaN**：输出长度与输入一致，lookback 范围内的前导位置填充 `math.NaN()`，与 TA-Lib C 的输出长度约定一致
- **lookback 函数**：每个指标函数对应一个 `*Lookback()` 函数，返回有效输出前需要跳过的元素数
- **精度目标**：所有函数输出与 TA-Lib C 的偏差 `< 1e-10`

### 基本用法

```go
package main

import (
    "fmt"
    "github.com/bambuo/talib"
)

func main() {
    input := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    result, err := talib.SMA(input, 3)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // [NaN, NaN, 2, 3, 4, 5, 6, 7, 8, 9]
}
```

## 已实现功能

### 重叠研究 (Overlap Studies)

| 函数 | 说明 |
|------|------|
| `MIDPOINT` | 区间中点价格 |
| `MIDPRICE` | 中位数价格 |
| `SAR` | 抛物线转向系统（Parabolic SAR） |
| `SAREXT` | 扩展抛物线 SAR |
| `MAVP` | 可变周期移动平均 |

### 移动平均 (Moving Averages)

| 函数 | 说明 |
|------|------|
| `SMA` | 简单移动平均 |
| `EMA` | 指数移动平均 |
| `WMA` | 加权移动平均 |
| `DEMA` | 双重指数移动平均 |
| `TEMA` | 三重指数移动平均 |
| `TRIMA` | 三角移动平均 |
| `KAMA` | 自适应移动平均（Kaufman） |
| `MAMA` | 自适应移动平均（MESA） |
| `MA` | 通用移动平均调度器 |

### 动量指标 (Momentum Indicators)

| 函数 | 说明 |
|------|------|
| `BOP` | 平衡成交量 |
| `CMO` | 钱德动量摆动指标 |
| `DX` | 方向指标 |
| `MOM` | 动量 |
| `ROC` | 变动率 |
| `ROCP` | 变动率百分比 |
| `ROCR` | 变动率比率 |
| `ROCR100` | 变动率比率 × 100 |
| `RSI` | 相对强弱指标 |
| `STOCH` | 随机指标 |
| `STOCHF` | 快速随机指标 |
| `STOCHRSI` | 随机 RSI |
| `TRIX` | 三重指数平滑移动平均 |
| `ULTOSC` | 终极摆动指标 |
| `WILLIAMS_R` | 威廉指标 %R |
| `ADX` | 平均趋向指数 |
| `ADXR` | 平均趋向指数评级 |
| `PLUS_DM` / `MINUS_DM` | 正/负方向变动 |
| `PLUS_DI` / `MINUS_DI` | 正/负方向指标 |
| `AROON` / `AROONOSC` | Aroon / Aroon 振荡器 |

### 成交量指标 (Volume Indicators)

| 函数 | 说明 |
|------|------|
| `OBV` | 能量潮 |
| `AD` | 累积/分配线 |
| `ADOSC` | 累积/分配振荡器 |

### 波动率指标 (Volatility Indicators)

| 函数 | 说明 |
|------|------|
| `ATR` | 平均真实波幅 |
| `NATR` | 归一化平均真实波幅 |
| `TRANGE` | 真实波幅 |

### 价格变换 (Price Transform)

| 函数 | 说明 |
|------|------|
| `AVGPRICE` | 平均价格 `(O+H+L+C)/4` |
| `MEDPRICE` | 中位价格 `(H+L)/2` |
| `TYPPRICE` | 典型价格 `(H+L+C)/3` |
| `WCLPRICE` | 加权收盘价 `(H+L+2C)/4` |

### 数学运算 (Math Operators)

| 函数 | 说明 |
|------|------|
| `ADD` | 向量加法 |
| `DIV` | 向量除法 |
| `MULT` | 向量乘法 |
| `SUB` | 向量减法 |
| `SUM` | 求和 |
| `MAX` / `MAXINDEX` | 最大值 / 最大值索引 |
| `MIN` / `MININDEX` | 最小值 / 最小值索引 |
| `MINMAX` / `MINMAXINDEX` | 最小最大值 / 最小最大值索引 |

### 数学变换 (Math Transform)

| 函数 | 说明 |
|------|------|
| `ACOS` / `ASIN` / `ATAN` | 反三角函数 |
| `CEIL` / `FLOOR` | 向上/向下取整 |
| `COS` / `COSH` / `SIN` / `SINH` / `TAN` / `TANH` | 三角函数 |
| `EXP` / `LN` / `LOG10` | 指数/对数 |
| `SQRT` | 平方根 |

### 统计函数 (Statistics)

| 函数 | 说明 |
|------|------|
| `STDDEV` | 标准差 |
| `VAR` | 方差 |
| `CORREL` | 皮尔逊相关系数 |
| `BETA` | Beta 系数 |
| `LINEARREG` | 线性回归 |
| `LINEARREG_SLOPE` | 线性回归斜率 |
| `LINEARREG_INTERCEPT` | 线性回归截距 |
| `LINEARREG_ANGLE` | 线性回归角度 |
| `TSF` | 时间序列预测 |

### 振荡器 (Oscillators)

| 函数 | 说明 |
|------|------|
| `APO` | 绝对价格振荡器 |
| `PPO` | 百分比价格振荡器 |
| `MACD` | 指数平滑异同移动平均 |
| `MACDEXT` | 可控制类型的 MACD |
| `MACDFIX` | 固定参数的 MACD |
| `T3` | T3 移动平均 (Tillson) |
| `CCI` | 商品通道指数 |
| `MFI` | 资金流量指数 |
| `BBands` | 布林带 |

### 希尔伯特变换 (Hilbert Transform)

| 函数 | 说明 |
|------|------|
| `HT_TRENDLINE` | 希尔伯特变换 — 瞬时趋势线 |
| `HT_TRENDMODE` | 希尔伯特变换 — 趋势模式 |
| `HT_DCPERIOD` | 希尔伯特变换 — 主导周期 |
| `HT_DCPHASE` | 希尔伯特变换 — 主导周期相位 |
| `HT_PHASOR` | 希尔伯特变换 — 相量分量 |
| `HT_SINE` | 希尔伯特变换 — 正弦波 |
| `HT_DCCOMPONENT` | 希尔伯特变换 — 去趋势价格振荡器 |

### 形态识别 (Pattern Recognition)

`talib/cdl` 子包提供 61 种 K 线形态识别函数，返回 `-100` / `0` / `100` 信号强度：

| 函数 | 说明 |
|------|------|
| `CDL2CROWS` | 两只乌鸦 |
| `CDL3BLACKCROWS` | 三只乌鸦 |
| `CDL3INSIDE` | 三内部上涨/下跌 |
| `CDL3LINESTRIKE` | 三线打击 |
| `CDL3OUTSIDE` | 三外部上涨/下跌 |
| `CDL3STARSINSOUTH` | 南方三星 |
| `CDL3WHITESOLDIERS` | 三白兵 |
| `CDLABANDONEDBABY` | 弃婴形态 |
| `CDLADVANCEBLOCK` | 推进红三兵 |
| `CDLBELTHOLD` | 捉腰带线 |
| `CDLBREAKAWAY` | 脱离形态 |
| `CDLCLOSINGMARUBOZU` | 收盘无影线 |
| `CDLCONCEALBABYSWALLOW` | 藏婴吞没 |
| `CDLCOUNTERATTACK` | 反击线 |
| `CDLDARKCLOUDCOVER` | 乌云盖顶 |
| `CDLDOJI` | 十字星 |
| `CDLDOJISTAR` | 十字星形态 |
| `CDLDRAGONFLYDOJI` | 蜻蜓十字星 |
| `CDLENGULFING` | 吞噬形态 |
| `CDLEVENINGDOJISTAR` | 黄昏十字星 |
| `CDLEVENINGSTAR` | 黄昏星 |
| `CDLGAPSIDESIDEWHITE` | 向上跳空并列阳线 |
| `CDLGRAVESTONEDOJI` | 墓碑十字星 |
| `CDLHAMMER` | 锤子线 |
| `CDLHANGINGMAN` | 上吊线 |
| `CDLHARAMI` | 孕线形态 |
| `CDLHARAMICROSS` | 十字孕线 |
| `CDLHIGHWAVE` | 长脚十字 |
| `CDLHIKKAKEMODIFI` | 上升三法变形 |
| `CDLHIKKAKE` | 上升三法 |
| `CDLHOMINGPIGEON` | 家鸽形态 |
| `CDLIDENTICAL3CROWS` | 三只乌鸦（相同收盘价） |
| `CDLINNECK` | 颈内线 |
| `CDLINVERTEDHAMMER` | 倒锤子线 |
| `CDLKICKING` | 踢形态 |
| `CDLKICKINGBYLENGTH` | 踢形态（长度确认） |
| `CDLLADDERBOTTOM` | 梯底形态 |
| `CDLLONGLEGGEDDOJI` | 长腿十字星 |
| `CDLLONGLINE` | 长实体线 |
| `CDLMARUBOZU` | 光头光脚 |
| `CDLMATCHINGLOW` | 低价相匹配 |
| `CDLMATHOLD` | 三线法 |
| `CDLMORNINGDOJISTAR` | 早晨十字星 |
| `CDLMORNINGSTAR` | 晨星 |
| `CDLONNECK` | 颈上线 |
| `CDLPIERCING` | 刺透形态 |
| `CDLRICKSHAWMAN` | 黄包车夫 |
| `CDLRISEFALL3METHODS` | 上升/下降三法 |
| `CDLSEPARATINGLINES` | 分离线 |
| `CDLSHOOTINGSTAR` | 射击之星 |
| `CDLSHORTLINE` | 短实体线 |
| `CDLSPINNINGTOP` | 陀螺线 |
| `CDLSTALLEDPATTERN` | 停顿形态 |
| `CDLSTICKSANDWICH` | 三明治形态 |
| `CDLTAKURI` | 探水杆 |
| `CDLTASUKIGAP` | 跳空并列阴阳线 |
| `CDLTHRUSTING` | 插入线 |
| `CDLTRISTAR` | 三星形态 |
| `CDLUNIQUE3RIVER` | 独特的三河形态 |
| `CDLUPSIDEGAP2CROWS` | 向上跳空两只乌鸦 |
| `CDLXSIDEGAP3METHODS` | 侧跳空三法 |

## 精度与验证

- 核心指标（SMA, EMA, RSI, MACD, BBands, STOCH, ATR 等）通过黄金数据对比验证
- 与 TA-Lib Python 版本输出完全一致（偏差 < 1e-10）
- 所有函数使用表驱动测试覆盖正常路径和边界条件

## 构建与测试

```bash
# 构建
go build ./...

# 运行所有测试
go test -count=1 ./...

# 运行基准测试
go test -bench=. ./...
```

## 项目结构

```
lib/talib/
├── cdl/               # 形态识别子包（61 种 K 线形态）
│   ├── cdl.go         # 共享辅助函数
│   ├── complex.go     # 复杂多烛形态
│   ├── single_k.go    # 单烛形态
│   ├── star.go        # 星线形态
│   ├── three_k.go     # 三烛形态
│   └── two_k.go       # 双烛形态
├── ta.go              # 包文档
├── types.go           # MAType 枚举
├── compute.go         # 内部计算辅助函数
├── validate.go        # 输入验证
├── ma.go              # 移动平均线
├── overlap.go         # 重叠研究
├── momentum.go        # 动量指标
├── oscillators.go     # 振荡器
├── statistics.go      # 统计函数
├── math_operators.go  # 数学运算
├── math_transform.go  # 数学变换
├── price_transform.go # 价格变换
├── ht.go              # 希尔伯特变换
├── mama.go            # MESA 自适应移动平均
├── macd.go            # MACD
├── bb.go              # 布林带
├── obv.go             # OBV
├── ad.go              # 累积/分配
├── atr.go             # ATR
├── adx.go             # ADX 系列
├── aroon.go           # Aroon
├── cci.go             # CCI
├── kama.go            # KAMA
├── mfi.go             # MFI
├── rsi.go             # RSI
├── stochastic.go      # 随机指标
├── ultimate.go        # ULTOSC
├── williams.go        # Williams %R
└── *_test.go          # 各模块测试
```

## 许可

MIT
