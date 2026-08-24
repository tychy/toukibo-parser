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
		case 59764: // U+E974
			sb.WriteRune('櫛')
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
		case 61066:
			sb.WriteRune('四')
		case 57421: // U+E04D
			sb.WriteRune('墅')
		case 57436: // U+E05C
			sb.WriteRune('橋')
		case 57443: // U+E063
			sb.WriteRune('齊')
		case 57532: // U+E0BC
			sb.WriteRune('脇')
		case 57675: // U+E14B
			sb.WriteRune('今')
		case 57765: // U+E1A5
			sb.WriteRune('鄉')
		case 57766: // U+E1A6
			sb.WriteRune('㋐')
		case 57767: // U+E1A7
			sb.WriteRune('㋑')
		case 57768: // U+E1A8
			sb.WriteRune('㋒')
		case 57787: // U+E1BB
			sb.WriteRune('㋥')
		case 57790: // U+E1BE
			sb.WriteRune('芳')
		case 57791: // U+E1BF
			sb.WriteRune('㋩')
		case 57807: // U+E1CF
			sb.WriteRune('㋺')
		case 57831: // U+E1E7
			sb.WriteRune('🄰')
		case 57832: // U+E1E8
			sb.WriteRune('🄱')
		case 57833: // U+E1E9
			sb.WriteRune('🄲')
		case 57866: // U+E20A
			sb.WriteRune('汇')
		case 57872: // U+E210
			sb.WriteRune('英')
		case 57897: // U+E229
			sb.WriteRune('茂')
		case 57928: // U+E248
			sb.WriteRune('茨')
		case 57988: // U+E284
			sb.WriteRune('兒')
		case 58128: // U+E310
			sb.WriteRune('渣')
		case 58133: // U+E315
			sb.WriteRune('愼')
		case 58162: // U+E332
			sb.WriteRune('梁')
		case 58263: // U+E397
			sb.WriteRune('菂')
		case 58267: // U+E39B
			sb.WriteRune('菅')
		case 58277: // U+E3A5
			sb.WriteRune('菊')
		case 58291: // U+E3B3
			sb.WriteRune('菟')
		case 58486: // U+E476
			sb.WriteRune('萩')
		case 58536: // U+E4A8
			sb.WriteRune('葉')
		case 58541: // U+E4AD
			sb.WriteRune('熙')
		case 58616: // U+E4F8
			sb.WriteRune('鯉')
		case 58789: // U+E5A5
			sb.WriteRune('喜')
		case 58807: // U+E5B7
			sb.WriteRune('初')
		case 58941: // U+E63D
			sb.WriteString("ⅩⅢ")
		case 58966: // U+E656
			sb.WriteRune('紀')
		case 59071: // U+E6BF
			sb.WriteRune('３')
		case 59073: // U+E6C1
			sb.WriteRune('５')
		case 59074: // U+E6C2
			sb.WriteRune('６')
		case 59084: // U+E6CC
			sb.WriteRune('１')
		case 59085: // U+E6CD
			sb.WriteRune('２')
		case 59308: // U+E7AC
			sb.WriteRune('蔭')
		case 59372: // U+E7EC
			sb.WriteString("σ²")
		case 59373: // U+E7ED
			sb.WriteString("√Ｔ")
		case 59374: // U+E7EE
			sb.WriteString("－ｑＴ")
		case 59375: // U+E7EF
			sb.WriteString("－ｒＴ")
		case 59439: // U+E82F
			sb.WriteRune('兎')
		case 59472: // U+E850
			sb.WriteRune('啓')
		case 59500: // U+E86C
			sb.WriteRune('厩')
		case 59553: // U+E8A1
			sb.WriteRune('藏')
		case 59554: // U+E8A2
			sb.WriteRune('ｎ')
		case 59564: // U+E8AC
			sb.WriteRune('屑')
		case 59565: // U+E8AD
			sb.WriteRune('屠')
		case 59573: // U+E8B5
			sb.WriteRune('惠')
		case 59609: // U+E8D9
			sb.WriteRune('彦')
		case 59646: // U+E8FE
			sb.WriteString("－ｒｔ")
		case 59647: // U+E8FF
			sb.WriteRune('赳')
		case 59649: // U+E901
			sb.WriteString("√ｔ")
		case 59650: // U+E902
			sb.WriteString("－λｔ")
		case 59661: // U+E90D
			sb.WriteRune('捲')
		case 59668: // U+E914
			sb.WriteRune('眞')
		case 59678: // U+E91E
			sb.WriteRune('揃')
		case 59701: // U+E935
			sb.WriteRune('浩')
		case 59712: // U+E940
			sb.WriteRune('港')
		case 59724: // U+E94C
			sb.WriteRune('厩')
		case 59749: // U+E965
			sb.WriteRune('厩')
		case 59750: // U+E966
			sb.WriteRune('﨑')
		case 59761: // U+E971
			sb.WriteRune('漆')
		case 59853: // U+E9CD
			sb.WriteRune('煎')
		case 59854: // U+E9CE
			sb.WriteString("－ｑｔ")
		case 59867: // U+E9DB
			sb.WriteString("－λＴ")
		case 59885: // U+E9ED
			sb.WriteRune('藤')
		case 59929: // U+EA19
			sb.WriteRune('祇')
		case 59951: // U+EA2F
			sb.WriteRune('起')
		case 59952: // U+EA30
			sb.WriteRune('稗')
		case 59989: // U+EA55
			sb.WriteRune('進')
		case 59991: // U+EA57
			sb.WriteRune('箭')
		case 59992: // U+EA58
			sb.WriteRune('箸')
		case 60013: // U+EA6D
			sb.WriteRune('邦')
		case 60014: // U+EA6E
			sb.WriteRune('篇')
		case 60018: // U+EA72
			sb.WriteRune('籾')
		case 60063: // U+EA9F
			sb.WriteRune('葛')
		case 60065: // U+EAA1
			sb.WriteRune('薩')
		case 60075: // U+EAAB
			sb.WriteRune('爌')
		case 60077: // U+EAAD
			sb.WriteRune('細')
		case 60097: // U+EAC1
			sb.WriteRune('蔽')
		case 60099: // U+EAC3
			sb.WriteRune('蓬')
		case 60109: // U+EACD
			sb.WriteRune('襖')
		case 60140: // U+EAEC
			sb.WriteRune('こ')
		case 60176: // U+EB10
			sb.WriteRune('満')
		case 60201: // U+EB29
			sb.WriteRune('靜')
		case 60202: // U+EB2A
			sb.WriteRune('江')
		case 60264: // U+EB68
			sb.WriteRune('津')
		case 60345: // U+EBB9
			sb.WriteRune('石')
		case 60379: // U+EBDB
			sb.WriteRune('器')
		case 60429: // U+EC0D
			sb.WriteRune('鄭')
		case 60430: // U+EC0E
			sb.WriteRune('錆')
		case 60432: // U+EC10
			sb.WriteRune('鑓')
		case 60480: // U+EC40
			sb.WriteRune('鞄')
		case 60481: // U+EC41
			sb.WriteRune('鞘')
		case 60488: // U+EC48
			sb.WriteRune('餌')
		case 60511: // U+EC5F
			sb.WriteRune('曽')
		case 60564: // U+EC94
			sb.WriteRune('英')
		case 60604: // U+ECBC
			sb.WriteRune('）')
		case 60606: // U+ECBE
			sb.WriteRune('）')
		case 60608: // U+ECC0
			sb.WriteRune('（')
		case 60610: // U+ECC2
			sb.WriteRune('（')
		case 60617: // U+ECC9
			sb.WriteRune('＋')
		case 60618: // U+ECCA
			sb.WriteRune('／')
		case 60621: // U+ECCD
			sb.WriteRune('○')
		case 60657: // U+ECF1
			sb.WriteRune('餅')
		case 60672: // U+ED00
			sb.WriteRune('）')
		case 60673: // U+ED01
			sb.WriteRune('（')
		case 60736: // U+ED40
			sb.WriteRune('幸')
		case 60833: // U+EDA1
			sb.WriteRune('厩')
		case 60834: // U+EDA2
			sb.WriteRune('鵠')
		case 60842: // U+EDAA
			sb.WriteRune('噌')
		case 60845: // U+EDAD
			sb.WriteRune('挽')
		case 60846: // U+EDAE
			sb.WriteRune('樽')
		case 60851: // U+EDB3
			sb.WriteRune('汲')
		case 60855: // U+EDB7
			sb.WriteRune('瀞')
		case 60856: // U+EDB8
			sb.WriteRune('灘')
		case 60865: // U+EDC1
			sb.WriteRune('迦')
		case 60868: // U+EDC4
			sb.WriteRune('逗')
		case 60869: // U+EDC5
			sb.WriteRune('遡')
		case 60871: // U+EDC7
			sb.WriteRune('溝')
		case 60886: // U+EDD6
			sb.WriteRune('曙')
		case 60888: // U+EDD8
			sb.WriteRune('虻')
		case 60892: // U+EDDC
			sb.WriteRune('梢')
		case 60903: // U+EDE7
			sb.WriteRune('舛')
		case 60909: // U+EDED
			sb.WriteRune('片')
		case 60916: // U+EDF4
			sb.WriteRune('巽')
		case 60918: // U+EDF6
			sb.WriteRune('棚')
		case 60955: // U+EE1B
			sb.WriteRune('梍')
		case 60959: // U+EE1F
			sb.WriteRune('滕')
		case 60964: // U+EE24
			sb.WriteRune('炳')
		case 60983: // U+EE37
			sb.WriteRune('蓮')
		case 60996: // U+EE44
			sb.WriteRune('廣')
		case 61014: // U+EE56
			sb.WriteRune('翔')
		case 61020: // U+EE5C
			sb.WriteRune('桒')
		case 61030: // U+EE66
			sb.WriteRune('藤')
		case 61031: // U+EE67
			sb.WriteRune('藤')
		case 61035: // U+EE6B
			sb.WriteRune('蘭')
		case 61047: // U+EE77
			sb.WriteRune('辰')
		case 61057: // U+EE81
			sb.WriteRune('鞆')
		case 61164: // U+EEEC
			sb.WriteRune('芒')
		case 61170: // U+EEF2
			sb.WriteRune('泰')
		case 61284: // U+EF64
			sb.WriteRune('均')
		case 61286: // U+EF66
			sb.WriteRune('達')
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
