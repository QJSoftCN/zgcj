package main

import (
	"../zgcj/sd"
	"../zgcj/sa"
	"log"
	"fmt"
	"github.com/qjsoftcn/gutils"
)

func main() {

	dls, err := sd.GetDayLines()
	if err != nil {
		log.Println(err)
	}

	ds := sd.NewDayLines(dls)

	//振幅最大，平均获利指数最高的个股
	fmt.Println("振幅最大，平均获利指数最好且成交量较大的个股")
	ds=sa.Toss(ds)
	for index, t := range ds.Get(0,10) {

		fmt.Println(index, t.Name, gutils.FormatFloat(t.LatestPrice, "2"),
			gutils.FormatFloat(t.GetAvgPrice(), "2"),
			gutils.FormatFloat(t.GetAmplitude()*100, "%2")+"%",
			gutils.FormatFloat(t.GetNADeviation()*100, "%2")+"%")
	}



}
