package toukibo

// 1万円 -> 10000
// １億円 -> 100000000
// １兆円 -> 1000000000000
// 30億2000万1000円 -> 3020001000

// 上限は兆円まで
func YenToNumber(yen string) int {
	yen = ZenkakuToHankaku(yen)

	sums := 0
	cur := 0
	for _, v := range yen {
		if v == '金' {
			continue
		}
		if v >= '0' && v <= '9' {
			cur = cur*10 + int(v-'0')
			continue
		}

		switch v {
		case '万':
			sums += cur * 10000
			cur = 0
		case '億':
			sums += cur * 100000000
			cur = 0
		case '兆':
			sums += cur * 1000000000000
			cur = 0
		case '円':
			sums += cur
			cur = 0
		}
	}
	return sums
}
