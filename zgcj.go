package main

import (
	"../zgcj/sd"
	"fmt"
)

func main() {

	codes, _ := sd.GetStockCodes()
	i := 0
	for {
		str, _ := sd.GetReal(codes[i:i+5])
		fmt.Println(str)
		i = i + 5
		if i >= len(codes) {
			break
		}
	}

}
