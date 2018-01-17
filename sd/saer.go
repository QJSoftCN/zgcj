package sd

import (
	"fmt"
	"github.com/qjsoftcn/gutils"
	"sort"
	"github.com/tealeg/xlsx"
	"path/filepath"
)

func (this *DayLines) FindLowestDay(days int) (*DayLine, int) {
	size := len(this.dls)
	if size == 0 {
		return nil, 0
	}

	if days >= size {
		days = size
	}

	ldl := &this.dls[0]
	day := 0
	for index, dl := range this.dls {
		if index >= days {
			break
		}

		if dl.LOW < ldl.LOW {
			ldl = &this.dls[index]
			day = index + 1
		}
	}

	return ldl, day
}

type StockLowest struct {
	Code    string
	DevAmp  float64
	DevDays int
	Low     *DayLine
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

	sls := make(StockLowests, 0)
	for _, stk := range this.stocks {
		ldl, day := stk.Before.FindLowestDay(days)
		if ldl != nil && ldl.LOW != 0 {
			pc := (stk.Now.NDAY.TCLOSE - ldl.LOW) / ldl.LOW
			sl := StockLowest{stk.Code, pc, day, ldl}
			sls = append(sls, sl)
		}
	}

	sort.Sort(sls)
	xfile := xlsx.NewFile()

	sheet, _ := xfile.AddSheet("sl")
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.SetString("序号")
	cell = row.AddCell()
	cell.SetString("代码")
	cell = row.AddCell()
	cell.SetString("名称")
	cell = row.AddCell()
	cell.SetString("现价")
	cell = row.AddCell()
	cell.SetString("最低价")
	cell = row.AddCell()
	cell.SetString("最低日期")
	cell = row.AddCell()
	cell.SetString("相对涨幅(%)")
	cell = row.AddCell()
	cell.SetString("调整天数")

	for index, sl := range sls {
		stk := this.Get(sl.Code)
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetInt(index + 1)
		cell = row.AddCell()
		cell.SetString(stk.Code)

		cell = row.AddCell()
		cell.SetString(stk.Name)
		cell = row.AddCell()
		cell.SetFloat (stk.Now.NDAY.TCLOSE)
		cell = row.AddCell()
		cell.SetFloat(sl.Low.LOW)
		cell = row.AddCell()
		cell.SetDate(sl.Low.UTIME)
		cell = row.AddCell()
		cell.SetString(gutils.FormatFloat(sl.DevAmp*100, "2%"))
		cell = row.AddCell()
		cell.SetInt(sl.DevDays)
		if index >= 19 {
			break
		}
	}

	fp := filepath.Join(Root_Dir, OUT_DIR, fmt.Sprint("sls_", days, ".xlsx"))
	xfile.Save(fp)
}
