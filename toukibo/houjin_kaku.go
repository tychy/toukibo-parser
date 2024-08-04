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
	HoujinKakuIppanZaisan     HoujinkakuType = "一般財産法人"
	HoujinKakuIppanZaidan     HoujinkakuType = "一般財団法人"
	HoujinKakuNPO             HoujinkakuType = "NPO法人"
	HoujinKakuTokuteiHieiri   HoujinkakuType = "特定非営利活動法人"
	HoujinKakuPoliticalParty  HoujinkakuType = "政党"
	HoujinKakuUniversity      HoujinkakuType = "国立大学法人"
	HoujinKakuGakko           HoujinkakuType = "学校法人"
)

func FindHoujinKaku(name, s string) HoujinkakuType {
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
	} else if strings.Contains(name, "一般財産法人") {
		return HoujinKakuIppanZaisan
	} else if strings.Contains(name, "一般財団法人") {
		return HoujinKakuIppanZaidan
	} else if strings.Contains(name, "NPO法人") || strings.Contains(name, "ＮＰＯ法人") {
		return HoujinKakuNPO
	} else if strings.Contains(name, "特定非営利活動法人") {
		return HoujinKakuTokuteiHieiri
	} else if strings.Contains(name, "政党") {
		// このパターンは今のところ存在しない 2024/08/04
		return HoujinKakuPoliticalParty
	} else if strings.Contains(name, "国立大学法人") {
		return HoujinKakuUniversity
	} else if strings.Contains(name, "学校法人") {
		return HoujinKakuGakko
	} else {
		if strings.Contains(s, "宗教法人") || strings.Contains(s, "境内建物、境内地") {
			return HoujinKakuShukyo
		}

		if name == "自由民主党" || name == "日本共産党" || name == "公明党" || name == "立憲民主党" ||
			name == "国民民主党" || name == "社会民主党" || name == "日本維新の会" || name == "れいわ新選組" {
			return HoujinKakuPoliticalParty
		}
	}

	return HoujinKakuUnknown
}
