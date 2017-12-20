package sd

import (
	"time"
	"github.com/qjsoftcn/gutils"
	"sort"
)

type Stock struct {
	Code       string
	Name       string
	Now DayLine
	Before DayLines
}

//股票日线行情
type DayLine struct {
	UTIME        time.Time//更新时间
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

//获取涨跌幅
func (this DayLine) GetQuoteChange() string {
	return gutils.FormatFloat(this.PCHG*100, "%2") + "%"
}

//获取振幅
func (this DayLine) GetAmplitude() float64 {
	return (this.HIGH - this.LOW) / this.LCLOSE
}

//获取平均成交价格
func (this DayLine) GetAvgPrice() float64 {
	if this.VOTURNOVER == 0 {
		return 0
	}
	return this.VATURNOVER / (float64(this.VOTURNOVER) * 100)
}

//(现价-均价)/均价
func (this DayLine) GetNADeviation() float64 {
	ap := this.GetAvgPrice()
	if ap == 0 {
		return 0
	}
	return (this.TCLOSE - ap) / ap
}

/*func (this DayLine) String() string {

	return fmt.Sprintln(gutils.FormatToSecond(this.Day),
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
}*/

type DayLines struct {
	dls    []DayLine
	dlsMap map[string]*DayLine
	by     By
}



func (this *DayLines) Get(fromIndex, endIndex int) []DayLine {
	return this.dls[fromIndex:endIndex]
}

func (this *DayLines) Len() int {
	return len(this.dls)
}

func (this *DayLines) Swap(i, j int) {
	this.dls[i], this.dls[j] = this.dls[j], this.dls[i]
}

func (this *DayLines) Less(i, j int) bool {
	return this.by(&this.dls[i], &this.dls[j])
}

func (this *DayLines) Sorts(bys ...By) {
	for _, by := range bys {
		this.by = by
		sort.Sort(this)
	}
}

func StockIsSh(code string) bool {
	return code[0] == '6'
}
