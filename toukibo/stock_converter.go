package toukibo

import (
	"fmt"
	"regexp"
	"strings"
)

func GetStockNumber(s string) (int, string) {
	sums := 0
	cur := 0
	foundKabu := false
	for idx, v := range s {
		if foundKabu {
			return sums, s[idx:]
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
		case '株':
			sums += cur
			cur = 0 // 意味はない
			foundKabu = true
			// 発行済株式の総数４万８２４９株各種の株式の数普通株式　　　３万２４９株Ａ種優先株式　１万株Ｂ種優先株式　８０００株
			//　のようなパターンでは最初の株でReturnさせる
			continue
		}
	}
	return sums, ""
}

type PreferredStock struct {
	Name string
	Num  int
}

type Stock struct {
	Total     int
	Preferred []PreferredStock
}

func (s Stock) Sum() int {
	sum := 0
	for _, p := range s.Preferred {
		sum += p.Num
	}
	return sum
}

func StockToNumber(stock string) int {
	stock = ZenkakuToHankaku(stock)
	stock = trimAllSpace(stock)
	res := Stock{}

	for {
		if stock == "" {
			break
		}

		if strings.HasPrefix(stock, "発行済株式の総数") {
			stock = strings.Replace(stock, "発行済株式の総数", "", -1)
			if strings.HasPrefix(stock, "普通株式") {
				continue // sample1082用のハック
			}

			sums, s := GetStockNumber(stock)
			stock = s
			res.Total = sums
			continue
		}

		if strings.HasPrefix(stock, "各種の株式の数") {
			stock = strings.Replace(stock, "各種の株式の数", "", -1)
			continue
		}

		if strings.HasPrefix(stock, "普通株式") {
			stock = strings.Replace(stock, "普通株式", "", -1)
			normal, s := GetStockNumber(stock)
			stock = s
			res.Preferred = append(res.Preferred, PreferredStock{Name: "普通株式", Num: normal})
			continue
		}

		// *優先株式　で始まる場合
		pattern := fmt.Sprintf("([%s]+優先株式)", ZenkakuNoNumberStringPattern)
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(stock)
		if len(matches) > 0 {
			stock = strings.Replace(stock, matches[1], "", -1)
			num, s := GetStockNumber(stock)
			stock = s
			res.Preferred = append(res.Preferred, PreferredStock{Name: matches[1], Num: num})
			continue
		}

		break
	}
	if res.Total == 0 {
		res.Total = res.Sum()
	}

	return res.Total
}
