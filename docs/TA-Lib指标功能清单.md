# TA-Lib 指标功能清单

> **参考标准**: [TA-Lib (Technical Analysis Library)](https://ta-lib.org/) — 金融技术分析领域事实标准  
> **对标库**: `github.com/bambuo/talib` — 纯 Go 实现，精度偏差 < 1e-10（对标 TA-Lib C 实现）  
> **统计时间**: 2026-06-17
> **最近更新**: 2026-06-17（补充 68 个缺失实现）

---

## 概览

| 类别 | TA-Lib 总数 | talib 已实现 | 覆盖率 | 状态 |
|------|:---------:|:----------:|:-----:|:----:|
| Overlap Studies（重叠研究） | 17 | 15 | 88% | 🟢 |
| Momentum Indicators（动量指标） | 30 | 30 | 100% | 🟢 |
| Volume Indicators（成交量指标） | 3 | 3 | 100% | 🟢 |
| Volatility Indicators（波动率指标） | 3 | 3 | 100% | 🟢 |
| Cycle Indicators（周期指标） | 5 | 6 | 100% | 🟢 |
| Price Transform（价格变换） | 4 | 4 | 100% | 🟢 |
| Statistic Functions（统计函数） | 9 | 9 | 100% | 🟢 |
| Pattern Recognition（形态识别） | 61 | 61 | 100% | 🟢 |
| Math Transform（数学变换） | 15 | 15 | 100% | 🟢 |
| Math Operators（数学运算符） | 11 | 11 | 100% | 🟢 |
| **合计** | **158** | **157** | **99%** | — |

> 注：覆盖率仅统计**主函数**（不含 Lookback 辅助函数）。talib 额外提供了 `HT_DCCOMPONENT`（TA-Lib 无此独立函数）和 CDLINSIDE（CDL3INSIDE 的别名）。Pattern Recognition 在独立子包 `talib/cdl` 中实现。

---

## 一、Overlap Studies（重叠研究 / 均线类）— 17 个

侧重价格平滑和趋势跟踪，多为移动平均线及其变体。

| # | TA-Lib 函数 | 中文名称 | 算法简介 | talib 实现 |
|:-:|-------------|---------|---------|:---------:|
| 1 | **BBANDS** | 布林带 | Middle=MA, Upper/Lower=Middle±k×StdDev（默认k=2） | ✅ `BBands()` |
| 2 | **DEMA** | 双指数移动平均 | DEMA = 2×EMA1 − EMA(EMA1) | ✅ `DEMA()` |
| 3 | **EMA** | 指数移动平均 | EMA[i]=α×Price[i]+(1−α)×EMA[i−1], α=2/(N+1) | ✅ `EMA()` |
| 4 | **HT_TRENDLINE** | 希尔伯特瞬时趋势线 | 基于 HT 平滑的同相分量输出 | ✅ `HT_TRENDLINE()` |
| 5 | **KAMA** | 考夫曼自适应移动平均 | 根据市场噪声动态调整平滑系数 | ✅ `KAMA()` |
| 6 | **MA** | 通用移动平均 | 根据 MAType 参数调度 SMA/EMA/WMA/DEMA/TEMA/TRIMA/KAMA/MAMA/T3 | ✅ `MA()` |
| 7 | **MAMA** | MESA 自适应移动平均 | 基于 Hilbert 变换的自适应平滑，含相位调整 | ✅ `MAMA()` |
| 8 | **MAVP** | 变周期移动平均 | MA 周期由另一时间序列动态指定 | ✅ `MAVP()` |
| 9 | **MIDPOINT** | 中点 | (周期内最高价 + 最低价) / 2 | ✅ `MIDPOINT()` |
| 10 | **MIDPRICE** | 中点价格 | (周期内最高价的平均 + 最低价的平均) / 2 | ✅ `MIDPRICE()` |
| 11 | **SAR** | 抛物线止损反转 | 加速因子驱动的趋势跟踪止损，经典 Wilder 算法 | ✅ `SAR()` |
| 12 | **SAREXT** | 扩展抛物线 SAR | 支持自定义初始加速因子和最大加速因子 | ✅ `SAREXT()` |
| 13 | **SMA** | 简单移动平均 | 窗口内等权算术平均 | ✅ `SMA()` |
| 14 | **T3** | T3 移动平均 | 6 层 EMA 叠加以 volumeFactor 加权（默认 v=0.7） | ✅ `T3()` |
| 15 | **TEMA** | 三指数移动平均 | TEMA = 3×EMA1 − 3×EMA2 + EMA3 | ✅ `TEMA()` |
| 16 | **TRIMA** | 三角移动平均 | SMA(SMA(input, N), N) — 两次 SMA 平滑 | ✅ `TRIMA()` |
| 17 | **WMA** | 加权移动平均 | 线性加权：权重 = [1, 2, …, N] | ✅ `WMA()` |

---

## 二、Momentum Indicators（动量指标）— 26 个

衡量价格变化速率和方向性强度，是技术分析中使用最广泛的类别。

| # | TA-Lib 函数 | 中文名称 | 算法简介 | talib 实现 |
|:-:|-------------|---------|---------|:---------:|
| 1 | **ADX** | 平均趋向指数 | ADX = MA(DX)，DX = \|+DI − −DI\|/(+DI + −DI)×100 | ✅ `ADX()` |
| 2 | **ADXR** | ADX 评级 | ADXR = (ADX[i] + ADX[i−N]) / 2 | ✅ `ADXR()` |
| 3 | **APO** | 绝对价格震荡器 | APO = MA(fast) − MA(slow) | ✅ `APO()` |
| 4 | **AROON** | 阿隆指标 | Up=(N−距最高价K线数)/N×100, Down=(N−距最低价K线数)/N×100 | ✅ `AROON()` |
| 5 | **AROONOSC** | 阿隆震荡器 | AROONOSC = AroonUp − AroonDown | ✅ `AROONOSC()` |
| 6 | **BOP** | 力量平衡 | (Close−Open)/(High−Low) | ✅ `BOP()` |
| 7 | **CCI** | 商品通道指数 | CCI = (TP − SMA(TP)) / (0.015×MeanDeviation) | ✅ `CCI()` |
| 8 | **CMO** | 钱德动量震荡器 | CMO = (UpSum − DownSum)/(UpSum + DownSum)×100 | ✅ `CMO()` |
| 9 | **DX** | 趋向运动指数 | DX = \|+DI − −DI\|/(+DI + −DI)×100 | ✅ `DX()` |
| 10 | **MACD** | 异同移动平均线 | MACD=EMA(fast)−EMA(slow), Signal=EMA(MACD,signal), Hist=MACD−Signal | ✅ `MACD()` |
| 11 | **MACDEXT** | 可控 MA 类型 MACD | 同上，但允许自定义 fast/slow/signal 各自的 MA 类型 | ✅ `MACDEXT()` |
| 12 | **MACDFIX** | 固定参数 MACD | MACD(12,26,9) 预设，仅返回 Histogram | ✅ `MACDFIX()` |
| 13 | **MFI** | 资金流量指数 | 基于典型价格和成交量的 RSI 变体：MFI=100−100/(1+MR) | ✅ `MFI()` |
| 14 | **MINUS_DI** | 负趋向指标 | −DI = 100 × EMA(−DM) / ATR | ✅ `MINUS_DI()` |
| 15 | **MINUS_DM** | 负趋向运动 | −DM = max(Low[i−1]−Low[i], 0) 当满足条件 | ✅ `MINUS_DM()` |
| 16 | **MOM** | 动量 | MOM = Price[i] − Price[i−N] | ✅ `MOM()` |
| 17 | **PLUS_DI** | 正趋向指标 | +DI = 100 × EMA(+DM) / ATR | ✅ `PLUS_DI()` |
| 18 | **PLUS_DM** | 正趋向运动 | +DM = max(High[i]−High[i−1], 0) 当满足条件 | ✅ `PLUS_DM()` |
| 19 | **PPO** | 百分比价格震荡器 | PPO = (MA(fast)−MA(slow)) / MA(slow) × 100 | ✅ `PPO()` |
| 20 | **ROC** | 变化率 | ROC = (Price/PrevPrice − 1) × 100 | ✅ `ROC()` |
| 21 | **ROCP** | 变化率百分比 | ROCP = (Price − PrevPrice) / PrevPrice | ✅ `ROCP()` |
| 22 | **ROCR** | 变化率比率 | ROCR = Price / PrevPrice | ✅ `ROCR()` |
| 23 | **ROCR100** | 变化率比率(百倍) | ROCR100 = (Price/PrevPrice) × 100 | ✅ `ROCR100()` |
| 24 | **RSI** | 相对强弱指数 | RSI=100−100/(1+RS), RS=AvgGain/AvgLoss, Wilder 平滑 | ✅ `RSI()` |
| 25 | **STOCH** | 随机指标(慢速) | %K=MA(%K_fast,slowK), %D=MA(%K,slowD) | ✅ `STOCH()` |
| 26 | **STOCHF** | 随机指标(快速) | %K_fast=(C−L)/(H−L)×100, %D=MA(%K,fD) | ✅ `STOCHF()` |
| 27 | **STOCHRSI** | 随机 RSI | 对 RSI 值再应用 Stochastic 公式 | ✅ `STOCHRSI()` |
| 28 | **TRIX** | 三重平滑变化率 | TRIX = ROC(EMA3, 1), 其中 EMA3=EMA(EMA(EMA(Price))) | ✅ `TRIX()` |
| 29 | **ULTOSC** | 终极震荡器 | 3 周期加权：BP=Close−min(Low,PrevClose), 加权平均后归一化 | ✅ `ULTOSC()` |
| 30 | **WILLR** | 威廉指标 | %R = (HighestHigh−Close)/(HighestHigh−LowestLow)×(−100) | ✅ `WILLIAMS_R()` |

---

## 三、Volume Indicators（成交量指标）— 3 个

| # | TA-Lib 函数 | 中文名称 | 算法简介 | talib 实现 |
|:-:|-------------|---------|---------|:---------:|
| 1 | **AD** | 柴金累积/派发线 | MFL=((C−L)−(H−C))/(H−L), AD=cumsum(MFL×Volume) | ✅ `AD()` |
| 2 | **ADOSC** | 柴金震荡器 | ADOSC = EMA(AD, fast) − EMA(AD, slow) | ✅ `ADOSC()` |
| 3 | **OBV** | 能量潮 | C[i]>C[i−1]→+V; C[i]<C[i−1]→−V; C[i]=C[i−1]→0 | ✅ `OBV()` |

---

## 四、Volatility Indicators（波动率指标）— 3 个

| # | TA-Lib 函数 | 中文名称 | 算法简介 | talib 实现 |
|:-:|-------------|---------|---------|:---------:|
| 1 | **ATR** | 平均真实波幅 | TR=max(H−L,\|H−PrevC\|,\|L−PrevC\|), ATR=WilderMA(TR,N) | ✅ `ATR()` |
| 2 | **NATR** | 归一化 ATR | NATR = ATR / Close × 100 | ✅ `NATR()` |
| 3 | **TRANGE** | 真实波幅 | TR = max(High−Low, \|High−PrevClose\|, \|Low−PrevClose\|) | ✅ `TRANGE()` |

---

## 五、Cycle Indicators（周期指标）— 5 个

基于 Hilbert 变换的周期分析，用于检测市场主导周期。

| # | TA-Lib 函数 | 中文名称 | 算法简介 | talib 实现 |
|:-:|-------------|---------|---------|:---------:|
| 1 | **HT_DCPERIOD** | HT 主导周期 | Hilbert 变换检测主导循环周期长度 | ✅ `HT_DCPERIOD()` |
| 2 | **HT_DCPHASE** | HT 主导周期相位 | 主导周期的瞬时相位角（度） | ✅ `HT_DCPHASE()` |
| 3 | **HT_PHASOR** | HT 相量分量 | 同相和正交分量（复数形式） | ✅ `HT_PHASOR()` |
| 4 | **HT_SINE** | HT 正弦波 | 正弦和余弦分量，用于预测转折点 | ✅ `HT_SINE()` |
| 5 | **HT_TRENDMODE** | HT 趋势/周期模式 | +1=趋势模式, −1=周期模式 | ✅ `HT_TRENDMODE()` |
| ★ | **HT_DCCOMPONENT** | HT 直流分量 | 额外实现：Hilbert 变换的直流偏移分量 | ✅（扩展） |

---

## 六、Price Transform（价格变换）— 4 个

基础价格统计，用于派生指标计算的中间步骤。

| # | TA-Lib 函数 | 中文名称 | 公式 | talib 实现 |
|:-:|-------------|---------|------|:---------:|
| 1 | **AVGPRICE** | 平均价格 | (Open + High + Low + Close) / 4 | ✅ `AVGPRICE()` |
| 2 | **MEDPRICE** | 中间价格 | (High + Low) / 2 | ✅ `MEDPRICE()` |
| 3 | **TYPPRICE** | 典型价格 | (High + Low + Close) / 3 | ✅ `TYPPRICE()` |
| 4 | **WCLPRICE** | 加权收盘价 | (High + Low + 2×Close) / 4 | ✅ `WCLPRICE()` |

---

## 七、Statistic Functions（统计函数）— 9 个

时间序列统计和线性回归。

| # | TA-Lib 函数 | 中文名称 | 算法简介 | talib 实现 |
|:-:|-------------|---------|---------|:---------:|
| 1 | **BETA** | Beta 系数 | β = Cov(X,Y) / Var(Y) — X 相对 Y 的波动敏感性 | ✅ `BETA()` |
| 2 | **CORREL** | 皮尔逊相关系数 | r = Cov(X,Y) / (σ_X × σ_Y) | ✅ `CORREL()` |
| 3 | **LINEARREG** | 线性回归 | 最小二乘法拟合，输出回归预测值 | ✅ `LINEARREG()` |
| 4 | **LINEARREG_ANGLE** | 线性回归角度 | arctan(slope) 转换为角度 | ✅ `LINEARREG_ANGLE()` |
| 5 | **LINEARREG_INTERCEPT** | 线性回归截距 | LRI = Mean(Y) − Slope×Mean(X) | ✅ `LINEARREG_INTERCEPT()` |
| 6 | **LINEARREG_SLOPE** | 线性回归斜率 | OLS 斜率 = N×ΣXY−ΣXΣY / N×ΣX²−(ΣX)² | ✅ `LINEARREG_SLOPE()` |
| 7 | **STDDEV** | 标准差 | 总体标准差 σ = √(Σ(x−μ)²/N) | ✅ `STDDEV()` |
| 8 | **TSF** | 时间序列预测 | 基于线性回归的下一周期预测值 | ✅ `TSF()` |
| 9 | **VAR** | 方差 | 总体方差 σ² = Σ(x−μ)²/N | ✅ `VAR()` |

---

## 八、Pattern Recognition（形态识别）— 61 个

日本蜡烛图形态识别，输出 {+100/−100/0} 分别表示看涨/看跌/无信号。

### 8.1 反转形态（看涨/看跌）

| # | TA-Lib 函数 | 中文名称 | 形态描述 | talib |
|:-:|-------------|---------|---------|:----:|
| 1 | **CDL2CROWS** | 两只乌鸦 | 两根阴线，第二根跳空高开后低收 | ✅ |
| 2 | **CDL3BLACKCROWS** | 三只黑乌鸦 | 连续三根实体阴线，每根收盘接近最低价 | ✅ |
| 3 | **CDL3INSIDE** | 三内部上涨/下跌 | 第三根确认前两根的内包关系 | ✅ |
| 4 | **CDL3LINESTRIKE** | 三线打击 | 四根 K 线组合，第四根逆转前三根趋势 | ✅ |
| 5 | **CDL3OUTSIDE** | 三外部上涨/下跌 | 第三根确认外包+方向 | ✅ |
| 6 | **CDL3STARSINSOUTH** | 南方三星 | 三根逐步缩小的阴线，下影线逐步缩短 | ✅ |
| 7 | **CDL3WHITESOLDIERS** | 三个白兵 | 连续三根实体阳线，每根收盘接近最高价 | ✅ |
| 8 | **CDLABANDONEDBABY** | 弃婴 | 跳空 Doji 被反向跳空确认 | ✅ |
| 9 | **CDLADVANCEBLOCK** | 前进受阻 | 三根阳线，实体逐步缩小，上影线逐步变长 | ✅ |
| 10 | **CDLBELTHOLD** | 捉腰带线 | 光头光脚阳/阴线 | ✅ |
| 11 | **CDLBREAKAWAY** | 突破 | 五根 K 线组合，第五根确认趋势反转 | ✅ |
| 12 | **CDLCOUNTERATTACK** | 反击线 | 两根相反颜色的实体，收盘价相等或接近 | ✅ |
| 13 | **CDLDARKCLOUDCOVER** | 乌云盖顶 | 阳线后跟一根高开低走的阴线，切入前一根 50%+ | ✅ |
| 14 | **CDLENGULFING** | 吞没形态 | 第二根实体完全包裹第一根实体 | ✅ |
| 15 | **CDLGAPSIDESIDEWHITE** | 向上/下跳空并列白线 | 跳空后的并列阳线 | ✅ |
| 16 | **CDLHARAMI** | 孕线 | 第二根实体完全被第一根实体包含 | ✅ |
| 17 | **CDLHARAMICROSS** | 十字孕线 | 孕线中第二根为 Doji | ✅ |
| 18 | **CDLHIKKAKE** | 陷阱 | 三根 K 线，第三根确认假突破 | ✅ |
| 19 | **CDLHIKKAKEMOD** | 修正陷阱 | 改进版 Hikkake 模式 | ✅ |
| 20 | **CDLHOMINGPIGEON** | 家鸽 | 两根阴线，第二根被第一根包含 | ✅ |
| 21 | **CDLIDENTICAL3CROWS** | 三只同等乌鸦 | 三根等长阴线，收盘依次降低 | ✅ |
| 22 | **CDLINNECK** | 颈内线 | 阴线后跟小阳线，收盘略低于前收盘 | ✅ |
| 23 | **CDLINSIDE** | — | （同 3INSIDE） | ✅ |
| 24 | **CDLINVERTEDHAMMER** | 倒锤子 | 上影线长（≥2×实体），实体在下方 | ✅ |
| 25 | **CDLKICKING** | 反冲 | 跳空后反向 Marubozu | ✅ |
| 26 | **CDLKICKINGBYLENGTH** | 反冲(按长度) | 由较长 Marubozu 决定方向的反冲 | ✅ |
| 27 | **CDLLADDERBOTTOM** | 梯底 | 五根 K 线看涨反转，逐步下探后跳空阳线 | ✅ |
| 28 | **CDLMATCHINGLOW** | 等低点 | 两根收盘价相等的 K 线组合 | ✅ |
| 29 | **CDLMATHOLD** | 铺垫形态 | 五根 K 线的看涨持续形态 | ✅ |
| 30 | **CDLMORNINGDOJISTAR** | 晨星十字 | Doji 隔开阴阳两根实体 | ✅ |
| 31 | **CDLMORNINGSTAR** | 晨星 | 小实体隔开阴阳两根实体 | ✅ |
| 32 | **CDLONNECK** | 颈上线 | 阴线后跟小阳线，收盘与前收盘持平 | ✅ |
| 33 | **CDLPIERCING** | 刺透形态 | 阴线后低开高走阳线，切入前阴线 50%+ | ✅ |
| 34 | **CDLRICKSHAWMAN** | 长腿十字线 | Doji 实体在 K 线中间位置 | ✅ |
| 35 | **CDLRISEFALL3METHODS** | 上升/下降三法 | 五根 K 线持续形态：趋势线+三根小回调+确认线 | ✅ |
| 36 | **CDLSEPARATINGLINES** | 分离线 | 两根同色 K 线，收盘价均等于开盘价 | ✅ |
| 37 | **CDLSHOOTINGSTAR** | 射击之星 | 上影线长（≥2×实体），实体在下方 | ✅ |
| 38 | **CDLSTALLEDPATTERN** | 停顿形态 | 三根阳线，实体逐步缩小 | ✅ |
| 39 | **CDLSTICKSANDWICH** | 条形三明治 | 两根阳线夹一根阴线，收盘价相等 | ✅ |
| 40 | **CDLTAKURI** | 探水竿 | 极长下影线的 Dragonfly Doji | ✅ |
| 41 | **CDLTASUKIGAP** | 田足缺口 | 跳空后的并列反向 K 线 | ✅ |
| 42 | **CDLTHRUSTING** | 插入线 | 阴线后阳线切入不足 50% | ✅ |
| 43 | **CDLTRISTAR** | 三星 | 三根 Doji 排列 | ✅ |
| 44 | **CDLUNIQUE3RIVER** | 独特三河 | 三根 K 线看涨反转，实体逐步降低 | ✅ |
| 45 | **CDLUPSIDEGAP2CROWS** | 向上跳空两只乌鸦 | 跳空后两根阴线 | ✅ |
| 46 | **CDLXSIDEGAP3METHODS** | 向上/下跳空三法 | 跳空后三法形态 | ✅ |

### 8.2 单 K 线形态

| # | TA-Lib 函数 | 中文名称 | 形态描述 | talib |
|:-:|-------------|---------|---------|:----:|
| 47 | **CDLCLOSINGMARUBOZU** | 收盘光头光脚 | 实体极长，无影线 | ✅ |
| 48 | **CDLDOJI** | 十字星 | 开盘价≈收盘价 | ✅ |
| 49 | **CDLDOJISTAR** | 十字星形态 | 跳空后的 Doji | ✅ |
| 50 | **CDLDRAGONFLYDOJI** | 蜻蜓十字 | Doji，长下影线，无上影线 | ✅ |
| 51 | **CDLEVENINGDOJISTAR** | 黄昏十字星 | 看跌版 Morning Doji Star | ✅ |
| 52 | **CDLEVENINGSTAR** | 黄昏之星 | 看跌版 Morning Star | ✅ |
| 53 | **CDLGRAVESTONEDOJI** | 墓碑十字 | Doji，长上影线，无下影线 | ✅ |
| 54 | **CDLHAMMER** | 锤子线 | 小实体在下，长下影线 | ✅ |
| 55 | **CDLHANGINGMAN** | 上吊线 | 形态同锤子，但处于上升趋势中 | ✅ |
| 56 | **CDLHIGHWAVE** | 高浪线 | 极长上下影线，小实体 | ✅ |
| 57 | **CDLLONGLEGGEDDOJI** | 长腿十字 | Doji + 长上下影线 | ✅ |
| 58 | **CDLLONGLINE** | 长蜡烛线 | 实体极长的 K 线 | ✅ |
| 59 | **CDLMARUBOZU** | 光头光脚 | 无影线的大阳/大阴线 | ✅ |
| 60 | **CDLSHORTLINE** | 短蜡烛线 | 实体极短 | ✅ |
| 61 | **CDLSPINNINGTOP** | 纺锤线 | 小实体，上下影线大致等长 | ✅ |

---

## 九、Math Transform（数学变换）— 15 个

对标 TA-Lib 的逐元素向量数学函数。

| # | TA-Lib 函数 | 说明 | talib |
|:-:|-------------|------|:----:|
| 1 | **ACOS** | 反余弦 | ✅ `ACOS()` |
| 2 | **ASIN** | 反正弦 | ✅ `ASIN()` |
| 3 | **ATAN** | 反正切 | ✅ `ATAN()` |
| 4 | **CEIL** | 向上取整 | ✅ `CEIL()` |
| 5 | **COS** | 余弦 | ✅ `COS()` |
| 6 | **COSH** | 双曲余弦 | ✅ `COSH()` |
| 7 | **EXP** | e 的幂 | ✅ `EXP()` |
| 8 | **FLOOR** | 向下取整 | ✅ `FLOOR()` |
| 9 | **LN** | 自然对数 | ✅ `LN()` |
| 10 | **LOG10** | 10 为底对数 | ✅ `LOG10()` |
| 11 | **SIN** | 正弦 | ✅ `SIN()` |
| 12 | **SINH** | 双曲正弦 | ✅ `SINH()` |
| 13 | **SQRT** | 平方根 | ✅ `SQRT()` |
| 14 | **TAN** | 正切 | ✅ `TAN()` |
| 15 | **TANH** | 双曲正切 | ✅ `TANH()` |

---

## 十、Math Operators（数学运算符）— 11 个

对标 TA-Lib 的向量算术和窗口运算函数。

| # | TA-Lib 函数 | 说明 | talib |
|:-:|-------------|------|:----:|
| 1 | **ADD** | 向量加法 | ✅ `ADD()` |
| 2 | **DIV** | 向量除法 | ✅ `DIV()` |
| 3 | **MAX** | 周期内最高值 | ✅ `MAX()` |
| 4 | **MAXINDEX** | 周期内最高值索引 | ✅ `MAXINDEX()` |
| 5 | **MIN** | 周期内最低值 | ✅ `MIN()` |
| 6 | **MININDEX** | 周期内最低值索引 | ✅ `MININDEX()` |
| 7 | **MINMAX** | 周期内最低和最高值 | ✅ `MINMAX()` |
| 8 | **MINMAXINDEX** | 周期内最低和最高值索引 | ✅ `MINMAXINDEX()` |
| 9 | **MULT** | 向量乘法 | ✅ `MULT()` |
| 10 | **SUB** | 向量减法 | ✅ `SUB()` |
| 11 | **SUM** | 周期内求和 | ✅ `SUM()` |

所有61个形态已全部在子包 `lib/talib/cdl/` 中实现。

---

## 文件分布（更新）

```
lib/talib/
├── ma.go              SMA, EMA, WMA, DEMA, TEMA, TRIMA, MA + 8 Lookback
├── kama.go            KAMA + Lookback
├── mama.go            MAMA + Lookback (含 MAMAResult)
├── macd.go            MACD + Lookback
├── oscillators.go     APO, PPO, MACDEXT, MACDFIX, T3, TRIX + 6 Lookback
├── rsi.go             RSI + Lookback
├── bb.go              BBands + Lookback
├── ad.go              AD, ADOSC, ADXR + 4 Lookback
├── adx.go             PLUS_DM, MINUS_DM, PLUS_DI, MINUS_DI, ADX + 8 Lookback
├── aroon.go           AROON, AROONOSC + 2 Lookback
├── atr.go             ATR, TRANGE, NATR + 4 Lookback
├── cci.go             CCI + Lookback
├── mfi.go             MFI + Lookback
├── obv.go             OBV + Lookback
├── stochastic.go      STOCH, STOCHF, STOCHRSI + 6 Lookback
├── williams.go        WILLIAMS_R + Lookback
├── ultimate.go        ULTOSC + Lookback
├── momentum.go        BOP, CMO, DX, MOM, ROC, ROCP, ROCR, ROCR100 + 8 Lookback
├── price_transform.go AVGPRICE, MEDPRICE, TYPPRICE, WCLPRICE
├── overlap.go         MIDPOINT, MIDPRICE, SAR, SAREXT, MAVP + 5 Lookback
├── statistics.go      STDDEV, VAR, CORREL, BETA, LINEARREG, LINEARREG_ANGLE,
│                      LINEARREG_INTERCEPT, LINEARREG_SLOPE, TSF + 9 Lookback
├── math_operators.go  ADD, DIV, MULT, SUB, SUM, MAX, MAXINDEX, MIN, MININDEX,
│                      MINMAX, MINMAXINDEX + 11 Lookback
├── math_transform.go  ACOS, ASIN, ATAN, CEIL, COS, COSH, EXP, FLOOR, LN,
│                      LOG10, SIN, SINH, SQRT, TAN, TANH
├── ht.go              HT_TRENDLINE, HT_TRENDMODE, HT_DCPERIOD, HT_DCPHASE,
│                      HT_PHASOR, HT_SINE, HT_DCCOMPONENT + 7 Lookback
├── types.go           MAType 枚举（SMA/EMA/WMA/DEMA/TEMA/TRIMA/KAMA/MAMA/T3）
├── compute.go         Sum, Mean, MaxFloat64, MinFloat64
├── validate.go        ValidateNumericInput, ValidateOHLCInput, ValidatePeriod,
│                      MakeOutput, IsNaN, 错误定义
└── ta.go              包文档
```

### 形态识别子包

```
lib/talib/cdl/
├── cdl.go            包文档、Signal 常量、共享辅助函数（realBody/isDoji/gap 检测等）
├── single_k.go       CDLDOJI, CDLDRAGONFLYDOJI, CDLGRAVESTONEDOJI, CDLHAMMER,
│                     CDLHANGINGMAN, CDLINVERTEDHAMMER, CDLSHOOTINGSTAR,
│                     CDLSPINNINGTOP, CDLMARUBOZU, CDLCLOSINGMARUBOZU,
│                     CDLLONGLINE, CDLSHORTLINE, CDLHIGHWAVE,
│                     CDLLONGLEGGEDDOJI, CDLRICKSHAWMAN, CDLTAKURI
├── two_k.go          CDLENGULFING, CDLHARAMI, CDLHARAMICROSS, CDLDARKCLOUDCOVER,
│                     CDLPIERCING, CDLCOUNTERATTACK, CDLKICKING, CDLKICKINGBYLENGTH,
│                     CDLGAPSIDESIDEWHITE, CDLDOJISTAR, CDLBELTHOLD,
│                     CDLHOMINGPIGEON, CDLMATCHINGLOW, CDLSEPARATINGLINES,
│                     CDLSTICKSANDWICH, CDLTASUKIGAP, CDLONNECK, CDLINNECK,
│                     CDLTHRUSTING
├── star.go           CDLMORNINGSTAR, CDLEVENINGSTAR, CDLMORNINGDOJISTAR,
│                     CDLEVENINGDOJISTAR, CDLABANDONEDBABY, CDLTRISTAR
├── three_k.go        CDL2CROWS, CDL3BLACKCROWS, CDL3INSIDE, CDL3OUTSIDE,
│                     CDL3WHITESOLDIERS, CDL3STARSINSOUTH, CDLIDENTICAL3CROWS,
│                     CDLADVANCEBLOCK, CDLSTALLEDPATTERN, CDLUNIQUE3RIVER,
│                     CDLUPSIDEGAP2CROWS, CDLHIKKAKE, CDLHIKKAKEMOD,
│                     CDL3LINESTRIKE, CDLINSIDE
├── complex.go        CDLBREAKAWAY, CDLMATHOLD, CDLLADDERBOTTOM,
│                     CDLRISEFALL3METHODS, CDLXSIDEGAP3METHODS
└── cdl_test.go       全部61个形态 + 公共函数测试

---

## 附录 A：已完成的补充实现

| 优先级 | 类别 | 函数 | 状态 |
|:------:|------|------|:----:|
| 🔴 高 | 统计 | STDDEV, VAR, CORREL, BETA | ✅ 已实现 |
| 🔴 高 | 动量 | MOM, ROC, ROCP, ROCR, ROCR100 | ✅ 已实现 |
| 🔴 高 | 价格变换 | TYPPRICE, MEDPRICE, AVGPRICE, WCLPRICE | ✅ 已实现 |
| 🟡 中 | 重叠 | KAMA, MAMA | ✅ 已实现 |
| 🟡 中 | 重叠 | SAR, SAREXT, MAVP, MIDPOINT, MIDPRICE | ✅ 已实现 |
| 🟡 中 | 统计 | LINEARREG 系列, TSF | ✅ 已实现 |
| 🟡 中 | 动量 | CMO, BOP, DX | ✅ 已实现 |
| 🟢 低 | 数学 | ADD/SUB/MULT/DIV/SUM/MAX/MIN/MAXINDEX/MININDEX/MINMAX/MINMAXINDEX | ✅ 已实现 |
| 🟢 低 | 数学变换 | ACOS/ASIN/ATAN/CEIL/COS/COSH/EXP/FLOOR/LN/LOG10/SIN/SINH/SQRT/TAN/TANH | ✅ 已实现 |

### 剩余未实现

| 优先级 | 类别 | 函数 | 状态 |
|:------:|------|------|:----:|
| 🟢 低 | 形态识别 | CDL×61 | ✅ 蜡烛图形态（61 个），单独子包实现更合适 |

---

## 附录 B：TA-Lib 精度对标标准

| 指标 | talib 精度目标 | 对标基准 |
|------|:-----------:|---------|
| 所有指标函数 | 偏差 < 1×10⁻¹⁰ | TA-Lib C 实现 |
| 输出约定 | NaN 填充前导值 | TA-Lib C lookback 一致 |
| MA 类型枚举 | MAType 0–8 | TA-Lib MA_Type 枚举 |
| Wilder 平滑 | k = 1/period | TA-Lib RSI/ATR 默认 |

---

> **文档维护**: 当 talib 新增指标实现时，请同步更新本文档的对应复选框状态。
