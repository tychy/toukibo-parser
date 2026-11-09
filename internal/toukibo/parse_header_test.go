package toukibo

import "testing"

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantName    string
		wantAddress string
		wantClosed  bool
	}{
		{
			name:        "通常の登記簿",
			input:       "２０２２／０３／０８　１８：４０　現在の情報です。 　 　東京都中央区日本橋二丁目２番３号サンプルビル４階 　株式会社サンプル  ",
			wantName:    "株式会社サンプル",
			wantAddress: "東京都中央区日本橋二丁目２番３号サンプルビル４階",
			wantClosed:  false,
		},
		{
			// 閉鎖された登記簿は、作成日時の行と本店・商号の間に告知の行が入る。
			name:        "閉鎖された登記簿",
			input:       "２０２６／１１／０９　１２：５０　現在の情報です。 　 これは閉鎖された登記簿です。 　 　東京都港区六本木六丁目８－１０ 　株式会社サンプル  ",
			wantName:    "株式会社サンプル",
			wantAddress: "東京都港区六本木六丁目８－１０",
			wantClosed:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ParseHeader(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if h.CompanyName != tt.wantName {
				t.Fatalf("name: want %q, got %q", tt.wantName, h.CompanyName)
			}
			if h.CompanyAddress != tt.wantAddress {
				t.Fatalf("address: want %q, got %q", tt.wantAddress, h.CompanyAddress)
			}
			if h.IsClosed != tt.wantClosed {
				t.Fatalf("isClosed: want %v, got %v", tt.wantClosed, h.IsClosed)
			}
		})
	}
}
