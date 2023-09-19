package toukibo

import (
	"fmt"
	"regexp"
	"strings"
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

	separatorPattern := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s", separator1, separator2, separator3, separator4, separator5, separator6, separator7)
	re := regexp.MustCompile(separatorPattern)

	parts := re.Split(content, -1)
	return header, parts, nil
}

type Houjin struct {
	header *HoujinHeader
	body   *HoujinBody
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
	header, body, err := DivideToukiboContent(input)
	if err != nil {
		return nil, err
	}
	houjinHeader, err := ParseHeader(header)
	if err != nil {
		return nil, err
	}

	fmt.Println(houjinHeader.String())

	houjinBody, err := ParseBody(body)
	if err != nil {
		return nil, err
	}

	fmt.Println(houjinBody.String())
	return &Houjin{header: houjinHeader, body: houjinBody}, nil
}
