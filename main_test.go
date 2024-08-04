package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tychy/toukibo-parser/toukibo"
	"gopkg.in/yaml.v3"
)

type TestData struct {
	HoujinKaku                string                         `yaml:"HoujinKaku"`
	HoujinName                string                         `yaml:"HoujinName"`
	HoujinAddress             string                         `yaml:"HoujinAddress"`
	HoujinExecutiveValues     []toukibo.HoujinExecutiveValue `yaml:"HoujinExecutiveValues"`
	HoujinExecutiveNames      []string                       `yaml:"HoujinExecutiveNames"`
	HoujinRepresentativeNames []string                       `yaml:"HoujinRepresentativeNames"`
	HoujinCapital             string                         `yaml:"HoujinCapital"`
	HoujinCreatedAt           string                         `yaml:"HoujinCreatedAt"`
	HoujinBankruptedAt        string                         `yaml:"HoujinBankruptedAt"`
	HoujinDissolvedAt         string                         `yaml:"HoujinDissolvedAt"`
	HoujinContinuedAt         string                         `yaml:"HoujinContinuedAt"`
}

func TestToukiboParser(t *testing.T) {
	testCount := 901

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

			if h.GetHoujinKaku() != td.HoujinKaku {
				t.Fatalf("kaku is not match,\nwant : %s,\ngot  : %s,", td.HoujinKaku, h.GetHoujinKaku())
			}

			if h.GetHoujinName() != td.HoujinName {
				t.Fatalf("name is not match,\nwant : %s,\ngot  : %s,", td.HoujinName, h.GetHoujinName())
			}

			if h.GetHoujinAddress() != td.HoujinAddress {
				t.Fatalf("address is not match,\nwant : %s,\ngot  : %s,", td.HoujinAddress, h.GetHoujinAddress())
			}

			// Exec
			execs, err := h.GetHoujinExecutives()
			if err != nil {
				t.Fatal(err)
			}

			if len(execs) != len(td.HoujinExecutiveValues) {
				t.Fatalf("executive count is not match,\nwant : %d,\ngot  : %d", len(td.HoujinExecutiveValues), len(execs))
			}
			for i, v := range execs {
				if v.Name != td.HoujinExecutiveValues[i].Name {
					t.Fatalf("executive name is not match,\nwant : %s,\ngot  : %s", td.HoujinExecutiveValues[i].Name, v.Name)
				}
				if v.Position != td.HoujinExecutiveValues[i].Position {
					t.Fatalf("executive position is not match,\nwant : %s,\ngot  : %s", td.HoujinExecutiveValues[i].Position, v.Position)
				}
			}

			// ExecutiveNames
			execNames, err := h.GetHoujinExecutiveNames()
			if err != nil {
				t.Fatal(err)
			}

			if len(execNames) != len(td.HoujinExecutiveNames) {
				t.Fatalf("executive name count is not match,\nwant : %d,\ngot  : %d", len(td.HoujinExecutiveNames), len(execNames))
			}
			for i, v := range execNames {
				if v != td.HoujinExecutiveNames[i] {
					t.Fatalf("executive name is not match,\nwant : %s,\ngot  : %s", td.HoujinExecutiveNames[i], v)
				}
			}

			// RepresentativeNames
			repNames, err := h.GetHoujinRepresentativeNames()
			if err != nil {
				t.Fatal(err)
			}
			if len(repNames) != len(td.HoujinRepresentativeNames) {
				t.Fatalf("representative name count is not match,\nwant : %d,\ngot  : %d", len(td.HoujinRepresentativeNames), len(repNames))
			}
			for i, v := range repNames {
				if v != td.HoujinRepresentativeNames[i] {
					t.Fatalf("representative name is not match,\nwant : %s,\ngot  : %s", td.HoujinRepresentativeNames[i], v)
				}
			}

			if fmt.Sprint(h.GetHoujinCapital()) != td.HoujinCapital {
				t.Fatalf("capital is not match,\nwant : %s,\ngot  : %d,", td.HoujinCapital, h.GetHoujinCapital())
			}

			if h.GetHoujinCreatedAt() != td.HoujinCreatedAt {
				t.Fatalf("created_at is not match,\nwant : %s,\ngot  : %s,", td.HoujinCreatedAt, h.GetHoujinCreatedAt())
			}

			if h.GetHoujinBankruptedAt() != td.HoujinBankruptedAt {
				t.Fatalf("bankrupted_at is not match,\nwant : %s,\ngot  : %s,", td.HoujinBankruptedAt, h.GetHoujinBankruptedAt())
			}

			if h.GetHoujinDissolvedAt() != td.HoujinDissolvedAt {
				t.Fatalf("dissolved_at is not match,\nwant : %s,\ngot  : %s,", td.HoujinDissolvedAt, h.GetHoujinDissolvedAt())
			}
			if h.GetHoujinContinuedAt() != td.HoujinContinuedAt {
				t.Fatalf("continued_at is not match,\nwant : %s,\ngot  : %s,", td.HoujinContinuedAt, h.GetHoujinContinuedAt())
			}
		})
	}
}
