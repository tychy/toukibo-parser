package toukibo

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

func zenkakuToHankaku(s string) string {
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
