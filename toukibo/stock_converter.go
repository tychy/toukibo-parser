package toukibo

import "strings"

func StockToNumber(stock string) int {
	// 発行済株式の総数５万株 → 50000
	stock = strings.Replace(stock, "発行済株式の総数", "", -1)
	stock = ZenkakuToHankaku(stock)

	sums := 0
	cur := 0
	for _, v := range stock {
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
		case '株':
			sums += cur
			cur = 0
		}
	}
	return sums
}
