package toukibo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type HoujinBody struct {
	HoujinNumber       string
	HoujinKaku         HoujinkakuType
	HoujinName         HoujinValueArray
	HoujinAddress      HoujinValueArray
	HoujinPurpose      HoujinValueArray
	HoujinKoukoku      string
	HoujinCreatedAt    string
	HoujinBankruptedAt string
	HoujinDissolvedAt  string
	HoujinContinuedAt  string
	HoujinCapital      HoujinValueArray
	HoujinStock        HoujinValueArray
	HoujinToukiRecord  HoujinValueArray
	HoujinExecutive    []HoujinExecutiveValueArray
}

func (h *HoujinBody) String() string {
	out := fmt.Sprintf("Body\n法人番号 : %s\n法人名  : %s\n法人住所 : %s\n目的   : %s\n公告   : %s\n成立年月日: %s\n解散年月日: %s\n資本金  : %s\n登記記録 : %s\n",
		h.HoujinNumber,
		h.HoujinName,
		h.HoujinAddress,
		h.HoujinPurpose,
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

func (h *HoujinBody) GetHoujinKaku() HoujinkakuType {
	if h.HoujinKaku == HoujinKakuUnknown {
		return HoujinKakuUnknown
	}
	return h.HoujinKaku
}

var warekiDatePattern = regexp.MustCompile(`^(明治|大正|昭和|平成|令和)(元|[0-9]+)年([0-9]+)月([0-9]+)日$`)

func warekiDateValue(s string) (int, bool) {
	normalized := strings.Join(strings.Fields(ZenkakuToHankaku(s)), "")
	matches := warekiDatePattern.FindStringSubmatch(normalized)
	if len(matches) != 5 {
		return 0, false
	}
	eraStart := map[string]int{
		"明治": 1868,
		"大正": 1912,
		"昭和": 1926,
		"平成": 1989,
		"令和": 2019,
	}[matches[1]]
	year := 1
	var err error
	if matches[2] != "元" {
		year, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, false
		}
	}
	month, err := strconv.Atoi(matches[3])
	if err != nil {
		return 0, false
	}
	day, err := strconv.Atoi(matches[4])
	if err != nil {
		return 0, false
	}
	return (eraStart+year-1)*10000 + month*100 + day, true
}

func (h *HoujinBody) isCurrentlyDissolved() bool {
	if h.HoujinDissolvedAt == "" {
		return false
	}
	if h.HoujinContinuedAt == "" {
		return true
	}
	dissolvedAt, dissolvedOK := warekiDateValue(h.HoujinDissolvedAt)
	continuedAt, continuedOK := warekiDateValue(h.HoujinContinuedAt)
	return dissolvedOK && continuedOK && dissolvedAt > continuedAt
}

func postProcessResponsibilityChangesGlobal(evsArr []HoujinExecutiveValue) {
	// 1. ResignedAtに「責任変更」が含まれている役員を無効化
	for i := range evsArr {
		if !evsArr[i].IsValid {
			continue
		}

		if strings.Contains(evsArr[i].ResignedAt, "責任変更") {
			evsArr[i].IsValid = false
			if DebugOn {
				fmt.Printf("Invalidating (resigned): %s %s\n", evsArr[i].Name, evsArr[i].Position)
			}
		}
	}

	// 2. 無限責任社員と有限責任社員の同一人物がいる場合、有限責任社員を無効化
	//    （有限→無限への責任変更を想定）
	for i := range evsArr {
		if !evsArr[i].IsValid {
			continue
		}

		if evsArr[i].Position == "無限責任社員" {
			if DebugOn {
				fmt.Printf("Found 無限責任社員: %s\n", evsArr[i].Name)
			}
			// 同じ名前の有限責任社員を探して無効化
			for j := range evsArr {
				if !evsArr[j].IsValid || i == j {
					continue
				}

				if evsArr[j].Name == evsArr[i].Name && evsArr[j].Position == "有限責任社員" {
					if DebugOn {
						fmt.Printf("Invalidating 有限責任社員: %s (found matching 無限責任社員)\n", evsArr[j].Name)
					}
					evsArr[j].IsValid = false
				}
			}
		}
	}
}

func (h *HoujinBody) GetHoujinExecutives() ([]HoujinExecutiveValue, error) {
	if len(h.HoujinExecutive) == 0 {
		if h.isCurrentlyDissolved() {
			// 法人が解散していれば役員はいなくても良い
			return []HoujinExecutiveValue{}, nil
		}
		return []HoujinExecutiveValue{}, fmt.Errorf("not found executives")
	}

	// まず全ての役員を集める
	var all []HoujinExecutiveValue
	for _, e := range h.HoujinExecutive {
		all = append(all, e...)
	}

	// 責任変更の後処理を全役員に対して実行
	postProcessResponsibilityChangesGlobal(all)

	// IsValidな役員のみを返す
	var res []HoujinExecutiveValue
	for _, v := range all {
		if v.IsValid {
			res = append(res, v)
		}
	}

	if len(res) > 0 {
		return res, nil
	}
	if h.isCurrentlyDissolved() {
		return []HoujinExecutiveValue{}, nil
	}
	return []HoujinExecutiveValue{}, fmt.Errorf("not found executives")
}

func (h *HoujinBody) GetHoujinRepresentatives() ([]HoujinExecutiveValue, error) {
	if len(h.HoujinExecutive) == 0 {
		if h.isCurrentlyDissolved() {
			// 法人が解散していれば代表はいなくても良い
			return []HoujinExecutiveValue{}, nil
		}
		return []HoujinExecutiveValue{}, fmt.Errorf("not found representative")
	}

	// 代表清算人が代表となる場合
	res := h.FindExecutivesByPosition("代表清算人")
	if len(res) > 0 {
		return res, nil
	}

	// 清算人が代表となる場合
	res = h.FindExecutivesByPosition("清算人")
	if len(res) > 0 {
		return res, nil
	}

	// 破産管財人が代表となる場合
	res = h.FindExecutivesByPosition("破産管財人")
	if len(res) > 0 {
		return res, nil
	}
	// 保全管財人が代表となる場合
	// 今のところこのケースは見つけていないが、sample1047を見て必要だと判断
	res = h.FindExecutivesByPosition("保全管財人")
	if len(res) > 0 {
		return res, nil
	}

	// 解散後は取締役等を代表者として扱わない。清算人等が登記されて
	// いない場合も、現在の代表者なしとして正常に返す。
	if h.isCurrentlyDissolved() {
		return []HoujinExecutiveValue{}, nil
	}

	res = h.FindExecutivesByPosition("代表取締役", "代表理事", "代表社員", "会長",
		"代表役員", "代表者", "理事長", "会頭", "学長", "代表執行役")
	if len(res) > 0 {
		return res, nil
	}

	houjinKaku := h.GetHoujinKaku()

	// 特定目的会社、有限会社は取締役が代表となる
	if houjinKaku == HoujinKakuYugen || houjinKaku == HoujinKakuTokuteiMokuteki {
		res = h.FindExecutivesByPosition("取締役")
		if len(res) > 0 {
			return res, nil
		}
	}

	if houjinKaku == HoujinKakuGousi {
		res = h.FindExecutivesByPosition("無限責任社員")
		if len(res) > 0 {
			return res, nil
		}
	}

	// 理事が代表となる場合
	res = h.FindExecutivesByPosition("理事")
	if len(res) > 0 {
		return res, nil
	}
	// 監査役が代表となる場合
	res = h.FindExecutivesByPosition("監査役")
	if len(res) > 0 {
		return res, nil
	}
	// 社員が代表となる場合
	res = h.FindExecutivesByPosition("社員")
	if len(res) > 0 {
		return res, nil
	}
	return []HoujinExecutiveValue{}, fmt.Errorf("not found representative")
}

// FindExecutivesByPosition は指定された役職の役員を検索する共通関数
func (h *HoujinBody) FindExecutivesByPosition(positions ...string) []HoujinExecutiveValue {
	var res []HoujinExecutiveValue
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if v.IsValid && contains(positions, v.Position) {
				res = append(res, v)
			}
		}
	}
	return res
}

// contains はスライスに要素が含まれているかチェックする
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
