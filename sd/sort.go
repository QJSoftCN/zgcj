package sd

type By func(t1, t2 *DayLine) bool

func ByAmp(asc bool) By {
	return func(d1, d2 *DayLine) bool {
		zf1 := d1.GetAmplitude()
		zf2 := d2.GetAmplitude()
		if asc {
			return zf1 < zf2
		} else {
			return zf1 > zf2
		}
	}
}

func ByNad(asc bool) By {
	return func(d1, d2 *DayLine) bool {
		n1 := d1.GetNADeviation()
		n2 := d2.GetNADeviation()
		if asc {
			return n1 < n2
		} else {
			return n1 > n2
		}
	}
}

func ByVolume(asc bool) By {
	return func(d1, d2 *DayLine) bool {
		n1 := d1.Volume
		n2 := d2.Volume
		if asc {
			return n1 < n2
		} else {
			return n1 > n2
		}
	}
}
