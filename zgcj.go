package main

import (
	"../zgcj/sd"
	"log"
	"fmt"
)

func main() {

	//fmt.Println(runtime.NumCPU(),runtime.GOMAXPROCS(1))
	dls,err:=sd.GetDayLines()
	if err!=nil{
		log.Println(err)
	}
	//fmt.Println(runtime.NumCPU(),runtime.GOMAXPROCS(-1))
	for index,dl:=range dls{
		fmt.Println(index,dl)
	}

}
