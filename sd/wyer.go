package sd

import (
	"net/http"
	"github.com/qjsoftcn/gutils"
	"time"
	"os"
	"io"
	"log"
	"strings"
	"errors"
)

const (
	His_Url = "http://quotes.money.163.com/service/chddata.html?code=${code}&start=${start}&end=${end}&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"
)

func getMarketCode(code string) string {
	switch code[0] {
	case '0', '3':
		return "1"
	case '6':
		return "0"
	default:
		return "-1"
	}
}

func makeUrl(code string, start, end time.Time) string {

	url := strings.Replace(His_Url, "${code}", code, 1)
	url = strings.Replace(url, "${start}", gutils.Format(start, "yyyyMMdd"), 1)
	url = strings.Replace(url, "${end}", gutils.Format(end, "yyyyMMdd"), 1)
	return url
}

func GetHistory(code string, start, end time.Time) ([]Stock, error) {
	mc := getMarketCode(code)
	if mc == "-1" {
		return nil, errors.New("code:" + code + " err")
	}

	url := makeUrl(mc+code, start, end)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("read wy hisdata err:", err)
		return nil, err
	}

	ts := gutils.FormatToSecondForFileName(time.Now())
	xlsFname := dlsFname(code, ts, "csv")

	fpath := dlsFpath(xlsFname)

	f, err := os.Create(fpath)
	if err != nil {
		log.Fatal("ReadDailyLines err:", err)
		return nil, err
	}

	bnum, err := io.Copy(f, res.Body)
	if err != nil {
		log.Fatal("ReadDailyLines err:", err)
		return nil, err
	}

	bnum = bnum / 1024
	log.Println(f.Name(), "read finished,file size:", bnum, "KB")
	return nil, nil
}
