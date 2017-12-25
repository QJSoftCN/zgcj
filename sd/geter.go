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
	"encoding/csv"
	"strconv"
)

const (
	His_Url      = "http://quotes.money.163.com/service/chddata.html?code=${code}&start=${start}&end=${end}&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP"
	Stock_List   = "http://quote.eastmoney.com/stocklist.html"
	Real_HQ_Url  = "http://qt.gtimg.cn/q=${code}"
	Real_SHQ_Url = "http://qt.gtimg.cn/q=s_${code}"

	Root_Dir      = "zgcj"
	DLS_DIR       = "dls"
	HF_EXT        = "csv"
	Market_SH     = "sh"
	Market_SZ     = "sz"
	Var_Code      = "${code}"
	Var_Start     = "${start}"
	Var_End       = "${end}"
	Code_Regx_Exp = "(6|0|3)[\\d]{5}"

	ENC_CODE = "gbk"
	R_HQ_SEP = ";"
)

var enc = mahonia.NewDecoder(ENC_CODE)

//获取所有股票的代码
func GetStockCodes() ([]string, error) {
	doc, err := goquery.NewDocument(Stock_List)
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(Code_Regx_Exp)
	codes := make([]string, 0)
	doc.Find("div[id='quotesearch'] ul li a").Each(func(i int, cs *goquery.Selection) {
		str := cs.Text()
		str = enc.ConvertString(str)
		code := reg.FindString(str)
		if code != "" {
			codes = append(codes, code)
		}
	})

	return codes, nil
}

func ValidCodes(codes []string) []string {
	suc := 0
	i := 0
	size := len(codes)
	nCodes := make([]string, 0)
	for {
		str, _ := GetReal(Real_SHQ_Url, codes[i:i+C_Step])
		m := SplitRealStr(str)
		suc += len(m)
		for key, val := range m {
			tcap, _ := strconv.ParseFloat(val[len(m[key])-1], 64)
			if tcap != 0 {
				nCodes = append(nCodes, key)
			} else {
				fmt.Println(key)
			}
		}

		if i+C_Step >= size {
			break
		} else {
			i += C_Step
		}
	}

	return nCodes
}

func getMarketCodeForWy(code string) string {
	switch code[0] {
	case '0', '3':
		return "1"
	case '6':
		return "0"
	default:
		return "-1"
	}
}

func getMarketCodeForTx(code string) string {
	switch code[0] {
	case '0', '3':
		return Market_SZ
	case '6':
		return Market_SH
	default:
		return "-1"
	}
}

func makeUrl(code, sd, ed string) string {

	url := strings.Replace(His_Url, Var_Code, code, 1)
	url = strings.Replace(url, Var_Start, sd, 1)
	url = strings.Replace(url, Var_End, ed, 1)
	return url
}

func DirCodesPath() string {
	return filepath.Join(Root_Dir, DLS_DIR, "codes."+HF_EXT)
}

func DirHisPath(code string) string {
	return filepath.Join(Root_Dir, DLS_DIR, code+"."+HF_EXT)
}

func backupCode(code, sd, ed string) bool {
	mc := getMarketCodeForWy(code)
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
	defer f.Close()
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

func BackupCodes() {
	codes, err := GetStockCodes()
	if err != nil {
		log.Println(err)
		return
	}
	//验证编码
	codes = ValidCodes(codes)

	fp := DirCodesPath()
	f, err := os.Create(fp)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return
	}

	cw := csv.NewWriter(f)
	cw.Write(codes)
}

func Backup(days int) {
	//backup codes
	BackupCodes()
	//backup historys
	BackupDays(days)
}

func BackupDays(days int) bool {
	now := time.Now()
	start := now.Add(-40 * 24 * time.Hour)

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

func makeRealUrl(burl string, codes []string) string {
	str := ""
	for _, code := range codes {
		if len(code) == 0 {
			continue
		}
		str += getMarketCodeForTx(code) + code + ","
	}
	str = str[:len(str)-1]
	url := strings.Replace(burl, Var_Code, str, 1)
	return url
}

const (
	RHQ_T = "~"
	RHQ_Q = "\""
)

func SplitRealStr(str string) map[string][]string {
	rhqs := strings.Split(str, R_HQ_SEP)

	s_map := make(map[string][]string)
	for _, rhq := range rhqs {
		f := strings.Index(rhq, RHQ_Q)
		if f != -1 {
			hq := strings.Split(rhq[f+1:len(rhq)-1], RHQ_T)
			s_map[hq[2]] = hq[1:]
		}
	}

	return s_map
}

func GetReal(burl string, code []string) (string, error) {
	url := makeRealUrl(burl, code)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("GetReal 1 ", err)
		return "", err
	}

	str := doc.Text()
	if err != nil {
		log.Println("GetReal 2", err)
		return "", err
	}
	str = enc.ConvertString(str)
	return str, nil
}

func init() {
	gutils.Dir(filepath.Join(Root_Dir, DLS_DIR))
}
