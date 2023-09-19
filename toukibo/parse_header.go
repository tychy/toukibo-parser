package toukibo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type HoujinHeader struct {
	CreatedAt      time.Time
	CompanyName    string
	CompanyAddress string
}

func (h *HoujinHeader) String() string {
	return fmt.Sprintf("Header\nPDF作成日時: %v\n商号: %s\n本店住所: %s\n",
		h.CreatedAt,
		h.CompanyName,
		h.CompanyAddress,
	)
}
func ReadCreatedAt(s string) (time.Time, error) {
	// 正規表現パターン: 全角数字で構成された日付と時刻
	pattern := "([０-９]{2,4}／[０-９]{1,2}／[０-９]{1,2})　*([０-９]{1,2}：[０-９]{1,2})"
	regex := regexp.MustCompile(pattern)

	// 抽出された日付と時刻を表示
	matches := regex.FindStringSubmatch(s)
	if len(matches) > 0 {
		// 全角数字を半角数字に変換
		dateStr := zenkakuToHankaku(matches[1])
		timeStr := zenkakuToHankaku(matches[2])

		// 日付と時刻を time.Time 型に変換
		layout := "2006/01/02 15:04"
		dt, err := time.Parse(layout, fmt.Sprintf("%s %s", dateStr, timeStr))
		if err != nil {
			return time.Time{}, fmt.Errorf("日付と時刻の変換に失敗しました: %w", err)
		}
		return dt, nil
	} else {
		return time.Time{}, fmt.Errorf("日付と時刻が見つかりませんでした")
	}
}

func ParseHeader(s string) (*HoujinHeader, error) {
	arr := strings.Split(s, " 　")
	createdAt, err := ReadCreatedAt(arr[0])
	if err != nil {
		return nil, err
	}
	companyAddress := arr[2]
	companyName := arr[3]
	return &HoujinHeader{CreatedAt: createdAt, CompanyAddress: companyAddress, CompanyName: companyName}, nil
}
