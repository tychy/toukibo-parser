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
	splitExecutive2 = "┃　　　　　　　　├─────────────────────────────────────┨"
	splitExecutive3 = "┃　　　　　　　　├───────────────────────┴─────────────┨"
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
	Position   string
	Address    string
	IsValid    bool
	RegisterAt string
	ResignedAt string
}

func (h *HoujinExecutiveValue) String() string {
	return fmt.Sprintf("name: %s, position: %s, address: %s, isValid: %v, registerAt: %s, resignedAt: %s",
		h.Name, h.Position, h.Address, h.IsValid, h.RegisterAt, h.ResignedAt)
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
	HoujinDissolvedAt string
	HoujinCapital     HoujinValueArray
	HoujinToukiRecord HoujinValueArray
	HoujinExecutive   []HoujinExecutiveValueArray
}

func (h *HoujinBody) String() string {
	out := fmt.Sprintf("Body\n法人番号 : %s\n法人名  : %s\n法人住所 : %s\n公告   : %s\n成立年月日: %s\n解散年月日: %s\n資本金  : %s\n登記記録 : %s\n",
		h.HoujinNumber,
		h.HoujinName,
		h.HoujinAddress,
		h.HoujinKoukoku,
		h.HoujinCreatedAt,
		h.HoujinDissolvedAt,
		h.HoujinCapital,
		h.HoujinToukiRecord,
	)
	out += "役員  : \n"
	for _, e := range h.HoujinExecutive {
		out += "[" + e.String() + "],\n"
	}
	return out
}

func (h *HoujinBody) GetHoujinKaku() (HoujinkakuType, error) {
	if len(h.HoujinName) == 0 {
		return HoujinKakuUnknown, fmt.Errorf("not found houjin name")
	}
	for _, v := range h.HoujinName {
		if v.IsValid {
			return FindHoujinKaku(v.Value), nil
		}
	}
	return HoujinKakuUnknown, fmt.Errorf("not found houjin name")
}

func (h *HoujinBody) GetHoujinRepresentatives() ([]HoujinExecutiveValue, error) {
	if len(h.HoujinExecutive) == 0 {
		if h.HoujinDissolvedAt != "" {
			// 法人が解散していれば代表はいなくても良い
			return []HoujinExecutiveValue{}, nil
		}
		return []HoujinExecutiveValue{}, fmt.Errorf("not found representative")
	}

	var res []HoujinExecutiveValue
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "代表取締役" || v.Position == "代表理事" || v.Position == "代表社員" || v.Position == "会長" || v.Position == "代表役員" || v.Position == "代表者" || v.Position == "会頭") && v.IsValid {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		return res, nil
	}

	houjinKaku, err := h.GetHoujinKaku()
	if err != nil {
		return nil, err
	}

	// 特定目的会社、有限会社は取締役が代表となる
	if houjinKaku == HoujinKakuYugen || houjinKaku == HoujinKakuTokuteiMokuteki {
		var res []HoujinExecutiveValue
		for _, e := range h.HoujinExecutive {
			for _, v := range e {
				if (v.Position == "取締役") && v.IsValid {
					res = append(res, v)
				}
			}
		}
		if len(res) > 0 {
			return res, nil
		}
	}

	if houjinKaku == HoujinKakuGousi {
		var res []HoujinExecutiveValue
		for _, e := range h.HoujinExecutive {
			for _, v := range e {
				if (v.Position == "無限責任社員") && v.IsValid {
					res = append(res, v)
				}
			}
		}
		if len(res) > 0 {
			return res, nil
		}
	}

	// 理事が代表となる場合
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "理事") && v.IsValid {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	// 清算人が代表となる場合
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "清算人") && v.IsValid {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	// 監査役が代表となる場合
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "監査役") && v.IsValid {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	// 社員が代表となる場合
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "社員") && v.IsValid {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	return []HoujinExecutiveValue{}, fmt.Errorf("not found representative")
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
	pattern := fmt.Sprintf("%s|%s|%s", splitExecutive1, splitExecutive2, splitExecutive3)
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
	pattern := fmt.Sprintf(`│([%s]+)　*(辞任|退任|死亡|抹消|廃止|解任)`, ZenkakuStringPattern)
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
					// 登記が記載されていない場合無視する
					fmt.Printf("failed to get registerAt from %s", parts[i])
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

func getShain(s string) (string, string, string) {
	pattern := fmt.Sprintf("┃　*│　*社員　*([%s]+)", ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		s = regex.ReplaceAllString(s, "┃　　　　　　　　")
		return s, "社員", trimAllSpace(matches[1])
	}
	return s, "", ""
}

func getExecutiveNameAndPosition(s string) (string, string, string) {
	positions := "代表取締役|取締役|監査役|会計監査人|代表理事|理事|監事|代表社員|業務執行社員|会長|清算人|代表役員|会計参与|無限責任社員|有限責任社員|破産管財人|評議員|代表者|会頭"
	pattern := fmt.Sprintf("(%s)　*([%s]+)", positions, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimPattern(s, pattern), trimAllSpace(matches[1]), trimAllSpace(matches[2])
	}

	s, pos, name := getShain(s)
	if pos != "" {
		return s, pos, name
	}

	return s, "", ""
}

func isResigned(s string) bool {
	return strings.Contains(s, "辞任") || strings.Contains(s, "退任") || strings.Contains(s, "死亡") ||
		strings.Contains(s, "抹消") || strings.Contains(s, "廃止") || strings.Contains(s, "解任")
}

func GetHoujinExecutiveValue(s string) (HoujinExecutiveValueArray, error) {
	parts := splitReverts(s)
	res := make(HoujinExecutiveValueArray, 0, len(parts))
	for i, s := range parts {
		isLast := i == len(parts)-1

		s, position, name := getExecutiveNameAndPosition(s)
		if position == "" || name == "" {
			if isResigned(s) {
				if !isLast || i == 0 {
					return nil, fmt.Errorf("resign is not the last %s", s)
				}
				resignedAt, err := getResignedAt(s)
				if err != nil {
					return nil, err
				}
				res[i-1].ResignedAt = resignedAt
				res[i-1].IsValid = false
				break
			}
			return nil, fmt.Errorf("failed to get executive name and position from %s", s)
		}
		s = mergeLines(s)

		address, err := getValue(s)
		if err != nil {
			return nil, err
		}

		var registerAt string
		if !(i == 0 && isLast) {
			registerAt, err = getRegisterAt(s)
			if err != nil {
				// 登記が記載されていない場合無視する
				fmt.Printf("failed to get registerAt from %s", parts[i])
			}
		}
		var resignedAt string
		if isResigned(s) {
			resignedAt, err = getResignedAt(s)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, HoujinExecutiveValue{
			Name:       name,
			Address:    address,
			Position:   position,
			IsValid:    isLast && resignedAt == "", // 退任などが記載されないまま無効になる場合をisLastで判定する
			RegisterAt: registerAt,
			ResignedAt: resignedAt,
		})

	}
	fmt.Println(res)
	return res, nil
}

func (h *HoujinBody) ConsumeHoujinNumber(s string) bool {
	// 正規表現パターン: 全角数字で構成された法人番号
	pattern := "([０-９]{4}－[０-９]{2}－[０-９]{6})"
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
	pattern := `(公告をする方法|公告の方法|法人の公告方法)　*│　*(.+)┃`
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
	return strings.Contains(s, "役員に関する事項") || strings.Contains(s, "社員に関する事項")
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

func (h *HoujinBody) ConsumeHoujinDissolvedAt(s string) bool {
	// ex 北海道知事の命令により解散
	// ex 会社法４７２条第１項の規定により解散
	pattern := fmt.Sprintf("┃解　散　*│　*([%s]+日)([%s]*)により解散", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		h.HoujinDissolvedAt = strings.TrimSpace(matches[1])
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
		strings.Contains(s, "地区") || strings.Contains(s, "解散の事由") || strings.Contains(s, "監査役会設置会社") || strings.Contains(s, "資産の総額") ||
		strings.Contains(s, "地　区") || strings.Contains(s, "支　店") || strings.Contains(s, "従たる事務所") {
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
	if h.ConsumeHoujinDissolvedAt(s) {
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

	// fmt.Println("not consumed: ", s)
	// TOOD enable this for debug
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
