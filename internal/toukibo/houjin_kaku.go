package toukibo

import (
	"strings"
)

type HoujinkakuType string

const (
	HoujinKakuUnknown         HoujinkakuType = "不明"
	HoujinKakuKabusiki        HoujinkakuType = "株式会社"
	HoujinKakuYugen           HoujinkakuType = "有限会社"
	HoujinKakuGoudou          HoujinkakuType = "合同会社"
	HoujinKakuGousi           HoujinkakuType = "合資会社"
	HoujinKakuGoumei          HoujinkakuType = "合名会社"
	HoujinKakuTokuteiMokuteki HoujinkakuType = "特定目的会社"
	HoujinKakuKyodou          HoujinkakuType = "協同組合"
	HoujinKakuRoudou          HoujinkakuType = "労働組合"
	HoujinKakuSinrin          HoujinkakuType = "森林組合"
	HoujinKakuSeikatuEisei    HoujinkakuType = "生活衛生同業組合"
	HoujinKakuSinyou          HoujinkakuType = "信用金庫"
	HoujinKakuShokoukai       HoujinkakuType = "商工会"
	HoujinKakuKoueki          HoujinkakuType = "公益財団法人"
	HoujinKakuNouji           HoujinkakuType = "農事組合"
	HoujinKakuShukyo          HoujinkakuType = "宗教法人"
	HoujinKakuKanriKumiai     HoujinkakuType = "管理組合法人"
	HoujinKakuIryo            HoujinkakuType = "医療法人"
	HoujinKakuSihoshosi       HoujinkakuType = "司法書士法人"
	HoujinKakuZeirishi        HoujinkakuType = "税理士法人"
	HoujinKakuShakaifukusi    HoujinkakuType = "社会福祉法人"
	HoujinKakuIppanShadan     HoujinkakuType = "一般社団法人"
	HoujinKakuKouekiShadan    HoujinkakuType = "公益社団法人"
	HoujinKakuIppanZaisan     HoujinkakuType = "一般財産法人"
	HoujinKakuIppanZaidan     HoujinkakuType = "一般財団法人"
	HoujinKakuNPO             HoujinkakuType = "NPO法人"
	HoujinKakuTokuteiHieiri   HoujinkakuType = "特定非営利活動法人"
	HoujinKakuPoliticalParty  HoujinkakuType = "政党"
	HoujinKakuUniversity      HoujinkakuType = "国立大学法人"
	HoujinKakuGakko           HoujinkakuType = "学校法人"
	HoujinKakuBengoshi        HoujinkakuType = "弁護士法人"
	HoujinKakuDokuritsu       HoujinkakuType = "独立行政法人"

	// 認可法人と特殊法人の違い
	// https://www.gyoukaku.go.jp/siryou/tokusyu/about.pdf
	// 一覧　https://www.hatosan.com/utc/jdt.html

	// 認可法人
	// 特別の法律によって限定数設置されるが、特殊法人と異なり、「特別の設立行為」によ
	// って強制設立されるものではなく、法律の枠内において民間等の関係者が任意設立し、
	// 主務大臣の認可を受けたもの
	HoujinKakuNinka HoujinkakuType = "認可法人"

	// 特殊法人
	// 特別の法律により特別の設立行為をもって強制設立すべきものとされる法人。
	HoujinKakuTokushu HoujinkakuType = "特殊法人"

	// 特別民間法人
	HoujinKakuTokubetsuMinkan HoujinkakuType = "特別民間法人"
)

var tokushuHoujin = map[string]struct{}{
	// 内閣府所管
	"沖縄振興開発金融公庫":    {},
	"沖縄科学技術大学院大学学園": {},
	// 復興庁所管
	"福島国際研究機構": {},
	// 総務省所管
	"日本電信電話株式会社":  {},
	"東日本電信電話株式会社": {},
	"西日本電信電話株式会社": {},
	"日本放送協会":      {},
	"日本郵政株式会社":    {},
	"日本郵便株式会社":    {},
	// 財務省所管
	"日本たばこ産業株式会社":          {},
	"株式会社日本政策金融公庫":         {},
	"株式会社日本政策投資銀行":         {},
	"輸出入・港湾関連情報処理センター株式会社": {},
	"株式会社国際協力銀行":           {},
	// 文部科学省所管
	"日本私立学校振興・共済事業団": {},
	"放送大学学園":         {},
	// 厚生労働省所管
	"日本年金機構": {},
	// 農林水産省所管
	"日本中央競馬会": {},
	// 経済産業省所管
	"日本アルコール産業株式会社": {},
	"株式会社商工組合中央金庫":  {},
	"株式会社日本貿易保険":    {},
	// 国土交通省所管
	"新関西国際空港株式会社":    {},
	"北海道旅客鉄道株式会社":    {},
	"四国旅客鉄道株式会社":     {},
	"日本貨物鉄道株式会社":     {},
	"東京地下鉄株式会社":      {},
	"成田国際空港株式会社":     {},
	"東日本高速道路株式会社":    {},
	"中日本高速道路株式会社":    {},
	"西日本高速道路株式会社":    {},
	"首都高速道路株式会社":     {},
	"阪神高速道路株式会社":     {},
	"本州四国連絡高速道路株式会社": {},
	// 環境省所管
	"中間貯蔵・環境安全事業株式会社": {},
}

func IsTokushuHoujin(name string) bool {
	// https://www.gyoukaku.go.jp/siryou/tokusyu/taisyou_houjin.pdf
	_, ok := tokushuHoujin[name]
	return ok
}

var ninkaHoujin = map[string]struct{}{
	"銀行等保有株式取得機構":          {},
	"原子力損害賠償・廃炉等支援機構":      {},
	"株式会社地域経済活性化支援機構":      {},
	"株式会社民間資金等活用事業推進機構":    {},
	"株式会社東日本大震災事業者再生支援機構":  {},
	"大阪湾広域臨海環境整備センター":      {},
	"株式会社海外通信・放送・郵便事業支援機構": {},
	"外国人技能実習機構":            {},
	"日本銀行":                 {},
	"預金保険機構":               {},
	"日本赤十字社":               {},
	"農水産業協同組合貯金保険機構":       {},
	"株式会社農林漁業成長産業化支援機構":    {},
	"電力広域的運営推進機関":          {},
	"使用済燃料再処理機構":           {},
	"株式会社産業革新投資機構":         {},
	"株式会社海外需要開拓支援機構":       {},
	"株式会社海外交通・都市開発事業支援機構":  {},
	"株式会社脱炭素化支援機構":         {},
}

var tokubetsuMinkanHoujin = map[string]struct{}{
	// 内閣府所管
	"日本消防検定協会":         {},
	"消防団員等公務災害補償等共済基金": {},
	"自動車安全運転センター":      {},
	"日本公認会計士協会":        {},
	// 総務省所管
	"危険物保安技術協会":  {},
	"日本行政書士会連合会": {},
	// 法務省所管
	"日本司法書士会連合会":    {},
	"日本土地家屋調査士会連合会": {},
	// 財務省所管
	"日本税理士会連合会": {},
	// 厚生労働省所管
	"社会保険診療報酬支払基金":     {},
	"中央労働災害防止協会":       {},
	"建設業労働災害防止協会":      {},
	"陸上貨物運送事業労働災害防止協会": {},
	"林業・木材製造業労働災害防止協会": {},
	"港湾貨物運送事業労働災害防止協会": {},
	"中央職業能力開発協会":       {},
	"企業年金連合会":          {},
	"石炭鉱業年金基金":         {},
	"全国社会保険労務士会連合会":    {},
	"全国健康保険協会":         {},
	// 農林水産省所管
	"農林中央金庫":      {},
	"全国漁業共済組合連合会": {},
	// 経済産業省所管
	"東京中小企業投資育成株式会社":  {},
	"名古屋中小企業投資育成株式会社": {},
	"大阪中小企業投資育成株式会社":  {},
	"高圧ガス保安協会":        {},
	"日本電気計器検定所":       {},
	"日本商工会議所":         {},
	"全国商工会連合会":        {},
	"日本弁理士会":          {},
	// 国土交通省所管
	"日本勤労者住宅協会":  {},
	"軽自動車検査協会":   {},
	"日本小型船舶検査機構": {},
	"日本水先人会連合会":  {},
}

func IsNinkaHoujin(name string) bool {
	_, ok := ninkaHoujin[name]
	return ok
}

func IsTokubetsuMinkanHoujin(name string) bool {
	_, ok := tokubetsuMinkanHoujin[name]
	return ok
}

var politicalParty = map[string]struct{}{
	"公明党":      {},
	"国民民主党":    {},
	"参政党":      {},
	"社会民主党":    {},
	"自由民主党":    {},
	"日本維新の会":   {},
	"日本保守党":    {},
	"日本共産党":    {},
	"みんなでつくる党": {},
	"立憲民主党":    {},
	"れいわ新選組":   {},
}

func IsPoliticalParty(name string) bool {
	_, ok := politicalParty[name]
	return ok
}

func FindHoujinKaku(name, s string) HoujinkakuType {
	if IsPoliticalParty(name) {
		return HoujinKakuPoliticalParty
	}

	if IsNinkaHoujin(name) {
		return HoujinKakuNinka
	}

	if IsTokushuHoujin(name) {
		return HoujinKakuTokushu
	}

	if IsTokubetsuMinkanHoujin(name) {
		return HoujinKakuTokubetsuMinkan
	}

	if strings.Contains(name, "株式会社") {
		return HoujinKakuKabusiki
	} else if strings.Contains(name, "有限会社") {
		return HoujinKakuYugen
	} else if strings.Contains(name, "合同会社") {
		return HoujinKakuGoudou
	} else if strings.Contains(name, "合資会社") {
		return HoujinKakuGousi
	} else if strings.Contains(name, "合名会社") {
		return HoujinKakuGoumei
	} else if strings.Contains(name, "特定目的会社") {
		return HoujinKakuTokuteiMokuteki
	} else if strings.Contains(name, "協同組合") {
		return HoujinKakuKyodou
	} else if strings.Contains(name, "労働組合") {
		return HoujinKakuRoudou
	} else if strings.Contains(name, "森林組合") {
		return HoujinKakuSinrin
	} else if strings.Contains(name, "生活衛生同業組合") {
		return HoujinKakuSeikatuEisei
	} else if strings.Contains(name, "信用金庫") {
		return HoujinKakuSinyou
	} else if strings.Contains(name, "商工会") {
		return HoujinKakuShokoukai
	} else if strings.Contains(name, "公益財団法人") {
		return HoujinKakuKoueki
	} else if strings.Contains(name, "農事組合") {
		return HoujinKakuNouji
	} else if strings.Contains(name, "宗教法人") {
		// このパターンは今のところ存在しない 2024/08/04
		return HoujinKakuShukyo
	} else if strings.Contains(name, "管理組合法人") {
		return HoujinKakuKanriKumiai
	} else if strings.Contains(name, "医療法人") {
		return HoujinKakuIryo
	} else if strings.Contains(name, "司法書士法人") {
		return HoujinKakuSihoshosi
	} else if strings.Contains(name, "税理士法人") {
		return HoujinKakuZeirishi
	} else if strings.Contains(name, "社会福祉法人") {
		return HoujinKakuShakaifukusi
	} else if strings.Contains(name, "一般社団法人") {
		return HoujinKakuIppanShadan
	} else if strings.Contains(name, "公益社団法人") {
		return HoujinKakuKouekiShadan
	} else if strings.Contains(name, "一般財産法人") {
		return HoujinKakuIppanZaisan
	} else if strings.Contains(name, "一般財団法人") {
		return HoujinKakuIppanZaidan
	} else if strings.Contains(name, "NPO法人") || strings.Contains(name, "ＮＰＯ法人") {
		return HoujinKakuNPO
	} else if strings.Contains(name, "特定非営利活動法人") {
		return HoujinKakuTokuteiHieiri
	} else if strings.Contains(name, "国立大学法人") {
		return HoujinKakuUniversity
	} else if strings.Contains(name, "学校法人") {
		return HoujinKakuGakko
	} else if strings.Contains(name, "弁護士法人") {
		return HoujinKakuBengoshi
	} else if strings.Contains(name, "独立行政法人") {
		return HoujinKakuDokuritsu
	} else {
		if strings.Contains(s, "宗教法人") || strings.Contains(s, "境内建物、境内地") {
			return HoujinKakuShukyo
		}
	}
	return HoujinKakuUnknown
}
