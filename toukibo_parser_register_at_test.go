package toukibo_parser

import (
	"fmt"
	"testing"
)

type registerAtWant struct {
	Name       string
	Position   string
	RegisterAt string
}

func TestToukiboParserRegisterAt(t *testing.T) {
	tests := []struct {
		sample int
		want   []registerAtWant
	}{
		{
			sample: 1,
			want: []registerAtWant{
				{Name: "山岸大致", Position: "取締役"},
				{Name: "山岸大致", Position: "代表取締役"},
			},
		},
		{
			sample: 92,
			want: []registerAtWant{
				{Name: "関藤竜也", Position: "取締役", RegisterAt: "令和4年10月13日"},
				{Name: "河村晃平", Position: "取締役", RegisterAt: "令和4年10月13日"},
				{Name: "高杉慧", Position: "取締役", RegisterAt: "令和4年10月13日"},
				{Name: "柏木彩", Position: "取締役", RegisterAt: "令和4年10月13日"},
				{Name: "徳山耕平", Position: "取締役", RegisterAt: "令和4年10月13日"},
				{Name: "関藤竜也", Position: "代表取締役", RegisterAt: "令和4年10月13日"},
				{Name: "小川敬介", Position: "監査役", RegisterAt: "令和4年10月13日"},
				{Name: "田上沙織", Position: "監査役", RegisterAt: "令和4年5月17日"},
				{Name: "堀口拓也", Position: "監査役", RegisterAt: "令和4年5月17日"},
			},
		},
		{
			sample: 2,
			want: []registerAtWant{
				{Name: "大熊由美子", Position: "取締役"},
				{Name: "大熊徹也", Position: "取締役"},
				{Name: "大熊理恵子", Position: "取締役", RegisterAt: "平成31年4月24日"},
				{Name: "大熊徹也", Position: "代表取締役", RegisterAt: "平成29年5月26日"},
				{Name: "大熊堅二", Position: "監査役"},
			},
		},
		{
			sample: 10,
			want: []registerAtWant{
				{Name: "久保正人", Position: "取締役", RegisterAt: "令和4年7月11日"},
				{Name: "齋藤俊哉", Position: "取締役", RegisterAt: "令和4年7月11日"},
				{Name: "原口洋一", Position: "取締役", RegisterAt: "令和4年7月11日"},
				{Name: "渡利勝也", Position: "取締役", RegisterAt: "令和4年7月11日"},
				{Name: "長井智一", Position: "取締役", RegisterAt: "令和4年7月11日"},
				{Name: "松下勝則", Position: "取締役", RegisterAt: "令和4年7月11日"},
				{Name: "齋藤俊哉", Position: "代表取締役", RegisterAt: "令和4年7月11日"},
				{Name: "伊東幹雄", Position: "監査役", RegisterAt: "令和2年7月20日"},
			},
		},
		{
			sample: 20,
			want: []registerAtWant{
				{Name: "松井由美子", Position: "取締役", RegisterAt: "令和1年9月10日"},
				{Name: "荻野聡", Position: "取締役", RegisterAt: "令和2年10月21日"},
				{Name: "伊澤伸人", Position: "取締役", RegisterAt: "令和2年10月21日"},
				{Name: "圓岡孝文", Position: "取締役", RegisterAt: "令和2年10月29日"},
				{Name: "大澤祐介", Position: "取締役", RegisterAt: "令和2年10月29日"},
				{Name: "松井由美子", Position: "代表取締役", RegisterAt: "令和1年9月10日"},
			},
		},
		{
			sample: 227,
			want: []registerAtWant{
				{Name: "渥美勤", Position: "取締役"},
				{Name: "渥美則彰", Position: "取締役", RegisterAt: "平成16年5月12日"},
				{Name: "渥美かず子", Position: "取締役", RegisterAt: "平成20年5月13日"},
				{Name: "渥美勤", Position: "代表取締役", RegisterAt: "平成16年5月12日"},
			},
		},
		{
			sample: 691,
			want: []registerAtWant{
				{Name: "熊谷駒吉", Position: "無限責任社員"},
				{Name: "大堀唯延", Position: "有限責任社員"},
				{Name: "熊谷千代三郎", Position: "有限責任社員"},
				{Name: "藤木金藏", Position: "有限責任社員"},
				{Name: "熊谷森治", Position: "有限責任社員"},
				{Name: "熊谷藏十", Position: "有限責任社員"},
				{Name: "熊谷駒吉", Position: "代表社員", RegisterAt: "平成18年5月12日"},
			},
		},
		{
			sample: 1133,
			want: []registerAtWant{
				{Name: "池見幸浩", Position: "取締役", RegisterAt: "令和5年6月6日"},
				{Name: "大畑貴文", Position: "取締役", RegisterAt: "令和5年6月6日"},
				{Name: "田中祐輔", Position: "取締役", RegisterAt: "令和5年6月6日"},
				{Name: "雨宮玲於奈", Position: "取締役", RegisterAt: "令和5年6月6日"},
				{Name: "大滝直道", Position: "取締役", RegisterAt: "令和6年6月10日"},
				{Name: "池見幸浩", Position: "代表取締役", RegisterAt: "令和5年6月6日"},
				{Name: "田中仁", Position: "監査役", RegisterAt: "令和5年6月6日"},
				{Name: "福本麻美", Position: "監査役", RegisterAt: "令和6年6月10日"},
				{Name: "幸森鈴華", Position: "監査役", RegisterAt: "令和3年6月10日"},
			},
		},
		{
			sample: 1320,
			want: []registerAtWant{
				{Name: "富田仁", Position: "取締役", RegisterAt: "令和3年11月29日"},
				{Name: "橋本文敏", Position: "取締役", RegisterAt: "令和5年7月21日"},
				{Name: "富田仁", Position: "代表取締役", RegisterAt: "令和3年11月29日"},
				{Name: "三浦振一郎", Position: "監査役", RegisterAt: "令和2年3月12日"},
			},
		},
		{
			sample: 1644,
			want: []registerAtWant{
				{Name: "永田周一", Position: "取締役", RegisterAt: "令和6年4月9日"},
				{Name: "淺田慎二", Position: "取締役", RegisterAt: "令和6年4月9日"},
				{Name: "ナオタケ・ムラヤマ", Position: "取締役", RegisterAt: "令和6年4月9日"},
				{Name: "今野穣", Position: "取締役", RegisterAt: "令和6年4月9日"},
				{Name: "川田尚吾", Position: "取締役", RegisterAt: "令和7年4月15日"},
				{Name: "細村拓也", Position: "取締役", RegisterAt: "令和7年8月15日"},
				{Name: "永田周一", Position: "代表取締役", RegisterAt: "令和7年4月15日"},
				{Name: "外川香", Position: "監査役", RegisterAt: "令和7年4月15日"},
				{Name: "森谷均", Position: "監査役", RegisterAt: "令和5年12月27日"},
				{Name: "門松優介", Position: "監査役", RegisterAt: "令和6年4月19日"},
			},
		},
		{
			sample: 521,
			want: []registerAtWant{
				{Name: "渥美秀登", Position: "社員", RegisterAt: "令和3年3月31日"},
				{Name: "渥美友季大", Position: "代表社員", RegisterAt: "令和3年3月31日"},
			},
		},
		{
			sample: 796,
			want: []registerAtWant{
				{Name: "伸和工業株式会社", Position: "業務執行社員", RegisterAt: "平成25年10月24日"},
				{Name: "上海電力日本株式会社", Position: "業務執行社員", RegisterAt: "平成26年4月23日"},
				{Name: "伸和工業株式会社", Position: "代表社員", RegisterAt: "平成25年10月24日"},
				{Name: "西村浩", Position: "職務執行者", RegisterAt: "平成25年10月24日"},
				{Name: "上海電力日本株式会社", Position: "代表社員", RegisterAt: "令和5年7月12日"},
				{Name: "施伯紅", Position: "職務執行者", RegisterAt: "令和5年7月12日"},
			},
		},
		{
			sample: 891,
			want: []registerAtWant{
				{Name: "弓納持弘", Position: "無限責任社員", RegisterAt: "令和3年5月20日"},
				{Name: "弓納持陽子", Position: "有限責任社員"},
				{Name: "弓納持惠子", Position: "有限責任社員"},
				{Name: "弓納持和男", Position: "有限責任社員", RegisterAt: "平成18年5月24日"},
				{Name: "弓納持彰", Position: "有限責任社員", RegisterAt: "平成18年5月24日"},
				{Name: "弓納持弘", Position: "代表社員", RegisterAt: "令和3年5月20日"},
			},
		},
		{
			sample: 1047,
			want: []registerAtWant{
				{Name: "髙月勝守", Position: "取締役", RegisterAt: "令和3年11月4日"},
				{Name: "髙月廣海", Position: "取締役", RegisterAt: "令和3年11月4日"},
				{Name: "藥師寺大思", Position: "取締役", RegisterAt: "令和4年4月20日"},
				{Name: "髙月勝守", Position: "代表取締役", RegisterAt: "令和3年11月4日"},
				{Name: "藥師寺大思", Position: "代表取締役", RegisterAt: "令和6年3月6日"},
				{Name: "長谷雅史", Position: "監査役", RegisterAt: "令和4年4月20日"},
				{Name: "石井教文", Position: "破産管財人", RegisterAt: "令和6年7月1日"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("sample%d", tt.sample), func(t *testing.T) {
			t.Parallel()
			h, err := ParseByPDFPath(fmt.Sprintf("testdata/pdf/sample%d.pdf", tt.sample))
			if err != nil {
				t.Fatal(err)
			}
			execs, err := h.GetHoujinExecutives()
			if err != nil {
				t.Fatal(err)
			}
			if len(execs) != len(tt.want) {
				t.Fatalf("executive count: want %d, got %d (%+v)", len(tt.want), len(execs), execs)
			}
			for i, got := range execs {
				want := tt.want[i]
				if got.Name != want.Name || got.Position != want.Position || got.RegisterAt != want.RegisterAt {
					t.Fatalf("executive[%d]: got {Name:%q Position:%q RegisterAt:%q}, want %+v",
						i, got.Name, got.Position, got.RegisterAt, want)
				}
			}
		})
	}
}
