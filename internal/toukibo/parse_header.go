package toukibo

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// 閉鎖された登記簿のPDFは、作成日時の行と本店・商号の間にこの告知が入る。
// 本店・商号は位置で決めているため、取り除いてから読む。
const closedToukiboNotice = "これは閉鎖された登記簿です。"

type HoujinHeader struct {
	CreatedAt      time.Time
	CompanyName    string
	CompanyAddress string
	IsClosed       bool
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
		dateStr := ZenkakuToHankaku(matches[1])
		timeStr := ZenkakuToHankaku(matches[2])

		// 日付と時刻を time.Time 型に変換
		layout := "2006/01/02 15:04"
		dt, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", dateStr, timeStr), time.FixedZone("Asia/Tokyo", 9*3600))
		if err != nil {
			return time.Time{}, fmt.Errorf("日付と時刻の変換に失敗しました: %w", err)
		}
		return dt, nil
	} else {
		return time.Time{}, fmt.Errorf("日付と時刻が見つかりませんでした")
	}
}

func ParseHeader(s string) (*HoujinHeader, error) {
	isClosed := strings.Contains(s, closedToukiboNotice)
	arr := make([]string, 0, 4)
	for _, v := range strings.Split(s, " 　") {
		if strings.Contains(v, closedToukiboNotice) {
			v = strings.TrimSpace(strings.ReplaceAll(v, closedToukiboNotice, ""))
			if v == "" {
				continue
			}
		}
		arr = append(arr, v)
	}
	if len(arr) < 4 {
		return nil, fmt.Errorf("header has too few parts: %q", s)
	}
	createdAt, err := ReadCreatedAt(arr[0])
	if err != nil {
		return nil, err
	}
	companyAddress := trimLeadingTrailingSpace(arr[2])
	companyName := trimLeadingTrailingSpace(arr[3])
	return &HoujinHeader{CreatedAt: createdAt, CompanyAddress: companyAddress, CompanyName: companyName, IsClosed: isClosed}, nil
}
