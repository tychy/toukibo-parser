package pdf

import "fmt"

// PDFError represents errors that occur during PDF processing
type PDFError struct {
	Op      string // operation being performed
	Offset  int64  // file offset where error occurred
	Message string // error message
	Err     error  // underlying error, if any
}

func (e *PDFError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("pdf: %s at offset %d: %s: %v", e.Op, e.Offset, e.Message, e.Err)
	}
	return fmt.Sprintf("pdf: %s at offset %d: %s", e.Op, e.Offset, e.Message)
}

func (e *PDFError) Unwrap() error {
	return e.Err
}

// Common error types
var (
	ErrMalformedPDF   = fmt.Errorf("malformed PDF")
	ErrInvalidStream  = fmt.Errorf("invalid stream")
	ErrInvalidObject  = fmt.Errorf("invalid object")
	ErrMissingObject  = fmt.Errorf("missing object")
	ErrInvalidOperator = fmt.Errorf("invalid operator")
)