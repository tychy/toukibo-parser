package toukibo_parser

import (
	"bytes"

	"github.com/tychy/toukibo-parser/internal/pdf"
	"github.com/tychy/toukibo-parser/internal/toukibo"
)

func GetContentByPDFPath(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "s", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ParseByPDFPath(path string) (*toukibo.Houjin, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return nil, err
	}
	return ParseByPDFReader(r)
}

func ParseByPDFReader(r *pdf.Reader) (*toukibo.Houjin, error) {
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return nil, err
	}
	_, err = buf.ReadFrom(b)
	if err != nil {
		return nil, err
	}
	return toukibo.Parse(buf.String())
}

func ParseByPDFRawData(data []byte) (*toukibo.Houjin, error) {
	r, err := pdf.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	return ParseByPDFReader(r)
}
