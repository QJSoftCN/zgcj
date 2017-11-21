package sa

import "../sd"


func Toss(ds *sd.DayLines)*sd.DayLines{
	ds.Sorts(sd.ByAmp(false))
	ds=sd.NewDayLines(ds.Get(0,200))

	ds.Sorts(sd.ByNad(false))
	ds=sd.NewDayLines(ds.Get(0,100))

	ds.Sorts(sd.ByVolume(false))
	ds=sd.NewDayLines(ds.Get(0,50))
	return ds
}