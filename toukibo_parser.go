package toukibo_parser

import (
	"bytes"
	"fmt"

	"github.com/tychy/toukibo-parser/internal/pdf"
	"github.com/tychy/toukibo-parser/internal/toukibo"
)

func GetContentByPDFPath(path string) (string, error) {
	r, err := pdf.Open(path)
	if err != nil {
		return "", err
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
	
	content := buf.String()
	// Check if content is valid (not empty and contains expected characters)
	if len(content) == 0 {
		return "", fmt.Errorf("PDF content is empty")
	}
	
	// Check if content contains valid Japanese text or expected patterns
	// Encrypted PDFs often produce garbled text
	validChars := 0
	for _, r := range content {
		if r >= 0x3040 && r <= 0x309F || // Hiragana
		   r >= 0x30A0 && r <= 0x30FF || // Katakana
		   r >= 0x4E00 && r <= 0x9FAF || // Kanji
		   r >= 0x20 && r <= 0x7E {      // ASCII
			validChars++
		}
	}
	
	// If less than 50% of characters are valid, likely encrypted or corrupted
	// Also check for specific patterns that indicate valid toukibo content
	validRatio := float64(validChars) / float64(len([]rune(content)))
	if validRatio < 0.5 {
		// Check for common toukibo patterns
		hasToukiboPattern := false
		patterns := []string{"商　号", "本　店", "会社成立", "登記記録", "法人番号"}
		for _, pattern := range patterns {
			if bytes.Contains([]byte(content), []byte(pattern)) {
				hasToukiboPattern = true
				break
			}
		}
		if !hasToukiboPattern {
			return "", fmt.Errorf("PDF content appears to be encrypted or corrupted (valid ratio: %.2f)", validRatio)
		}
	}
	
	return content, nil
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
