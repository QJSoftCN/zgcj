package sd

import (
	"time"
	"github.com/qjsoftcn/gutils"
	"sort"
	"strconv"
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"io/ioutil"
	"strings"
)

const (
	C_Step int = 60
	P_SH       = '6'
)

type Stocks struct {
	stocks   []*Stock
	stockMap map[string]*Stock
}

func (this *Stocks) Get(code string) *Stock {
	return this.stockMap[code]
}

func (this *Stocks) GetStocks() []*Stock {
	return this.stocks
}

func (this *Stocks) getCodes(s, e int) []string {
	codes := make([]string, e-s+1)
	tl := len(this.stocks)
	if e > tl {
		e = tl
	}
	ss := this.stocks[s:e]
	for index, code := range ss {
		codes[index] = code.Code
	}
	return codes

}

func newNDay(hqs []string) *DayLine {
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
	now.TCLOSE, _ = strconv.ParseFloat(hqs[2], 64)
	now.VATURNOVER, _ = strconv.ParseFloat(hqs[36], 64)
	now.VOTURNOVER, _ = strconv.ParseFloat(hqs[35], 64)
	now.UTIME, _ = gutils.Parse(hqs[29][:8], "yyyyMMdd")
	return now
}

func (this *Stock) UpdateNow(hqs []string) {
	this.Name = hqs[0]
	this.Now = NewNowLine(hqs)
}

func newDayLine(hqs []string) *DayLine {
	now := new(DayLine)
	now.CHG, _ = strconv.ParseFloat(hqs[8], 64)
	now.HIGH, _ = strconv.ParseFloat(hqs[4], 64)
	now.LCLOSE, _ = strconv.ParseFloat(hqs[7], 64)
	now.LOW, _ = strconv.ParseFloat(hqs[5], 64)
	now.MCAP, _ = strconv.ParseFloat(hqs[14], 64)
	now.PCHG, _ = strconv.ParseFloat(hqs[9], 64)
	now.TCAP, _ = strconv.ParseFloat(hqs[13], 64)

	now.TOPEN, _ = strconv.ParseFloat(hqs[6], 64)
	now.TURNOVER, _ = strconv.ParseFloat(hqs[10], 64)
	now.TCLOSE, _ = strconv.ParseFloat(hqs[3], 64)
	now.VATURNOVER, _ = strconv.ParseFloat(hqs[12], 64)
	now.VOTURNOVER, _ = strconv.ParseFloat(hqs[11], 64)
	now.UTIME, _ = gutils.Parse(hqs[0], "yyyy-MM-dd")
	return now
}

func (this *Stock) ReadHistory() bool {
	fp := DirHisPath(this.Code)
	f, err := os.Open(fp)
	if err != nil {
		log.Println("read history code: ", this.Code, " err:", err)
		return false
	}

	r := csv.NewReader(f)
	rs, err := r.ReadAll()
	rl := len(rs)
	if rl <= 1 {
		log.Println("read history code: ", this.Code, " no history")
		return false
	}

	dls := make([]DayLine, len(rs)-1)
	for i, dl := range rs[1:] {
		dls[i] = *newDayLine(dl)
	}
	hs := new(DayLines)
	hs.dls = dls
	this.Before = hs
	return true
}

func NewNowLine(hqs []string) *NowLine {
	nl := new(NowLine)
	nl.NDAY = newNDay(hqs)
	nl.HLIMIT, _ = strconv.ParseFloat(hqs[46], 64)
	nl.LLIMIT, _ = strconv.ParseFloat(hqs[47], 64)
	nl.PEV, _ = strconv.ParseFloat(hqs[38], 64)
	nl.AMP, _ = strconv.ParseFloat(hqs[42], 64)
	nl.PBV, _ = strconv.ParseFloat(hqs[45], 64)

	return nl
}

func (this Stock) String() string {
	return fmt.Sprintln(this.Code, this.Name, this.Now, this.Before.Size())
}

func (this *Stocks) UpdateNow() int {
	suc := 0
	i := 0
	size := len(this.stocks)
	for {
		codes := this.getCodes(i, i+C_Step)
		str, _ := GetReal(false,codes)
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
	return suc
}

func (this *Stocks) UpdateHistory() int {
	suc := 0
	for _, stock := range this.stocks {
		isSuc := stock.ReadHistory()
		if isSuc {
			suc++
		}
	}

	return suc
}

func ReadCodes() ([]string, error) {
	fp := DirCodesPath()
	f, err := os.Open(fp)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bs,err:=ioutil.ReadAll(f)
	if err!=nil{
		log.Println(err)
		return nil, err
	}
	str:=string(bs)
	codes:=strings.Split(str," ")
	return codes, nil
}

func NewStocks() (*Stocks, error) {
	//首先读取所有Code
	codes, err := ReadCodes()
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
	suc := stocks.UpdateNow()
	log.Println("update now lines success:", suc, " total:", l)
	//更新历史行情
	suc = stocks.UpdateHistory()
	log.Println("update history lines success:", suc, " total:", l)

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
	state  int
	Now    *NowLine
	Before *DayLines
}

type NowLine struct {
	NDAY   *DayLine
	HLIMIT float64 //涨停价
	LLIMIT float64 //跌停价
	PEV    float64 //市盈率
	PBV    float64 //市净率
	AMP    float64 //振幅
}

func (this NowLine) String() string {
	return fmt.Sprintln(this.NDAY,
		gutils.FormatFloat(this.HLIMIT, "2"),
		gutils.FormatFloat(this.LLIMIT, "2"),
		gutils.FormatFloat(this.PEV, "2"),
		gutils.FormatFloat(this.PBV, "2"),
		gutils.FormatFloat(this.AMP, "2%"))
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

func (this DayLine) String() string {

	return fmt.Sprintln(gutils.Format(this.UTIME, "yyyyMMdd"),
		gutils.FormatFloat(this.TCLOSE, "2"),
		gutils.FormatFloat(this.HIGH, "2"),
		gutils.FormatFloat(this.LOW, "2"),
		gutils.FormatFloat(this.TOPEN, "2"),
		gutils.FormatFloat(this.LCLOSE, "2"),
		gutils.FormatFloat(this.CHG, "2"),
		gutils.FormatFloat(this.PCHG, "2%"),
		gutils.FormatFloat(this.TURNOVER, "2%"),
		this.VOTURNOVER, gutils.FormatFloat(this.VATURNOVER, "2"),
		gutils.FormatFloat(this.TCAP, "2"),
		gutils.FormatFloat(this.MCAP, "2"))
}

type DayLines struct {
	dls []DayLine
	by  By
}

func (this *DayLines) Size() int {
	return len(this.dls)
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
	return code[0] == P_SH
}
