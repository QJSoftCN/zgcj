package tx

import (
	"net/http"
	"os"
	"io"
	"log"
	"time"
	"path/filepath"
	"errors"
	"github.com/qjsoftcn/gutils"
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
