package toukibo_parser

import (
	"flag"
	"fmt"
	"strings"
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
		err := mainRun()
		if err != nil {
			fmt.Println(err)
		}
	case "find":
		err := mainFind(target)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("invalid mode")
	}
}

func mainRun() error {
	h, err := ParseByPDFPath(path)
	if err != nil {
		return err
	}
	repName, err := h.GetHoujinRepresentativeNames()
	if err != nil {
		return err
	}
	execs, err := h.GetHoujinExecutives()
	if err != nil {
		return err
	}

	execNames, err := h.GetHoujinExecutiveNames()
	if err != nil {
		return err
	}

	stock := h.GetHoujinStock()

	fmt.Println("HoujinKaku: " + h.GetHoujinKaku())
	fmt.Println("HoujinName: " + h.GetHoujinName())
	fmt.Println("HoujinAddress: " + h.GetHoujinAddress())
	fmt.Print("HoujinExecutiveValues: \n" + execs.String())
	fmt.Println("HoujinExecutiveNames: [" + strings.Join(execNames, ",") + "]")
	fmt.Println("HoujinRepresentativeNames: [" + strings.Join(repName, ",") + "]")
	fmt.Printf("HoujinCapital: %d\n", h.GetHoujinCapital())
	fmt.Printf("HoujinStock: %d\n", h.GetHoujinTotalStock())
	fmt.Print("HoujinPreferredStock: \n" + stock.String())
	fmt.Println("HoujinCreatedAt: " + h.GetHoujinCreatedAt())
	fmt.Println("HoujinBankruptedAt: " + h.GetHoujinBankruptedAt())
	fmt.Println("HoujinDissolvedAt: " + h.GetHoujinDissolvedAt())
	fmt.Println("HoujinContinuedAt: " + h.GetHoujinContinuedAt())
	return nil
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

func mainFind(s string) error {
	content, err := GetContentByPDFPath(path)
	if err != nil {
		return err
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
	return nil
}
