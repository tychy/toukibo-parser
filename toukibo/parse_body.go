package toukibo

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	revert1         = "┃　　　　　　　　├─────────────────────────────────────┨"
	revert2         = "┃　　　　　　　　├───────────────────────┬─────────────┨"
	revert3         = "┃　　　　　　　　├───────────────────────┼─────────────┨"
	revert4         = "┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├─────────────┨"
	splitExecutive1 = "┃　　　　　　　　├───────────────────────┼─────────────┨"
)

type HoujinValue struct {
	Value      string
	IsValid    bool
	RegisterAt string
}

func (h *HoujinValue) String() string {
	return fmt.Sprintf("value: %s, isValid: %v, registerAt: %s",
		h.Value, h.IsValid, h.RegisterAt)
}

type HoujinValueArray []HoujinValue

func (hva HoujinValueArray) String() string {
	var b strings.Builder
	for _, hv := range hva {
		b.WriteString("{")
		b.WriteString(hv.String())
		b.WriteString("},")
	}
	return b.String()
}

type HoujinExecutiveValue struct {
	Name       string
	Pisition   string
	Address    string
	IsValid    bool
	RegisterAt string
	ResignedAt string
}

func (h *HoujinExecutiveValue) String() string {
	return fmt.Sprintf("name: %s, position: %s, address: %s, isValid: %v, registerAt: %s, resignedAt: %s",
		h.Name, h.Pisition, h.Address, h.IsValid, h.RegisterAt, h.ResignedAt)
}

type HoujinExecutiveValueArray []HoujinExecutiveValue

func (hva HoujinExecutiveValueArray) String() string {
	var b strings.Builder
	for _, hv := range hva {
		b.WriteString("{")
		b.WriteString(hv.String())
		b.WriteString("},")
	}
	return b.String()
}

type HoujinBody struct {
	HoujinNumber      string
	HoujinName        HoujinValueArray
	HoujinAddress     HoujinValueArray
	HoujinKoukoku     string
	HoujinCreatedAt   string
	HoujinCapital     HoujinValueArray
	HoujinToukiRecord HoujinValueArray
	HoujinExecutive   []HoujinExecutiveValueArray
}

func (h *HoujinBody) String() string {
	out := fmt.Sprintf("Body\n法人番号 : %s\n法人名  : %s\n法人住所 : %s\n公告   : %s\n成立年月日: %s\n資本金  : %s\n登記記録 : %s\n",
		h.HoujinNumber,
		h.HoujinName,
		h.HoujinAddress,
		h.HoujinKoukoku,
		h.HoujinCreatedAt,
		h.HoujinCapital,
		h.HoujinToukiRecord,
	)
	fmt.Println("役員   : ")
	for _, e := range h.HoujinExecutive {
		out += "[" + e.String() + "],\n"
	}
	return out
}

func (h *HoujinBody) GetHoujinRepresentative() (HoujinExecutiveValue, error) {
	if len(h.HoujinExecutive) == 0 {
		return HoujinExecutiveValue{}, fmt.Errorf("not found representative")
	}

	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Pisition == "代表取締役" || v.Pisition == "代表理事") && v.IsValid {
				return v, nil
			}
		}
	}
	return HoujinExecutiveValue{}, fmt.Errorf("not found representative")
}

// revertで先に分割、複数行にまたがっているのを結合、それぞれの単位で値、登記日を取得する
// [{value, valid, since}, ]
// sinceの""は最初の登記
func splitReverts(s string) []string {
	revertPattern := fmt.Sprintf("%s|%s|%s|%s", revert1, revert2, revert3, revert4)
	re := regexp.MustCompile(revertPattern)
	parts := re.Split(s, -1)
	return parts
}

func splitExecutives(s string) []string {
	pattern := splitExecutive1
	re := regexp.MustCompile(pattern)
	parts := re.Split(s, -1)
	return parts
}

func trimPattern(s, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(s, "")
}

func mergeLines(s string) string {
	pattern := "　*┃ ┃　*│　*"
	return trimPattern(s, pattern)
}

func trimChangeAndRegisterAt(s string) (string, string, string) {
	// trim ┃　　　　　　　　│　　　　　　　　　平成３０年　７月３１日変更　　平成３０年　８月２７日登記┃
	pattern := fmt.Sprintf("┃　*│　*([%s]+)変更　*([%s]+)登記　*┃", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimPattern(s, pattern), trimAllSpace(matches[1]), trimAllSpace(matches[2])
	}
	return s, "", ""
}

func trimRegisterAt(s string) (string, string) {
	// trim ┃事項　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　　平成２０年　７月２５日登記┃
	pattern := fmt.Sprintf("┃[%s]+　*│　*([%s]+)(登記|移記)　*┃", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimPattern(s, pattern), trimAllSpace(matches[1])
	}
	return s, ""
}

func getValue(s string) (string, error) {
	pattern := fmt.Sprintf(`([%s]+)　*│　*([%s]+)`, ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[2]), nil
	}
	return "", fmt.Errorf("failed to get value from %s", s)
}

func getRegisterAt(s string) (string, error) {
	pattern := fmt.Sprintf(`│([%s]+)　*登記`, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimAllSpace(matches[1]), nil
	}
	return "", fmt.Errorf("failed to get registerAt from %s", s)
}

func getResignedAt(s string) (string, error) {
	pattern := fmt.Sprintf(`│([%s]+)　*辞任`, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimAllSpace(matches[1]), nil
	}
	return "", fmt.Errorf("failed to get resignedAt from %s", s)
}

func GetHoujinValue(s string) (HoujinValueArray, error) {
	parts := splitReverts(s)
	res := make(HoujinValueArray, len(parts))
	for i, s := range parts {
		var registerAt string
		s, _, registerAt = trimChangeAndRegisterAt(s)
		if registerAt == "" {
			s, registerAt = trimRegisterAt(s)
		}
		s = mergeLines(s)

		isLast := i == len(parts)-1
		value, err := getValue(s)
		if err != nil {
			return nil, err
		}

		if i == 0 {
			res[i] = HoujinValue{
				Value:      value,
				IsValid:    isLast,
				RegisterAt: registerAt,
			}
		} else {
			if registerAt == "" {
				registerAt, err = getRegisterAt(s)
				if err != nil {
					return nil, err
				}
			}
			res[i] = HoujinValue{
				Value:      value,
				IsValid:    isLast,
				RegisterAt: registerAt,
			}
		}
	}
	return res, nil
}

// trim all space
func trimAllSpace(s string) string {
	return strings.ReplaceAll(s, "　", "")
}

func getExecutiveNameAndPosition(s string) (string, string, string) {
	pattern := fmt.Sprintf("(代表取締役|取締役|監査役|会計監査人|代表理事)　*([%s]+)", ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimPattern(s, pattern), trimAllSpace(matches[1]), trimAllSpace(matches[2])
	}
	return s, "", ""
}

func GetHoujinExecutiveValue(s string) (HoujinExecutiveValueArray, error) {
	parts := splitReverts(s)
	res := make(HoujinExecutiveValueArray, 0, len(parts))
	for i, s := range parts {
		isLast := i == len(parts)-1

		s, position, name := getExecutiveNameAndPosition(s)
		if position == "" || name == "" {
			if strings.Contains(s, "辞任") {
				if !isLast || i == 0 {
					return nil, fmt.Errorf("resign is not the last %s", s)
				}
				resignedAt, err := getResignedAt(s)
				if err != nil {
					return nil, err
				}
				res[i-1].ResignedAt = resignedAt
				break
			}
			return nil, fmt.Errorf("failed to get executive name and position from %s", s)
		}
		s = mergeLines(s)

		value, err := getValue(s)
		if err != nil {
			return nil, err
		}

		registerAt, err := getRegisterAt(s)
		if err != nil {
			return nil, err
		}

		res = append(res, HoujinExecutiveValue{
			Name:       name,
			Address:    value,
			Pisition:   position,
			IsValid:    isLast,
			RegisterAt: registerAt,
		})
	}
	return res, nil
}

func (h *HoujinBody) ConsumeHoujinNumber(s string) bool {
	// 正規表現パターン: 全角数字で構成された法人番号
	pattern := "([０-９]{1,4}－[０-９]{1,2}－[０-９]{1,6})"
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		h.HoujinNumber = zenkakuToHankaku(matches[1])
		return true
	}
	return false
}

func (h *HoujinBody) ConsumeHoujinName(s string) bool {
	return strings.Contains(s, "商　号") || strings.Contains(s, "名　称")
}

func (h *HoujinBody) ConsumeHoujinAddress(s string) bool {
	return strings.Contains(s, "本　店") || strings.Contains(s, "主たる事務所")
}

func (h *HoujinBody) ConsumeHoujinKoukoku(s string) bool {
	pattern := `(公告をする方法|公告の方法)　*│　*(.+)┃`
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		v, err := GetHoujinValue(s)
		if err != nil || len(v) == 0 {
			return false
		}
		h.HoujinKoukoku = v[0].Value
		return true
	}
	return false
}

func (h *HoujinBody) ConsumeHoujinCapital(s string) bool {
	return strings.Contains(s, "資本金の額") || strings.Contains(s, "払込済出資総額")
}

func (h *HoujinBody) ConsumeHoujinToukiRecord(s string) bool {
	return strings.Contains(s, "登記記録に関する")
}

func (h *HoujinBody) ConsumeHoujinExecutive(s string) bool {
	return strings.Contains(s, "役員に関する事項")
}

func (h *HoujinBody) ConsumeHoujinCreatedAt(s string) bool {
	pattern := fmt.Sprintf(`(会社|法人)成立の年月日　*│　*([%s　]+)`, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		h.HoujinCreatedAt = strings.TrimSpace(matches[2])
		return true
	}
	return false
}

func (h *HoujinBody) ParseBodyMain(s string) error {
	if strings.Contains(s, "発行可能株式総数") || strings.Contains(s, "┃目　的") || strings.Contains(s, "┃目的等") ||
		strings.Contains(s, "出資１口の金額") || strings.Contains(s, "出資の総口数") || strings.Contains(s, "出資払込の方法") ||
		strings.Contains(s, "発行済株式の総数") || strings.Contains(s, "株式の譲渡制限") || strings.Contains(s, "株券を発行する旨") ||
		strings.Contains(s, "取締役等の会社") || strings.Contains(s, "非業務執行取締役") ||
		strings.Contains(s, "取締役会設置会社") || strings.Contains(s, "監査役設置会社") || strings.Contains(s, "会計監査人設置会") ||
		strings.Contains(s, "地区") {
		// skip
		return nil
	}

	if h.ConsumeHoujinNumber(s) {
		return nil
	}
	if h.ConsumeHoujinName(s) {
		v, err := GetHoujinValue(s)
		if err != nil {
			return err
		}

		h.HoujinName = v
		return nil
	}
	if h.ConsumeHoujinAddress(s) {
		v, err := GetHoujinValue(s)
		if err != nil {
			return err
		}
		h.HoujinAddress = v
		return nil
	}

	if h.ConsumeHoujinKoukoku(s) {
		return nil
	}
	if h.ConsumeHoujinCreatedAt(s) {
		return nil
	}
	if h.ConsumeHoujinCapital(s) {
		v, err := GetHoujinValue(s)
		if err != nil {
			return err
		}
		h.HoujinCapital = v
		return nil
	}
	if h.ConsumeHoujinToukiRecord(s) {
		v, err := GetHoujinValue(s)
		if err != nil {
			return err
		}
		h.HoujinToukiRecord = v
		return nil
	}
	if h.ConsumeHoujinExecutive(s) {
		executives := splitExecutives(s)
		h.HoujinExecutive = make([]HoujinExecutiveValueArray, len(executives))
		for i, e := range executives {
			v, err := GetHoujinExecutiveValue(e)
			if err != nil {
				return err
			}
			h.HoujinExecutive[i] = v
		}
		return nil
	}

	fmt.Println("not consumed: ", s)
	return nil
}

func ParseBody(pairs []string) (*HoujinBody, error) {
	h := &HoujinBody{}
	for i := 0; i < len(pairs); i++ {
		err := h.ParseBodyMain(pairs[i])
		if err != nil {
			return nil, err
		}
	}
	return h, nil
}
