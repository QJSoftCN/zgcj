package sd

import (
	"time"
)

type Stock struct {
	Code       string
	Name       string
	Day        time.Time
	TCLOSE     float64 //收盘价
	HIGH       float64 //最高价
	LOW        float64 //最低价
	TOPEN      float64 //开盘价
	LCLOSE     float64 //前收盘
	CHG        float64 //涨跌额
	PCHG       float64 //涨跌幅
	TURNOVER   float64 //换手率
	VOTURNOVER float64 //成交量
	VATURNOVER float64 //成交金额
	TCAP       float64 //总市值
	MCAP       float64 //流通市值
}

func StockIsSh(code string) bool {
	return code[0] == '6'
}
