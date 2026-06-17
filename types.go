package talib

// MAType 定义移动平均计算方法。
// 值与 TA-Lib 的 MA_Type 枚举一致。
type MAType int

const (
	MASMA   MAType = 0 // 简单移动平均
	MAEMA   MAType = 1 // 指数移动平均
	MAWMA   MAType = 2 // 加权移动平均
	MADEMA  MAType = 3 // 双重指数移动平均
	MATEMA  MAType = 4 // 三重指数移动平均
	MATRIMA MAType = 5 // 三角移动平均
	MAKAMA  MAType = 6 // 考夫曼自适应移动平均
	MAMAMA  MAType = 7 // MESA 自适应移动平均
	MAT3    MAType = 8 // T3 移动平均
)

// String 返回 MAType 的名称。
func (m MAType) String() string {
	switch m {
	case MASMA:
		return "SMA"
	case MAEMA:
		return "EMA"
	case MAWMA:
		return "WMA"
	case MADEMA:
		return "DEMA"
	case MATEMA:
		return "TEMA"
	case MATRIMA:
		return "TRIMA"
	case MAKAMA:
		return "KAMA"
	case MAMAMA:
		return "MAMA"
	case MAT3:
		return "T3"
	default:
		return "Unknown"
	}
}
