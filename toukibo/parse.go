package toukibo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	beginContent = "┏━━━━━━━━┯━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓"
	endContent   = "┗━━━━━━━━┷━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛"
	separator1   = "┠────────┼─────────────────────────────────────┨"
	separator2   = "┣━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫"
	separator3   = "┣━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━┫"
	separator4   = "┣━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━┫"
	separator5   = "┠────────┼───────────────────────┴─────────────┨"
	separator6   = "┠────────┼───────────────────────┬─────────────┨"
	separator7   = "┠────────┼───────────────────────┼─────────────┨"
	separator8   = "┣━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━━━━━━┫"
)

func findBeginContent(content string) (int, error) {
	index := strings.Index(content, beginContent)
	if index == -1 {
		return 0, fmt.Errorf("not found begin content")
	}
	return index, nil
}

func findEndContent(content string) (int, error) {
	index := strings.Index(content, endContent)
	if index == -1 {
		return 0, fmt.Errorf("not found end content")
	}
	return index, nil
}

func GetHeaderAndContent(content string) (string, string, error) {
	beginContentIdx, err := findBeginContent(content)
	if err != nil {
		return "", "", err
	}

	endContentIdx, err := findEndContent(content)
	if err != nil {
		return "", "", err
	}

	return content[:beginContentIdx], content[beginContentIdx+len(beginContent) : endContentIdx], nil
}

func DivideToukiboContent(input string) (string, []string, error) {
	header, content, err := GetHeaderAndContent(input)
	if err != nil {
		return "", nil, err
	}

	separatorPattern := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s", separator1, separator2, separator3, separator4, separator5, separator6, separator7, separator8)
	re := regexp.MustCompile(separatorPattern)

	parts := re.Split(content, -1)
	return header, parts, nil
}

type Houjin struct {
	header *HoujinHeader
	body   *HoujinBody
}

func (h *Houjin) GetToukiboCreatedAt() time.Time {
	return h.header.CreatedAt
}

func (h *Houjin) GetHoujinKaku() string {
	return string(h.body.HoujinKaku)
}

func (h *Houjin) GetHoujinName() string {
	return h.header.CompanyName
}

func (h *Houjin) GetHoujinAddress() string {
	return h.header.CompanyAddress
}

func (h *Houjin) GetHoujinCreatedAt() string {
	return h.body.HoujinCreatedAt
}

func (h *Houjin) GetHoujinBankruptedAt() string {
	return h.body.HoujinBankruptedAt
}

func (h *Houjin) GetHoujinDissolvedAt() string {
	return h.body.HoujinDissolvedAt
}

func (h *Houjin) GetHoujinCapital() int {
	for _, v := range h.body.HoujinCapital {
		if v.IsValid {
			res := v.Value
			if len(res) < 1 {
				return 0
			}
			return YenToNumber(v.Value)
		}
	}
	return 0
}

func normalize_kanji(input string) string {
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
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func (h *Houjin) ListHoujinExecutives() ([]string, error) {
	execs, err := h.body.ListHoujinExecutives()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(execs))
	for i, v := range execs {
		names[i] = v.Name
	}
	return names, nil
}

func (h *Houjin) GetHoujinRepresentativeNames() ([]string, error) {
	r, err := h.body.GetHoujinRepresentatives()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(r))
	for i, v := range r {
		names[i] = v.Name
	}
	return names, nil
}

func Parse(input string) (*Houjin, error) {
	s := normalize_kanji(input)
	header, body, err := DivideToukiboContent(s)
	if err != nil {
		return nil, err
	}
	houjinHeader, err := ParseHeader(header)
	if err != nil {
		return nil, err
	}

	houjinBody, err := ParseBody(body)
	if err != nil {
		return nil, err
	}

	// Get HoujinKaku
	houjinKakuType := FindHoujinKaku(houjinHeader.CompanyName)
	if houjinKakuType == HoujinKakuUnknown {
		if strings.Contains(s, "宗教法人") {
			houjinBody.HoujinKaku = HoujinKakuShukyo
		}
	} else {
		houjinBody.HoujinKaku = houjinKakuType
	}

	return &Houjin{header: houjinHeader, body: houjinBody}, nil
}
