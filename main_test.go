package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tychy/toukibo_parser/toukibo"
	"gopkg.in/yaml.v3"
)

type TestData struct {
	HoujinName                string   `yaml:"HoujinName"`
	HoujinAddress             string   `yaml:"HoujinAddress"`
	HoujinRepresentativeNames []string `yaml:"HoujinRepresentativeNames"`
	HoujinDissolvedAt         string   `yaml:"HoujinDissolvedAt"`
}

func TestToukiboParser(t *testing.T) {
	testCount := 80

	for i := 1; i <= testCount; i++ {
		t.Run(fmt.Sprintf("test%d", i), func(t *testing.T) {
			pdfFileName := fmt.Sprintf("testdata/pdf/sample%d.pdf", i)
			yamlFileName := fmt.Sprintf("testdata/yaml/sample%d.yaml", i)
			content, err := readPdf(pdfFileName)
			if err != nil {
				t.Fatal(err)
			}
			h, err := toukibo.Parse(content)
			if err != nil {
				t.Fatal(err)
			}
			yamlContent, err := os.ReadFile(yamlFileName)
			if err != nil {
				t.Fatal(err)
			}
			td := TestData{}

			err = yaml.Unmarshal([]byte(yamlContent), &td)
			if err != nil {
				t.Fatal(err)
			}

			// check
			n, err := h.GetHoujinRepresentativeNames()
			if err != nil {
				t.Fatal(err)
			}
			if len(n) != len(td.HoujinRepresentativeNames) {
				t.Fatalf("representative name count is not match,\nwant : %d,\ngot  : %d", len(td.HoujinRepresentativeNames), len(n))
			}
			for i, v := range n {
				if v != td.HoujinRepresentativeNames[i] {
					t.Fatalf("representative name is not match,\nwant : %s,\ngot  : %s", td.HoujinRepresentativeNames[i], v)
				}
			}

			if h.GetHoujinName() != td.HoujinName {
				t.Fatalf("name is not match,\nwant : %s,\ngot  : %s,", td.HoujinName, h.GetHoujinName())
			}

			if h.GetHoujinAddress() != td.HoujinAddress {
				t.Fatalf("address is not match,\nwant : %s,\ngot  : %s,", td.HoujinAddress, h.GetHoujinAddress())
			}

			if h.GetHoujinDissolvedAt() != td.HoujinDissolvedAt {
				t.Fatalf("dissolved_at is not match,\nwant : %s,\ngot  : %s,", td.HoujinDissolvedAt, h.GetHoujinDissolvedAt())
			}

		})
	}
}
