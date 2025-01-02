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

func (h *Houjin) GetHoujinNameHistory() HoujinValueArray {
	return h.body.HoujinName
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

func (h *Houjin) GetHoujinContinuedAt() string {
	return h.body.HoujinContinuedAt
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

func (h *Houjin) GetHoujinStock() HoujinStock {
	for _, v := range h.body.HoujinStock {
		if v.IsValid {
			res := v.Value
			if len(res) < 1 {
				return HoujinStock{Total: 0}
			}
			return GetHoujinStock(v.Value)
		}
	}
	return HoujinStock{Total: 0}
}

func (h *Houjin) GetHoujinTotalStock() int {
	stock := h.GetHoujinStock()
	return stock.Total
}

func (h *Houjin) GetHoujinExecutives() (HoujinExecutiveValueArray, error) {
	return h.body.GetHoujinExecutives()
}

func (h *Houjin) GetHoujinExecutiveNames() ([]string, error) {
	execs, err := h.body.GetHoujinExecutives()
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
	s := normalizeKanji(input)
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

	houjinKakuType := FindHoujinKaku(houjinHeader.CompanyName, s)
	houjinBody.HoujinKaku = houjinKakuType

	return &Houjin{header: houjinHeader, body: houjinBody}, nil
}
