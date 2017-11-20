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
	ts := ds.Toss()

	for index, t := range ts.GetTosses() {
		if index > 20 {
			break
		}
		fmt.Println(index, t.Day.Name, gutils.FormatFloat(t.Day.LatestPrice, "2"), gutils.FormatFloat(t.Day.GetAmplitude()*100, "%2"))
	}

}
