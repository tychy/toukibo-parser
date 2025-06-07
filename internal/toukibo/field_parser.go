package toukibo

import "strings"

// FieldProcessor はフィールド処理のインターフェース
type FieldProcessor interface {
	// Matches はフィールドがこのプロセッサで処理可能かを判定する
	Matches(s string) bool
	// Process はフィールドを処理する
	Process(h *HoujinBody, s string) bool
}

// fieldProcessor は具体的なフィールド処理の実装
type fieldProcessor struct {
	matcher func(string) bool
	process func(*HoujinBody, string) bool
}

func (fp *fieldProcessor) Matches(s string) bool {
	return fp.matcher(s)
}

func (fp *fieldProcessor) Process(h *HoujinBody, s string) bool {
	return fp.process(h, s)
}

// NewFieldProcessor は新しいFieldProcessorを作成する
func NewFieldProcessor(matcher func(string) bool, process func(*HoujinBody, string) bool) FieldProcessor {
	return &fieldProcessor{
		matcher: matcher,
		process: process,
	}
}

// FieldProcessorRegistry はフィールドプロセッサのレジストリ
type FieldProcessorRegistry struct {
	processors []FieldProcessor
}

// NewFieldProcessorRegistry は新しいレジストリを作成する
func NewFieldProcessorRegistry() *FieldProcessorRegistry {
	return &FieldProcessorRegistry{
		processors: []FieldProcessor{
			// 商号
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "商　号") || strings.Contains(s, "名　称")
				},
				func(h *HoujinBody, s string) bool {
					return h.processHoujinName(s)
				},
			),
			// 本店
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "本　店") || strings.Contains(s, "主たる事務所")
				},
				func(h *HoujinBody, s string) bool {
					return h.processHoujinAddress(s)
				},
			),
			// 資本金
			NewFieldProcessor(
				func(s string) bool {
					// ConsumeHoujinCapitalはインスタンスメソッドなので、別途定義
					return strings.HasPrefix(s, " ┃資本金の額") || strings.HasPrefix(s, " ┃特定資本金の額") ||
						strings.Contains(s, " ┃資本金") ||
						strings.Contains(s, "払込済出資総額") || strings.Contains(s, "出資の総額") ||
						strings.Contains(s, "資産の総額") || strings.Contains(s, "基本財産の総額") || strings.Contains(s, "特定資本の額") ||
						strings.Contains(s, "払い込んだ出資の")
				},
				func(h *HoujinBody, s string) bool {
					return h.processHoujinCapital(s)
				},
			),
			// 発行済株式
			NewFieldProcessor(
				func(s string) bool {
					return strings.HasPrefix(s, " ┃発行済株式の総数")
				},
				func(h *HoujinBody, s string) bool {
					return h.processHoujinStock(s)
				},
			),
			// 登記記録
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "登記記録に関する")
				},
				func(h *HoujinBody, s string) bool {
					// 登記記録はスキップするだけ
					return true
				},
			),
			// 役員
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "役員に関する事項") || strings.Contains(s, "社員に関する事項")
				},
				func(h *HoujinBody, s string) bool {
					// 役員情報は別途処理される
					return true
				},
			),
			// 会社成立日
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "会社成立の年月日") || strings.Contains(s, "法人成立の年月日")
				},
				func(h *HoujinBody, s string) bool {
					return h.ConsumeHoujinCreatedAt(s)
				},
			),
			// 破産
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "┃破　産")
				},
				func(h *HoujinBody, s string) bool {
					return h.ConsumeHoujinBankruptedAt(s)
				},
			),
			// 解散
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "┃解　散")
				},
				func(h *HoujinBody, s string) bool {
					return h.ConsumeHoujinDissolvedAt(s)
				},
			),
			// 継続
			NewFieldProcessor(
				func(s string) bool {
					return strings.Contains(s, "┃会社継続")
				},
				func(h *HoujinBody, s string) bool {
					return h.ConsumeHoujinContinuedAt(s)
				},
			),
		},
	}
}

// Process は与えられた文字列を適切なプロセッサで処理する
func (r *FieldProcessorRegistry) Process(h *HoujinBody, s string) bool {
	for _, processor := range r.processors {
		if processor.Matches(s) {
			return processor.Process(h, s)
		}
	}
	return false
}