package sd

import (
	"net/http"
	"os"
	"io"
	"log"
	"time"
	"path/filepath"
	"errors"
	"github.com/qjsoftcn/gutils"
	"github.com/tealeg/xlsx"
	"strconv"
)

const (
	DailyLine_URL = "http://stock.gtimg.cn/data/get_hs_xls.php?id=rankash&type=1&metric=chr"
	DLS_DIR       = "dls"
)

func initDLSDir() {
	if !gutils.PathExists(DLS_DIR) {
		gutils.Dir(DLS_DIR)
	}
}

func init() {
	initDLSDir()
	log.Println("Init DLS Dir")
}

func dlsFname(timestamp, ext string) string {
	return "tx_dls_" + timestamp + "." + ext
}

func dlsFpath(fname string) string {
	return filepath.Join(DLS_DIR, fname)
}

func ReadDailyLines() (string, error) {
	res, err := http.Get(DailyLine_URL)
	if err != nil {
		log.Fatal("ReadDailyLines err:", err)
		return "", err
	}

	ts := gutils.FormatToSecondForFileName(time.Now())
	xlsFname := dlsFname(ts, "xls")

	fpath := dlsFpath(xlsFname)

	f, err := os.Create(fpath)
	if err != nil {
		log.Fatal("ReadDailyLines err:", err)
		return "", err
	}

	bnum, err := io.Copy(f, res.Body)
	if err != nil {
		log.Fatal("ReadDailyLines err:", err)
		return "", err
	}

	bnum = bnum / 1024
	log.Println(f.Name(), "read finished,file size:", bnum, "KB")

	if bnum > 0 {
		//cast to xlsx
		xlsxFname := dlsFname(ts, "xlsx")
		xfpath := dlsFpath(xlsxFname)
		absFpath, _ := filepath.Abs(fpath)
		absxFpath, _ := filepath.Abs(xfpath)

		isSuc := gutils.ToXlsx(absFpath, absxFpath)
		if isSuc {
			log.Println(xfpath, " convert success")
			return xfpath, nil
		} else {
			log.Fatal(xfpath, " convert failed")
			return "", errors.New("convert failed")
		}
	}
	return "", errors.New("empty file")

}

func GetDayLines() ([]DayLine, error) {
	dlsfile, err := ReadDailyLines()
	//读取日线数据
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ReadXlsx(dlsfile)
}

func ReadXlsx(xlsxFilePath string) ([]DayLine, error) {
	file, err := xlsx.OpenFile(xlsxFilePath)
	if err != nil {
		return nil, err
	}

	sheet := file.Sheets[0]

	row := sheet.Rows[0]
	//fmt:11-17 15:01:33
	dupCell := row.Cells[1]
	year := time.Now().Year()
	dup := strconv.Itoa(year) + "-" + dupCell.Value
	ut, _ := time.Parse(dup, "yyyy-MM-dd HH:mm:ss")

	dls := make([]DayLine, 0)
	for ri := 2; ri < sheet.MaxRow; ri++ {
		row := sheet.Rows[ri]

		dl := new(DayLine)
		dl.Code = row.Cells[0].Value
		dl.Name = row.Cells[1].Value
		dl.LatestPrice, _ = row.Cells[2].Float()
		dl.QuoteChange, _ = row.Cells[3].Float()
		dl.AmountChange, _ = row.Cells[4].Float()
		dl.Buy1Price, _ = row.Cells[5].Float()
		dl.Sell1Price, _ = row.Cells[6].Float()
		dl.Volume, _ = row.Cells[7].Int64()
		dl.TurnOver, _ = row.Cells[8].Float()
		dl.NowOpenPrice, _ = row.Cells[9].Float()
		dl.LastClosedPrice, _ = row.Cells[10].Float()
		dl.HighestPrice, _ = row.Cells[11].Float()
		dl.LowestPrice, _ = row.Cells[12].Float()
		dl.UpdateTime = ut

		dls = append(dls, *dl)
	}

	return dls, err
}
