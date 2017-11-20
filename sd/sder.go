package sd

import (
	"time"
	"github.com/qjsoftcn/gutils"
	"fmt"
	"sort"
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

//获取涨跌幅
func (this DayLine) GetQuoteChange() string {
	return gutils.FormatFloat(this.QuoteChange*100, "%2") + "%"
}

//获取振幅
func (this DayLine) GetAmplitude() float64 {
	return (this.HighestPrice - this.LowestPrice) / this.LastClosedPrice
}

//获取平均成交价格
func (this DayLine) GetAvgPrice() float64 {
	if this.Volume == 0 {
		return 0
	}
	return this.TurnOver / (float64(this.Volume) * 100)
}

//(现价-均价)/均价
func (this DayLine) GetNADeviation() float64 {
	ap := this.GetAvgPrice()
	if ap == 0 {
		return 0
	}
	return (this.LatestPrice - ap) / ap
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
	by     By
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

type By func(t1, t2 *DayLine) bool

func (this *DayLines) GetDLS() []DayLine {
	return this.dls
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

func GetByAmp(asc bool) By {
	return func(d1, d2 *DayLine) bool {
		zf1 := d1.GetAmplitude()
		zf2 := d2.GetAmplitude()
		if asc {
			return zf1 < zf2
		} else {
			return zf1 > zf2
		}
	}
}

func GetByNad(asc bool) By {
	return func(d1, d2 *DayLine) bool {
		n1 := d1.GetNADeviation()
		n2 := d2.GetNADeviation()
		if asc {
			return n1 < n2
		} else {
			return n1 > n2
		}
	}
}

func GetByVolume(asc bool) By {
	return func(d1, d2 *DayLine) bool {
		n1 := d1.Volume
		n2 := d2.Volume
		if asc {
			return n1 < n2
		} else {
			return n1 > n2
		}
	}
}

func (this *DayLines) Sorts(bys ...By) {
	for _, by := range bys {
		this.by = by
		sort.Sort(this)
	}
}
