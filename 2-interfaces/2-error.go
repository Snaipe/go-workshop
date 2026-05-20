package main

type stringErr string

func (e stringErr) Error() string {
	return string(e)
}

// Reimplementation of fmt.Errorf(...)
func Errorf(msg string, args ...any) error {
	return stringErr(fmt.Sprintf(msg, args...))
}

// An error with more context
type URLError struct {
	URL string
	Err error
}

// Returns the error message, and causes URLError to implement the error interface
func (e URLError) Error() string {
	return fmt.Sprintf("on URL %q: %v", e.URL, e.Err)
}

// Returns the underlying error; used by errors.Is and errors.As
func (e URLError) Unwrap() error {
	return e.Err
}
