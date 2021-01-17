package download

import "fmt"

var (
	_ error = (*InvalidResponseCode)(nil)
	_ error = (*DeadlineExceeded)(nil)
	_ error = (*Canceled)(nil)
)

// InvalidResponseCode is the error containing the invalid response code error information
type InvalidResponseCode struct {
	expected int
	got      int
}

// Error returns the InvalidResponseCode error string
func (e *InvalidResponseCode) Error() string {
	return fmt.Sprintf("Invalid response code, received '%d' expected '%d'", e.got, e.expected)
}

// DeadlineExceeded is the error containing the deadline exceeded error information
type DeadlineExceeded struct {
	url string
}

// Error returns the DeadlineExceeded error string
func (e *DeadlineExceeded) Error() string {
	return fmt.Sprintf("Download timeout exceeded for '%s'", e.url)
}

// Canceled is the error containing the cancelled error information
type Canceled struct {
	url string
}

// Error returns the Canceled error string
func (e *Canceled) Error() string {
	return fmt.Sprintf("Download canceled for '%s'", e.url)
}
