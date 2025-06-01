package toukibo

import (
	"fmt"
	"regexp"
	"strings"
)

func GetStockNumber(s string) (int, string) {
	// 株の位置を探す
	kabuIndex := strings.Index(s, "株")
	if kabuIndex == -1 {
		// 株が見つからない場合
		num, err := ParseJapaneseNumber(s, JapaneseNumberUnits)
		if err != nil {
			return 0, ""
		}
		return num, ""
	}
	
	// 株までの部分を抽出
	numPart := s[:kabuIndex]
	remaining := s[kabuIndex+len("株"):]
	
	// 数値をパース
	num, err := ParseJapaneseNumber(numPart, JapaneseNumberUnits)
	if err != nil {
		return 0, remaining
	}
	
	return num, remaining
}

func GetHoujinStock(stock string) HoujinStock {
	stock = ZenkakuToHankaku(stock)
	stock = trimAllSpace(stock)
	res := HoujinStock{}

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
			res.Preferred = append(res.Preferred, HoujinPreferredStock{Type: "普通株式", Amount: normal})
			continue
		}

		if strings.HasPrefix(stock, "優先株式") {
			stock = strings.Replace(stock, "優先株式", "", -1)
			normal, s := GetStockNumber(stock)
			stock = s
			res.Preferred = append(res.Preferred, HoujinPreferredStock{Type: "優先株式", Amount: normal})
			continue
		}

		// sample1428用のハック
		if strings.HasPrefix(stock, "カブアンド種類株式") {
			stock = strings.Replace(stock, "カブアンド種類株式", "", -1)
			normal, s := GetStockNumber(stock)
			stock = s
			res.Preferred = append(res.Preferred, HoujinPreferredStock{Type: "カブアンド種類株式", Amount: normal})
			continue
		}

		// trim （%d）
		// sample1042用のハック
		stock = regexp.MustCompile(`（[0-9]+）`).ReplaceAllString(stock, "")

		// *優先株式 or *種類株式で始まる場合
		pattern := fmt.Sprintf("([%s]+-[0-9]種優先株式|[%s]+[0-9]種優先株式|[%s]+[0-9]優先株式|[%s]+優先株式|[%s]+種類株式|[%s]+種株式)",
			ZenkakuNoNumberStringPattern, ZenkakuNoNumberStringPattern, ZenkakuNoNumberStringPattern, ZenkakuNoNumberStringPattern, ZenkakuNoNumberStringPattern, ZenkakuNoNumberStringPattern)
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(stock)
		if len(matches) > 0 {
			stock = strings.Replace(stock, matches[1], "", 1)
			num, s := GetStockNumber(stock)
			stock = s
			res.Preferred = append(res.Preferred, HoujinPreferredStock{Type: matches[1], Amount: num})
			continue
		}

		break
	}
	if res.Total == 0 {
		res.Total = res.Sum()
	}

	return res
}

func StockToNumber(stock string) int {
	res := GetHoujinStock(stock)
	return res.Total
}
