package sd

import (
	"time"
	"../tx"
	"fmt"
)

//股票日线行情
type DayLine struct {
	UpdateTime time.Time
	Code       string
	Name       string
	//最新价
	LatestPrice float64
	//涨跌幅
	QuoteChange float64
	//涨跌额
	AmountChange float64
	//买1价
	Buy1Price float64
	//卖1价
	Sell1Price float64
	//成交量
	Volume int64
	//成交额
	TurnOver float64
	//今开
	NowOpenPrice float64
	//昨收
	LastClosedPrice float64
	//最高
	LowestPrice float64
	//最低
	HighestPrice float64
}

//获取日线数据，通过腾讯的接口
func GetDayLines() ([]DayLine, error) {

	dlsfile,err:=tx.ReadDailyLines()

	fmt.Println(dlsfile,err)
	return nil, nil
}
