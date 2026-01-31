package xtime

// SSFFMM + SSFFMM => SSFFMM
func AddTime(a uint64, b uint64) int64 {
	A := int64(a)
	B := int64(b)

	MMa := A % 100
	MMb := B % 100

	FFa := (A / 100) % 100
	FFb := (B / 100) % 100

	SSa := (A / 10000) % 100
	SSb := (B / 10000) % 100

	MMx := MMa + MMb
	if MMx > 59 {
		FFb++
		MMx %= 60
	}
	FFx := FFa + FFb
	if FFx > 59 {
		SSb++
		FFx %= 60
	}
	SSx := SSa + SSb

	return SSx*10000 + FFx*100 + MMx
}

// SSFFMM - SSFFMM = SSFFMM
func SubTime(a uint64, b uint64) int64 {
	A := int64(a)
	B := int64(b)

	MMa := A % 100
	MMb := B % 100

	FFa := (A / 100) % 100
	FFb := (B / 100) % 100

	SSa := (A / 10000) % 100
	SSb := (B / 10000) % 100

	MMx := MMa - MMb
	if MMx < 0 {
		FFb++
		MMx += 60
	}
	FFx := FFa - FFb
	if FFx < 0 {
		SSb++
		FFx += 60
	}
	SSx := SSa - SSb

	return SSx*10000 + FFx*100 + MMx
}

// SSFFMM - SSFFMM = Seconds
func DiffTime(a uint64, b uint64) int64 {
	A := int64(a)
	B := int64(b)

	MMa := A % 100
	MMb := B % 100

	FFa := (A / 100) % 100
	FFb := (B / 100) % 100

	SSa := (A / 10000) % 100
	SSb := (B / 10000) % 100

	MMx := MMa - MMb
	if MMx < 0 {
		FFb++
		MMx += 60
	}
	FFx := FFa - FFb
	if FFx < 0 {
		SSb++
		FFx += 60
	}
	SSx := SSa - SSb

	return SSx*3600 + FFx*60 + MMx
}
