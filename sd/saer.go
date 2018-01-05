package sd

type SPrizer struct {
	Days  int
	High  float64
	Low   float64
	Avg   float64
	State int
}

type SPrizers struct {
	Code string
	D5   *SPrizer
	D10  *SPrizer
	D20  *SPrizer
	D30  *SPrizer
	D60  *SPrizer
	D400 *SPrizer
}

type SVolumner struct {
	Days  int
	High int64
	Low   int64
	Avg   float64
	State int
}

type SVolumners struct {
	Code string
	V3 *SVolumner
	V5 *SVolumner
	V10 *SVolumner
	V15 *SVolumner
}

func(this Stock)Prize(){

}

var allSPrizers=make(map[string]SPrizers)
//build all prizers
//at this prizes to analysis prize
func (this *Stocks) Prize(){


}

