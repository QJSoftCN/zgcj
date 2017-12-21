package sd

import (
	"time"
	"github.com/qjsoftcn/gutils"
	"sort"
	"strconv"
)

const (
	C_Step int = 60
)

type Stocks struct {
	stocks   []*Stock
	stockMap map[string]*Stock
}

func (this *Stocks) getCodes(s, e int) []string {
	codes := make([]string, e-s+1)
	ss := this.stocks[s:e]
	for index, code := range ss {
		codes[index] = code.Code
	}
	return codes

}

func NewNowHQ(hqs []string) *DayLine {
	now := new(DayLine)
	now.CHG, _ = strconv.ParseFloat(hqs[30], 64)
	now.HIGH, _ = strconv.ParseFloat(hqs[40], 64)
	now.LCLOSE, _ = strconv.ParseFloat(hqs[3], 64)
	now.LOW, _ = strconv.ParseFloat(hqs[41], 64)
	now.MCAP, _ = strconv.ParseFloat(hqs[43], 64)
	now.PCHG, _ = strconv.ParseFloat(hqs[31], 64)
	now.TCAP, _ = strconv.ParseFloat(hqs[44], 64)

	now.TOPEN, _ = strconv.ParseFloat(hqs[4], 64)
	now.TURNOVER, _ = strconv.ParseFloat(hqs[37], 64)
	//now.UTIME, _ = strconv.ParseFloat(hqs[31], 64)
	now.TCLOSE, _ = strconv.ParseFloat(hqs[2], 64)
	now.VATURNOVER, _ = strconv.ParseFloat(hqs[36], 64)
	now.VOTURNOVER, _ = strconv.ParseFloat(hqs[31], 64)

	return now
}

func (this *Stock) UpdateNow(hqs []string) {
	this.Code = hqs[2]
	this.Name = hqs[1]

	this.Now = NewNowHQ(hqs)
}

func (this *Stocks) UpdateReal() int {
	suc := 0
	i := 0
	size := len(this.stocks)
	for {
		codes := this.getCodes(i, i+C_Step)
		str, _ := GetReal(codes)
		m := SplitRealStr(str)
		suc += len(m)
		for key, val := range m {
			stk := this.stockMap[key]
			stk.UpdateNow(val)
		}
		if i+C_Step >= size {
			break
		} else {
			i += C_Step
		}
	}
	return 0
}

func (this *Stocks) UpdateHistory() int {

	return 0
}

func NewStocks() (*Stocks, error) {
	//首先读取所有Code
	codes, err := GetStockCodes()
	if err != nil {
		return nil, err
	}

	l := len(codes)

	stocks := new(Stocks)
	stocks.stocks = make([]*Stock, l)
	stocks.stockMap = make(map[string]*Stock, l)

	for index, code := range codes {
		stock := NewStock(code)
		stocks.stocks[index] = stock
		stocks.stockMap[code] = stock
	}

	//更新实时行情
	stocks.UpdateReal()
	//更新历史行情
	stocks.UpdateHistory()

	return stocks, nil
}

func NewStock(code string) *Stock {
	s := new(Stock)
	s.Code = code
	return s
}

type Stock struct {
	Code   string
	Name   string
	Now    *DayLine
	Before *DayLines
}

//股票日线行情
type DayLine struct {
	UTIME      time.Time //更新时间
	TCLOSE     float64   //收盘价
	HIGH       float64   //最高价
	LOW        float64   //最低价
	TOPEN      float64   //开盘价
	LCLOSE     float64   //前收盘
	CHG        float64   //涨跌额
	PCHG       float64   //涨跌幅
	TURNOVER   float64   //换手率
	VOTURNOVER float64   //成交量
	VATURNOVER float64   //成交金额
	TCAP       float64   //总市值
	MCAP       float64   //流通市值
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
