package sd

import (
	"fmt"
	"github.com/qjsoftcn/gutils"
)

func (this *DayLines) FindLowestDay(days int) (*DayLine,int) {
	size:=len(this.dls)
	if size== 0 {
		return nil,0
	}

	if days>=size{days=size}

	ldl := &this.dls[0]
	day:=0
	for index, dl := range this.dls {
		if index >= days {
			break
		}

		if dl.LOW < ldl.LOW {
			ldl = &this.dls[index]
			day=index+1
		}
	}

	return ldl,day
}
type StockLowest struct {
	Code string
	DevAmp float64
	DevDays int
	Low *DayLine
}

type StockLowests []StockLowest

func (c StockLowests) Len() int {
	return len(c)
}
func (c StockLowests) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c StockLowests) Less(i, j int) bool {
	return c[i].DevAmp < c[j].DevAmp
}
//跌幅最大
func (this *Stocks) DropMost(days int) {

	sls:=make(StockLowests,0)
	for _, stk := range this.stocks {
		ldl,day := stk.Before.FindLowestDay(days)
		if ldl!=nil&&ldl.LOW!=0 {
			pc := (stk.Now.NDAY.TCLOSE - ldl.LOW) / ldl.LOW
			sl := StockLowest{stk.Code, pc, day,ldl}
			sls = append(sls, sl)
		}
	}

	for index,sl:=range sls{
		stk:=this.Get(sl.Code)
		if index>20{
			break
		}
		fmt.Println(index,stk.Code,stk.Name,stk.Now.NDAY.TCLOSE,gutils.FormatFloat(sl.DevAmp,"%.2") ,sl.DevDays,sl.Low.LOW,gutils.Format(sl.Low.UTIME,"yyyyMMdd"))
	}


}


