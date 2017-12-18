package sd

import (
	"net/http"
	"github.com/qjsoftcn/gutils"
	"time"
	"os"
	"io"
	"log"
	"strings"
	"github.com/puerkitobio/goquery"
	"github.com/axgle/mahonia"
	"regexp"
	"fmt"
	"path/filepath"
)

const (
	His_Url    = "http://quotes.money.163.com/service/chddata.html?code=${code}&start=${start}&end=${end}&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"
	Stock_List = "http://quote.eastmoney.com/stocklist.html"

	Root_Dir = "zgcj"
	HF_EXT   = "csv"
)

//获取所有股票的代码
func GetStockCodes() ([]string, error) {
	doc, err := goquery.NewDocument(Stock_List)
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile("(6|0|3)[\\d]{5}")
	codes := make([]string, 0)
	doc.Find("div[id='quotesearch'] ul li a").Each(func(i int, cs *goquery.Selection) {
		str := cs.Text()
		enc := mahonia.NewDecoder("gbk")
		str = enc.ConvertString(str)

		code := reg.FindString(str)
		if code != "" {
			codes = append(codes, code)
		}
	})

	return codes, nil

}

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

func makeUrl(code, sd, ed string) string {

	url := strings.Replace(His_Url, "${code}", code, 1)
	url = strings.Replace(url, "${start}", sd, 1)
	url = strings.Replace(url, "${end}", ed, 1)
	return url
}

func DirHisPath(code string) string {
	return filepath.Join(Root_Dir, DLS_DIR, code+"."+HF_EXT)
}

func backupCode(code, sd, ed string) bool {
	mc := getMarketCode(code)
	if mc == "-1" {
		return false
	}

	url := makeUrl(mc+code, sd, ed)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(code, err)
		return false
	}

	bf := DirHisPath(code)
	f, err := os.Create(bf)
	if err != nil {
		log.Fatal(code, err)
		return false
	}

	bnum, err := io.Copy(f, res.Body)
	if err != nil {
		log.Fatal(code, err)
		return false
	}

	bnum = bnum / 1024
	log.Println(code, f.Name(), bnum, "KB")
	return true
}

func BackupDays(days int) bool {
	now := time.Now()
	dur := time.Duration(-days) * time.Duration(24*time.Hour)
	start := now.Add(dur)

	codes, err := GetStockCodes()
	if err != nil {
		log.Fatal("read codes ", err)
		return false
	}

	startDay := gutils.Format(start, "yyyyMMdd")
	endDay := gutils.Format(now, "yyyyMMdd")

	failedCodes := make([]string, 0)

	for index, code := range codes {
		isSuc := backupCode(code, startDay, endDay)
		if isSuc {
			fmt.Println(index, code, " backup success")
		} else {
			failedCodes = append(failedCodes, code)
		}
	}

	if len(failedCodes) == 0 {
		return true
	} else {
		log.Println(failedCodes)
		return false
	}
}

func init() {
	gutils.Dir(filepath.Join(Root_Dir, DLS_DIR))
}
