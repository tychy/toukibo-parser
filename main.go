package main

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/tychy/toukibo_parser/pdf"
	"github.com/tychy/toukibo_parser/toukibo"
)

func main() {
	f := flag.String("path", "sample1", "")
	flag.Parse()
	path := fmt.Sprintf("sample/%s.pdf", *f)
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
	fmt.Println("代表:", names)
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
