package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/tychy/toukibo_parser/pdf"
	"github.com/tychy/toukibo_parser/toukibo"
)

func main() {
	f := flag.String("path", "testdata/pdf/sample1.pdf", "")
	flag.Parse()
	path := fmt.Sprint(*f)
	content, err := readPdf(path)
	if err != nil {
		panic(err)
	}

	h, err := toukibo.Parse(content)
	if err != nil {
		panic(err)
	}
	names, err := h.GetHoujinRepresentativeNames()
	if err != nil {
		panic(err)
	}
	fmt.Println("HoujinKaku: " + h.GetHoujinKaku())
	fmt.Println("HoujinName: " + h.GetHoujinName())
	fmt.Println("HoujinAddress: " + h.GetHoujinAddress())
	fmt.Println("HoujinRepresentativeNames: [" + strings.Join(names, ",") + "]")
	fmt.Println("HoujinDissolvedAt: " + h.GetHoujinDissolvedAt())
	fmt.Println("HoujinCapital: ", h.GetHoujinCapital())
	return
}

func readPdf(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
