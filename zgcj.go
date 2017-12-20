package main

import (
	"../zgcj/sd"
	"fmt"
)

func main() {

	codes, _ := sd.GetStockCodes()
	i := 0
	step:=20
	for {
		e:=i+step
		str, _ := sd.GetReal(codes[i:e])
		fmt.Println(str)
		i = e
		if i >= len(codes) {
			break
		}
	}

}
