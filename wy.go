package main

import (
	"../zgcj/sd"
	"github.com/qjsoftcn/gutils"
)

func main() {
	s, _ := gutils.Parse("20171101", "yyyyMMdd")
	e, _ := gutils.Parse("20171214", "yyyyMMdd")
	sd.GetHistory("002912", s, e)
}
