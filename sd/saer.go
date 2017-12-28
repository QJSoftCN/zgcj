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

func (this *Stocks) Prize() {


}
