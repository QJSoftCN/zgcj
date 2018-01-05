package main

import (
	"../zgcj/sd"
	"fmt"
	"github.com/qjsoftcn/gutils"
)

func getOneDemo() {
	codes := make([]string, 1)
	codes[0] = "002129"
	str, _ := sd.GetReal(false, codes[:1])
	m := sd.SplitRealStr(str)
	for index, val := range m[codes[0]] {
		fmt.Println(index, "=", val)
	}
}

func load() {
	stocks, err := sd.NewStocks()
	if err != nil {
		fmt.Println(err)
		return
	}

	code := "603993"
	stk := stocks.Get(code)
	fmt.Println(stk.Now)

	t, _ := gutils.Parse("2017/12/21", "yyyy/MM/dd")
	fmt.Println(t)
}

func b() {
	sd.Backup(400)
}

func main() {
	b()
}
