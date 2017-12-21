package main

import (
	"../zgcj/sd"
	"fmt"
)

func main() {

	codes:=make([]string,1)
	codes[0]="603993"
	str,_:=sd.GetReal(codes[:1])
	m:=sd.SplitRealStr(str)
	for index,val:=range m[codes[0]]{
		fmt.Println(index,"=",val)
	}

}
