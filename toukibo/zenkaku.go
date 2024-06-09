package toukibo

import (
	"strings"
)

const (
	ZenkakuZero          = '０'
	ZenkakuNine          = '９'
	ZenkakuA             = 'Ａ'
	ZenkakuZ             = 'Ｚ'
	ZenkakuSmallA        = 'ａ'
	ZenkakuSmallZ        = 'ｚ'
	ZenkakuSpace         = '　'
	ZenkakuColon         = '：'
	ZenkakuSlash         = '／'
	ZenkakuHyphen        = '－'
	ZenkakuStringPattern = `\p{Han}\p{Hiragana}\p{Katakana}Ａ-Ｚａ-ｚ０-９A-Za-z0-9＆’，‐．・ー\s　。－、：／`
)

func ZenkakuToHankaku(s string) string {
	var result string
	for _, r := range s {
		if r >= ZenkakuZero && r <= ZenkakuNine {
			result += string(r - ZenkakuZero + '0')
		} else if r >= ZenkakuA && r <= ZenkakuZ {
			result += string(r - ZenkakuA + 'A')
		} else if r >= ZenkakuSmallA && r <= ZenkakuSmallZ {
			result += string(r - ZenkakuSmallA + 'a')
		} else if r == ZenkakuSlash {
			result += "/"
		} else if r == ZenkakuColon {
			result += ":"
		} else if r == ZenkakuSpace {
			result += " "
		} else if r == ZenkakuHyphen {
			result += "-"
		} else {
			result += string(r)
		}
	}
	return result
}

func normalizeKanji(input string) string {
	// https://www.natade.net/webapp/mojicode-kaiseki/
	// https://codepoints.net/
	var sb strings.Builder
	for _, r := range input {
		switch r {
		case 57451:
			sb.WriteRune('塚')
		case 57735:
			sb.WriteRune('西')
		case 60887:
			sb.WriteRune('逢')
		case 57788:
			sb.WriteRune('花')
		case 60906:
			sb.WriteRune('辻')
		case 57374:
			sb.WriteRune('土')
		case 60849:
			sb.WriteRune('樋')
		case 58450:
			sb.WriteRune('廣')
		case 57860:
			sb.WriteRune('若')
		case 59648:
			sb.WriteRune('藤')
		case 59470:
			sb.WriteRune('吉')
		case 60100:
			sb.WriteRune('蛸')
		case 63964, 59478:
			sb.WriteRune('隆')
		case 59424:
			sb.WriteRune('座')
		case 60939:
			sb.WriteRune('那')
		case 59911:
			sb.WriteRune('覇')
		case 59620:
			sb.WriteRune('徽')
		case 61033:
			sb.WriteRune('聰')
		case 60059:
			sb.WriteRune('芦')
		case 59788:
			sb.WriteRune('禮')
		case 59677:
			sb.WriteRune('原')
		case 59859:
			sb.WriteRune('牙')
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
