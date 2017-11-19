package sd

import (
	"time"
	"github.com/qjsoftcn/gutils"
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

func (this DayLine) GetQuoteChange() string {
	return gutils.FormatFloat(this.QuoteChange*100, "%2") + "%"
}

func(this DayLine)GetMaxPriceChange()float64{
	return this.HighestPrice-this.LowestPrice
}

func (this DayLine) String() string {

	return fmt.Sprintln(gutils.FormatToSecond(this.UpdateTime),
		this.Code, this.Name, gutils.FormatFloat(this.LatestPrice, "2"),
		this.GetQuoteChange(),
		gutils.FormatFloat(this.AmountChange, "2"),
		gutils.FormatFloat(this.Buy1Price, "2"),
		gutils.FormatFloat(this.Sell1Price, "2"),
		this.Volume, gutils.FormatFloat(this.TurnOver, "2"),
		gutils.FormatFloat(this.NowOpenPrice, "2"),
		gutils.FormatFloat(this.LastClosedPrice, "2"),
		gutils.FormatFloat(this.HighestPrice, "2"),
		gutils.FormatFloat(this.LowestPrice, "2"))
}

type DayLines struct {
	dls    []DayLine
	dlsMap map[string]*DayLine
}

func NewDayLines(ds []DayLine) *DayLines {
	dls := new(DayLines)
	dls.dls = ds

	dlsMap := make(map[string]*DayLine, len(ds))
	for _, d := range ds {
		dlsMap[d.Code] = &d
		dlsMap[d.Name] = &d
	}

	dls.dlsMap = dlsMap
	return dls
}

func (this *DayLines) Toss(){
	//价格变化
	for _,d:=range this.dls{
		pc:=d.GetMaxPriceChange()
		fmt.Println(pc)
	}

}
