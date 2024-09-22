package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/tychy/toukibo-parser/pdf"
	"github.com/tychy/toukibo-parser/toukibo"
)

var (
	mode   string
	path   string
	target string
)

func main() {
	flag.StringVar(&mode, "mode", "run", "run or find")
	flag.StringVar(&path, "path", "testdata/pdf/sample1.pdf", "pdf file path")
	flag.StringVar(&target, "target", "", "")
	flag.Parse()

	switch mode {
	case "run":
		mainRun()
	case "find":
		mainFind(target)
	default:
		fmt.Println("invalid mode")
	}
}

func mainRun() {
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

	stock := h.GetHoujinStock()

	fmt.Println("HoujinKaku: " + h.GetHoujinKaku())
	fmt.Println("HoujinName: " + h.GetHoujinName())
	fmt.Println("HoujinAddress: " + h.GetHoujinAddress())
	fmt.Print("HoujinExecutiveValues: \n" + execs.String())
	fmt.Println("HoujinExecutiveNames: [" + strings.Join(execNames, ",") + "]")
	fmt.Println("HoujinRepresentativeNames: [" + strings.Join(repName, ",") + "]")
	fmt.Printf("HoujinCapital: %d\n", h.GetHoujinCapital())
	fmt.Printf("HoujinStock: %d\n", stock.Total)
	fmt.Print("HoujinPreferredStock: \n" + stock.String())
	fmt.Println("HoujinCreatedAt: " + h.GetHoujinCreatedAt())
	fmt.Println("HoujinBankruptedAt: " + h.GetHoujinBankruptedAt())
	fmt.Println("HoujinDissolvedAt: " + h.GetHoujinDissolvedAt())
	fmt.Println("HoujinContinuedAt: " + h.GetHoujinContinuedAt())
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mainFind(s string) {
	content, err := readPdf(path)
	if err != nil {
		panic(err)
	}

	if strings.Contains(content, s) {
		fmt.Println("found in " + path)
		// 前後を表示
		for {
			idx := strings.Index(content, s)
			if idx == -1 {
				break
			}
			fmt.Println(content[max(0, idx-60):min(len(content), idx+240)])
			content = content[idx+1:]
		}
	}
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
