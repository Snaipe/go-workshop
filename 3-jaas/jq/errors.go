package jq

import "fmt"

// SyntaxError represents errors when the jq filter is malformed.
type SyntaxError struct {
	Filter  string
	Message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("syntax error on filter %q: %v", e.Filter, e.Message)
}

// ParseError represents errors when parsing JSON data.
type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("json parse error: %v", e.Message)
}
