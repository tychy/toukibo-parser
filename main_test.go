package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tychy/toukibo_parser/toukibo"
	"gopkg.in/yaml.v3"
)

type TestData struct {
	HoujinName    string `yaml:"HoujinName"`
	HoujinAddress string `yaml:"HoujinAddress"`
}

func TestToukiboParser(t *testing.T) {
	testCount := 10

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
			if h.GetHoujinName() != td.HoujinName {
				t.Errorf("name is not match, expected: %s, actual: %s", td.HoujinName, h.GetHoujinName())
			}

			if h.GetHoujinAddress() != td.HoujinAddress {
				t.Errorf("address is not match, expected: %s, actual: %s", td.HoujinAddress, h.GetHoujinAddress())
			}

		})
	}
}
