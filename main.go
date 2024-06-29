package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/tychy/toukibo-parser/pdf"
	"github.com/tychy/toukibo-parser/toukibo"
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
	repName, err := h.GetHoujinRepresentativeNames()
	if err != nil {
		panic(err)
	}
	execs, err := h.GetHoujinExecutives()
	if err != nil {
		panic(err)
	}

	execNames, err := h.GetHoujinExecutiveNames()
	if err != nil {
		panic(err)
	}

	fmt.Println("HoujinKaku: " + h.GetHoujinKaku())
	fmt.Println("HoujinName: " + h.GetHoujinName())
	fmt.Println("HoujinAddress: " + h.GetHoujinAddress())
	fmt.Print("HoujinExecutiveValues: \n" + execs.String())
	fmt.Println("HoujinExecutiveNames: [" + strings.Join(execNames, ",") + "]")
	fmt.Println("HoujinRepresentativeNames: [" + strings.Join(repName, ",") + "]")
	fmt.Println("HoujinCapital: ", h.GetHoujinCapital())
	fmt.Println("HoujinCreatedAt: " + h.GetHoujinCreatedAt())
	fmt.Println("HoujinBankruptedAt: " + h.GetHoujinBankruptedAt())
	fmt.Println("HoujinDissolvedAt: " + h.GetHoujinDissolvedAt())
	fmt.Println("HoujinContinuedAt: " + h.GetHoujinContinuedAt())
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
