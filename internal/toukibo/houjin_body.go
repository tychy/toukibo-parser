package toukibo

import "fmt"

type HoujinBody struct {
	HoujinNumber       string
	HoujinKaku         HoujinkakuType
	HoujinName         HoujinValueArray
	HoujinAddress      HoujinValueArray
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

func (h *HoujinBody) GetHoujinKaku() HoujinkakuType {
	if h.HoujinKaku == HoujinKakuUnknown {
		return HoujinKakuUnknown
	}
	return h.HoujinKaku
}

func (h *HoujinBody) GetHoujinExecutives() ([]HoujinExecutiveValue, error) {
	if len(h.HoujinExecutive) == 0 {
		if h.HoujinDissolvedAt != "" {
			// 法人が解散していれば代表はいなくても良い
			return []HoujinExecutiveValue{}, nil
		}
		return []HoujinExecutiveValue{}, fmt.Errorf("not found executives")
	}

	var res []HoujinExecutiveValue
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if v.IsValid {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		return res, nil
	}
	return []HoujinExecutiveValue{}, fmt.Errorf("not found executives")
}

func (h *HoujinBody) GetHoujinRepresentatives() ([]HoujinExecutiveValue, error) {
	if len(h.HoujinExecutive) == 0 {
		if h.HoujinDissolvedAt != "" {
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
