package main

import (
	"../zgcj/sd"
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
	fmt.Println("振幅最大，平均获利指数最好的个股")
	ds.Sorts(sd.GetByAmp(false), sd.GetByNad(false),sd.GetByVolume(false))
	for index, t := range ds.GetDLS() {
		if index > 20 {
			break
		}
		fmt.Println(index, t.Name, gutils.FormatFloat(t.LatestPrice, "2"),
			gutils.FormatFloat(t.GetAmplitude()*100, "%2"),
			gutils.FormatFloat(t.GetNADeviation()*100, "%2")+"%")
	}

	fmt.Println("振幅最大，平均获利指数最差的个股")
	//振幅最大，平均获利指数最差的个股
	ds.Sorts(sd.GetByAmp(false), sd.GetByNad(true),sd.GetByVolume(false))
	for index, t := range ds.GetDLS() {
		if index > 20 {
			break
		}
		fmt.Println(index, t.Name, gutils.FormatFloat(t.LatestPrice, "2"),
			gutils.FormatFloat(t.GetAmplitude()*100, "%2"),
			gutils.FormatFloat(t.GetNADeviation()*100, "%2")+"%")
	}

}
