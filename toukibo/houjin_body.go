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

	var res []HoujinExecutiveValue
	// 代表清算人が代表となる場合
	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "代表清算人") && v.IsValid {
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

	for _, e := range h.HoujinExecutive {
		for _, v := range e {
			if (v.Position == "代表取締役" || v.Position == "代表理事" || v.Position == "代表社員" || v.Position == "会長" || v.Position == "代表役員" || v.Position == "代表者" || v.Position == "理事長" || v.Position == "会頭" || v.Position == "学長") && v.IsValid {
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
