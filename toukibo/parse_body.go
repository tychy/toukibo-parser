package toukibo

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
)

const (
	revert1         = "┃　　　　　　　　├─────────────────────────────────────┨"
	revert2         = "┃　　　　　　　　├───────────────────────┬─────────────┨"
	revert3         = "┃　　　　　　　　├───────────────────────┼─────────────┨"
	revert4         = "┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├─────────────┨"
	revert5         = "┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　┌─────────────┨"
	splitExecutive1 = "┃　　　　　　　　├───────────────────────┼─────────────┨"
	splitExecutive2 = "┃　　　　　　　　├─────────────────────────────────────┨"
	splitExecutive3 = "┃　　　　　　　　├───────────────────────┴─────────────┨"
	splitExecutive4 = "┃　　　　　　　　│　（特定社員）　　　　　　　　　　　　　　　　　　　　　　　　　　　　　　┃ ┃　　　　　　　　├───────────────────────┬─────────────┨" // sample133を通すためのハック

	trimExecutive1 = "┃　　　　　　　　│　　　　　　　　　　　　　　　　　　　　　　　├－－－－－－－－－－－－－┨"

	positionsPattern = "代表取締役|取締役・監査等|取締役|監査役|会計監査人|代表理事|理事長|理事|監事|代表社員|" +
		"業務執行社員|会長|代表清算人|清算人|代表役員|会計参与|無限責任社員|有限責任社員|破産管財人|評議員|代表者|" +
		"会頭|学長|代表執行役|執行役|報酬委員|監査委員|指名委員|職務執行者|保全管財人"
)

// revertで先に分割、複数行にまたがっているのを結合、それぞれの単位で値、登記日を取得する
// [{value, valid, since}, ]
// sinceの""は最初の登記
func splitReverts(s string) []string {
	revertPattern := fmt.Sprintf("%s|%s|%s|%s|%s", revert1, revert2, revert3, revert4, revert5)
	return regexp.MustCompile(revertPattern).Split(s, -1)
}

func splitExecutives(s string) []string {
	pattern := fmt.Sprintf("%s|%s|%s|%s", splitExecutive1, splitExecutive2, splitExecutive3, splitExecutive4)
	return regexp.MustCompile(pattern).Split(s, -1)
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
	pattern := fmt.Sprintf(`([%s]+)　*登記`, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		return trimAllSpace(matches[1]), nil
	}
	return "", fmt.Errorf("failed to get registerAt from %s", s)
}

func getResignedAt(s string) (string, error) {
	pattern := fmt.Sprintf(`([%s]+)　*(辞任|退任|死亡|抹消|廃止|解任|退社)`, ZenkakuStringPattern)
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
					slog.Debug(fmt.Sprintf("GetHoujinValue: failed to get registerAt from %s", parts[i]))
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
	pattern := fmt.Sprintf("社員　*([%s]+)", ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		s = regex.ReplaceAllString(s, "┃　　　　　　　　")
		return s, "社員", trimAllSpace(matches[1])
	}
	return s, "", ""
}

func normalizeExecutiveName(s string) string {
	name := trimAllSpace(s)
	// 金額の記載がある場合、役員名から削除
	name = trimPattern(name, fmt.Sprintf("金[%s]+(?:万円|円)全部履行", ZenkakuNumberPattern))
	name = trimPattern(name, fmt.Sprintf("金[%s]+(?:万円|円)", ZenkakuNumberPattern)) // sample519用のハック

	// 生年月日の記載がある場合は削除
	// NOTE: 取締役に同姓同名の別人がいる場合は、例外的に生年月日が登記される
	name = trimPattern(name, "(大正|昭和|平成|令和)[0-9０-９]+年[0-9０-９]+月[0-9０-９]+日生")
	return name
}

func getMultipleExecutiveNamesAndPositions(s string) (result []struct{ Name, Position string }, three []string) {
	var onNameAndPos bool
	for _, l := range extractLines(s) {
		_, remain := getPartOne(l)
		b, c := getPartTwo(remain)
		three = append(three, c)

		if regexp.MustCompile("^　*$").MatchString(b) {
			// 空行挟んで名前や役職が続く場合があるため、何もせずにスキップ
			continue
		}

		if onNameAndPos {
			// 前の行に役職 + 名前が記述されていた場合

			// NOTE: 役職が「取締役・監査等委員」の場合は、役職部分が2行にまたがる
			if result[len(result)-1].Position == "取締役・監査等" && strings.HasPrefix(b, "　委員") {
				result[len(result)-1].Name += trimLeadingTrailingSpace(strings.TrimPrefix(b, "　委員"))
				result[len(result)-1].Position += "委員"
				continue
			}

			// 字下げされている場合は、名前の続きとして扱う
			if strings.HasPrefix(b, "　　") {
				result[len(result)-1].Name += trimLeadingTrailingSpace(b)
				continue
			}

			onNameAndPos = false
		}
		if match := regexp.MustCompile(fmt.Sprintf("(%s)　+([%s]+)", positionsPattern, ZenkakuStringPattern)).FindStringSubmatch(b); len(match) == 3 {
			// NOTE: 名前や役職が複数行にまたがる可能性があり、Name と Position は次行以降の内容も踏まえて確定させる必要がある
			onNameAndPos = true
			result = append(result, struct{ Name, Position string }{Name: trimAllSpace(match[2]), Position: trimAllSpace(match[1])})
			continue
		}

		if _, pos, name := getShain(b); pos != "" {
			result = append(result, struct{ Name, Position string }{Name: name, Position: pos})
			continue
		}
	}

	for i := 0; i < len(result); i++ {
		result[i].Name = normalizeExecutiveName(result[i].Name)
	}

	return
}

// ┃ * ┃ or ┃ * ┨ の中身を抽出
func extractLines(s string) []string {
	var res []string
	cur := ""
	for _, r := range s {
		if r == '┃' || r == '┨' {
			if len(cur) > 0 {
				res = append(res, cur)
			}
			cur = ""
			continue
		}
		cur += string(r)
	}
	return res
}

func getPartOne(s string) (string, string) {
	cur := ""
	for _, r := range s {
		if r == '│' {
			return cur, s[len(cur)+3:] // 3は仕切りの大きさ
		}
		cur += string(r)
	}
	return cur, ""
}

func getPartTwo(s string) (string, string) {
	cur := ""
	for _, r := range s {
		if r == '│' || r == '├' {
			return cur, s[len(cur)+3:] // 3は仕切りの大きさ
		}
		cur += string(r)
	}
	return cur, ""
}

func splitThree(s string) (string, string, string) {
	partOne, remain := getPartOne(s)
	partTwo, partThree := getPartTwo(remain)
	return trimLeadingTrailingSpace(partOne), trimLeadingTrailingSpace(partTwo), trimLeadingTrailingSpace(partThree)
}

func GetHoujinExecutiveValue(s string) (HoujinExecutiveValueArray, error) {
	if debug {
		PrintBar()
	}

	s = trimPattern(s, trimExecutive1) // 必要のない仕切りを削除
	parts := splitReverts(s)
	evsArr := make(HoujinExecutiveValueArray, 0, len(parts))

	var idx int
	for _, p := range parts {
		if debug {
			PrintSlice(extractLines(p))
		}

		var evs []HoujinExecutiveValue

		// 役員名、役職を取得
		posAndNames, three := getMultipleExecutiveNamesAndPositions(p)
		for _, posAndName := range posAndNames {
			ev := HoujinExecutiveValue{
				IsValid:  true,
				Name:     posAndName.Name,
				Position: posAndName.Position,
			}
			evs = append(evs, ev)
		}

		// 登記日、辞任日を取得
		var registerAt, resignedAt string
		for _, t := range three {
			if at, _ := getRegisterAt(t); at != "" {
				registerAt = at
			}
			if at, _ := getResignedAt(t); at != "" {
				resignedAt = at
			}
		}
		if registerAt != "" {
			for i := range evs {
				evs[i].RegisterAt = registerAt
			}
		}
		if resignedAt != "" {
			for i := range evs {
				evs[i].ResignedAt = resignedAt
				evs[i].IsValid = false
			}
		}

		if idx > 0 {
			if len(evs) == 0 {
				if registerAt != "" {
					evsArr[idx-1].RegisterAt = registerAt
					evsArr[idx-1].ResignedAt = resignedAt
					evsArr[idx-1].IsValid = false
				}
				continue
			}

			// 簡単のために、一つの欄に複数の役員情報が記載されている場合は無効処理等をスキップしている
			if len(evs) == 1 {
				// 同じ氏名、役職の役員が連続している場合、前の役員を無効にする
				if evsArr[idx-1].Name == evs[0].Name && evsArr[idx-1].Position == evs[0].Position {
					evsArr[idx-1].IsValid = false
				}

				// sample 30, 89, 106用のハック
				// XXXXの氏/名称変更がある場合、その前の役員は無効にする
				if strings.Contains(strings.Join(three, ""), evsArr[idx-1].Name+"の氏変更") ||
					strings.Contains(strings.Join(three, ""), evsArr[idx-1].Name+"の氏名変更") ||
					strings.Contains(strings.Join(three, ""), evsArr[idx-1].Name+"の名称変更") ||
					strings.Contains(strings.Join(three, ""), evsArr[idx-1].Name+"の名") {
					evsArr[idx-1].IsValid = false
				}
			}
		}

		idx += len(evs)
		evsArr = append(evsArr, evs...)
	}
	if debug {
		fmt.Println(evsArr)
	}
	return evsArr, nil

}

func (h *HoujinBody) ConsumeHoujinNumber(s string) bool {
	// 正規表現パターン: 全角数字で構成された法人番号
	pattern := "([０-９]{4}－[０-９]{2}－[０-９]{6})"
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		h.HoujinNumber = ZenkakuToHankaku(matches[1])
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
	return strings.HasPrefix(s, " ┃資本金の額") || strings.HasPrefix(s, " ┃特定資本金の額") ||
		strings.Contains(s, " ┃資本金") ||
		strings.Contains(s, "払込済出資総額") || strings.Contains(s, "出資の総額") ||
		strings.Contains(s, "資産の総額") || strings.Contains(s, "基本財産の総額") || strings.Contains(s, "特定資本の額") ||
		strings.Contains(s, "払い込んだ出資の") // sample 520用のハック
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
		h.HoujinCreatedAt = ZenkakuToHankaku(strings.TrimSpace(matches[2]))
		return true
	}
	return false
}

func (h *HoujinBody) ConsumeHoujinBankruptedAt(s string) bool {
	pattern := fmt.Sprintf("┃破　産　*│　*([%s]+日)([%s]*)", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		h.HoujinBankruptedAt = ZenkakuToHankaku(strings.TrimSpace(matches[1]))
		return true
	}
	return false
}

func (h *HoujinBody) ConsumeHoujinDissolvedAt(s string) bool {
	if !strings.Contains(s, "┃解　散") {
		return false
	}

	var s2 string
	for _, p := range extractLines(s) {
		_, b, _ := splitThree(p)
		s2 += b
	}

	// ex 北海道知事の命令により解散
	// ex 会社法４７２条第１項の規定により解散
	pattern := fmt.Sprintf("([%s]+日)([%s]*)により解散", ZenkakuStringPattern, ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s2)
	if len(matches) > 0 {
		h.HoujinDissolvedAt = ZenkakuToHankaku(strings.TrimSpace(matches[1]))
		return true
	}
	return false
}

func (h *HoujinBody) ConsumeHoujinContinuedAt(s string) bool {
	// ex 令和2年7月1日会社継続
	pattern := fmt.Sprintf("┃会社継続　*│　*([%s]+日)会社継続", ZenkakuStringPattern)
	regex := regexp.MustCompile(pattern)

	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		h.HoujinContinuedAt = ZenkakuToHankaku(strings.TrimSpace(matches[1]))
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
		strings.Contains(s, "解散の事由") || strings.Contains(s, "監査役会設置会社") ||
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
	if h.ConsumeHoujinBankruptedAt(s) {
		return nil
	}
	if h.ConsumeHoujinDissolvedAt(s) {
		return nil
	}
	if h.ConsumeHoujinContinuedAt(s) {
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
			if strings.Contains(e, "監査役の監査の範囲を会計に関するものに限定") {
				continue
			}

			v, err := GetHoujinExecutiveValue(e)
			if err != nil {
				return err
			}
			h.HoujinExecutive[i] = v
		}
		return nil
	}
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
