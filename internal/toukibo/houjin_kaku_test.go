package toukibo_test

import (
	"testing"

	"github.com/tychy/toukibo-parser/internal/toukibo"
)

func TestIsTokushuHoujin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "特殊法人のケース",
			input:    "日本電信電話株式会社",
			expected: true,
		},
		{
			name:     "特殊法人でないケース",
			input:    "株式会社テスト",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toukibo.IsTokushuHoujin(tt.input)
			if result != tt.expected {
				t.Errorf("IsTokushuHoujin(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsNinkaHoujin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "認可法人のケース",
			input:    "日本銀行",
			expected: true,
		},
		{
			name:     "認可法人でないケース",
			input:    "株式会社テスト",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toukibo.IsNinkaHoujin(tt.input)
			if result != tt.expected {
				t.Errorf("IsNinkaHoujin(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsTokubetsuMinkanHoujin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "特別民間法人のケース",
			input:    "日本公認会計士協会",
			expected: true,
		},
		{
			name:     "特別民間法人でないケース",
			input:    "株式会社テスト",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toukibo.IsTokubetsuMinkanHoujin(tt.input)
			if result != tt.expected {
				t.Errorf("IsTokubetsuMinkanHoujin(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPoliticalParty(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "政党のケース",
			input:    "自由民主党",
			expected: true,
		},
		{
			name:     "政党でないケース",
			input:    "株式会社テスト",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toukibo.IsPoliticalParty(tt.input)
			if result != tt.expected {
				t.Errorf("IsPoliticalParty(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindHoujinKaku(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		content  string
		expected toukibo.HoujinkakuType
	}{
		{
			name:     "政党のケース",
			input:    "自由民主党",
			content:  "",
			expected: toukibo.HoujinKakuPoliticalParty,
		},
		{
			name:     "認可法人のケース",
			input:    "日本銀行",
			content:  "",
			expected: toukibo.HoujinKakuNinka,
		},
		{
			name:     "特殊法人のケース",
			input:    "日本郵便株式会社",
			content:  "",
			expected: toukibo.HoujinKakuTokushu,
		},
		{
			name:     "特別民間法人のケース",
			input:    "日本公認会計士協会",
			content:  "",
			expected: toukibo.HoujinKakuTokubetsuMinkan,
		},
		{
			name:     "株式会社のケース",
			input:    "株式会社テスト",
			content:  "",
			expected: toukibo.HoujinKakuKabusiki,
		},
		{
			name:     "株式会社の前後に空白があるケース",
			input:    " 株式会社テスト ",
			content:  "",
			expected: toukibo.HoujinKakuKabusiki,
		},
		{
			name:     "有限会社のケース",
			input:    "有限会社テスト",
			content:  "",
			expected: toukibo.HoujinKakuYugen,
		},
		{
			name:     "合同会社のケース",
			input:    "合同会社テスト",
			content:  "",
			expected: toukibo.HoujinKakuGoudou,
		},
		{
			name:     "合資会社のケース",
			input:    "合資会社テスト",
			content:  "",
			expected: toukibo.HoujinKakuGousi,
		},
		{
			name:     "合名会社のケース",
			input:    "合名会社テスト",
			content:  "",
			expected: toukibo.HoujinKakuGoumei,
		},
		{
			name:     "特定目的会社のケース",
			input:    "特定目的会社テスト",
			content:  "",
			expected: toukibo.HoujinKakuTokuteiMokuteki,
		},
		{
			name:     "協同組合のケース",
			input:    "協同組合テスト",
			content:  "",
			expected: toukibo.HoujinKakuKyodou,
		},
		{
			name:     "労働組合のケース",
			input:    "労働組合テスト",
			content:  "",
			expected: toukibo.HoujinKakuRoudou,
		},
		{
			name:     "森林組合のケース",
			input:    "森林組合テスト",
			content:  "",
			expected: toukibo.HoujinKakuSinrin,
		},
		{
			name:     "生活衛生同業組合のケース",
			input:    "生活衛生同業組合テスト",
			content:  "",
			expected: toukibo.HoujinKakuSeikatuEisei,
		},
		{
			name:     "信用金庫のケース",
			input:    "信用金庫テスト",
			content:  "",
			expected: toukibo.HoujinKakuSinyou,
		},
		{
			name:     "商工会のケース",
			input:    "商工会テスト",
			content:  "",
			expected: toukibo.HoujinKakuShokoukai,
		},
		{
			name:     "公益財団法人のケース",
			input:    "公益財団法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuKoueki,
		},
		{
			name:     "農事組合のケース",
			input:    "農事組合テスト",
			content:  "",
			expected: toukibo.HoujinKakuNouji,
		},
		{
			name:     "管理組合法人のケース",
			input:    "管理組合法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuKanriKumiai,
		},
		{
			name:     "医療法人のケース",
			input:    "医療法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuIryo,
		},
		{
			name:     "司法書士法人のケース",
			input:    "司法書士法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuSihoshosi,
		},
		{
			name:     "税理士法人のケース",
			input:    "税理士法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuZeirishi,
		},
		{
			name:     "社会福祉法人のケース",
			input:    "社会福祉法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuShakaifukusi,
		},
		{
			name:     "一般社団法人のケース",
			input:    "一般社団法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuIppanShadan,
		},
		{
			name:     "公益社団法人のケース",
			input:    "公益社団法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuKouekiShadan,
		},
		{
			name:     "一般財産法人のケース",
			input:    "一般財産法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuIppanZaisan,
		},
		{
			name:     "一般財団法人のケース",
			input:    "一般財団法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuIppanZaidan,
		},
		{
			name:     "NPO法人のケース（半角）",
			input:    "NPO法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuNPO,
		},
		{
			name:     "NPO法人のケース（全角）",
			input:    "ＮＰＯ法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuNPO,
		},
		{
			name:     "特定非営利活動法人のケース",
			input:    "特定非営利活動法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuTokuteiHieiri,
		},
		{
			name:     "国立大学法人のケース",
			input:    "国立大学法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuUniversity,
		},
		{
			name:     "学校法人のケース",
			input:    "学校法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuGakko,
		},
		{
			name:     "弁護士法人のケース",
			input:    "弁護士法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuBengoshi,
		},
		{
			name:     "独立行政法人のケース",
			input:    "独立行政法人テスト",
			content:  "",
			expected: toukibo.HoujinKakuDokuritsu,
		},
		{
			name:     "宗教法人のケース（contentに宗教法人が含まれる）",
			input:    "テスト法人",
			content:  "宗教法人",
			expected: toukibo.HoujinKakuShukyo,
		},
		{
			name:     "宗教法人のケース（contentに境内建物、境内地が含まれる）",
			input:    "テスト法人",
			content:  "境内建物、境内地",
			expected: toukibo.HoujinKakuShukyo,
		},
		{
			name:     "不明なケース",
			input:    "テスト法人",
			content:  "",
			expected: toukibo.HoujinKakuUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toukibo.FindHoujinKaku(tt.input, tt.content)
			if result != tt.expected {
				t.Errorf("FindHoujinKaku(%q, %q) = %v, want %v", tt.input, tt.content, result, tt.expected)
			}
		})
	}
}
