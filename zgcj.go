package main

import (
	"../zgcj/sd"
	"fmt"
)

func getOneDemo() {
	codes := make([]string, 1)
	codes[0] = "002129"
	str, _ := sd.GetReal(codes[:1])
	m := sd.SplitRealStr(str)
	for index, val := range m[codes[0]] {
		fmt.Println(index, "=", val)
	}
}

func main() {
	stocks, err := sd.NewStocks()
	if err != nil {
		fmt.Println(err)
		return
	}

	for index, stock := range stocks.GetStocks() {
		fmt.Println(index, stock)
	}

}
