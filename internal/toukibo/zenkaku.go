package toukibo

import (
	"strings"
)

const (
	ZenkakuZero                  = '０'
	ZenkakuNine                  = '９'
	ZenkakuA                     = 'Ａ'
	ZenkakuZ                     = 'Ｚ'
	ZenkakuSmallA                = 'ａ'
	ZenkakuSmallZ                = 'ｚ'
	ZenkakuSpace                 = '　'
	ZenkakuColon                 = '：'
	ZenkakuSlash                 = '／'
	ZenkakuHyphen                = '－'
	ZenkakuNumberPattern         = `０-９`
	ZenkakuNoNumberStringPattern = `\p{Han}\p{Hiragana}\p{Katakana}Ａ-Ｚａ-ｚA-Za-z＆’，‐．・ー\s　。－、：／`
	ZenkakuStringPattern         = ZenkakuNoNumberStringPattern + ZenkakuNumberPattern + `0-9`
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
		case 57451, 57510:
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
		case 57670:
			sb.WriteRune('龍')
		case 60956:
			sb.WriteRune('媛')
		case 58307:
			sb.WriteRune('邊')
		case 60060:
			sb.WriteRune('茨')
		case 59765:
			sb.WriteRune('榊')
		case 58305:
			sb.WriteRune('角')
		case 58069:
			sb.WriteRune('荒')
		case 60981, 60997:
			sb.WriteRune('藤')
		case 60025:
			sb.WriteRune('邦')
		case 60848:
			sb.WriteRune('楢')
		case 57447:
			sb.WriteRune('橋')
		case 57687:
			sb.WriteRune('邉')
		// Cyrillic homoglyphs → Fullwidth Latin
		// PDF font encoding sometimes maps Latin glyphs to Cyrillic codepoints
		case 'А': // U+0410 → Ａ
			sb.WriteRune('Ａ')
		case 'В': // U+0412 → Ｂ
			sb.WriteRune('Ｂ')
		case 'Е': // U+0415 → Ｅ
			sb.WriteRune('Ｅ')
		case 'К': // U+041A → Ｋ
			sb.WriteRune('Ｋ')
		case 'М': // U+041C → Ｍ
			sb.WriteRune('Ｍ')
		case 'Н': // U+041D → Ｎ
			sb.WriteRune('Ｎ')
		case 'О': // U+041E → Ｏ
			sb.WriteRune('Ｏ')
		case 'Р': // U+0420 → Ｒ
			sb.WriteRune('Ｒ')
		case 'С': // U+0421 → Ｓ
			sb.WriteRune('Ｓ')
		case 'Т': // U+0422 → Ｔ
			sb.WriteRune('Ｔ')
		case 'Х': // U+0425 → Ｘ
			sb.WriteRune('Ｘ')
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
