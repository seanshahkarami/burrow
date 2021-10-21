package gopher

import (
	"bufio"
	"io"
	"strings"
)

type Scanner struct {
	scanner *bufio.Scanner
	code    string
	fields  []string
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
	}
}

func (s *Scanner) Scan() bool {
	if ok := s.scanner.Scan(); !ok {
		return false
	}
	if s.scanner.Text() == "." {
		return false
	}
	s.scanLine(s.scanner.Text())
	return true
}

func (s *Scanner) scanLine(line string) {
	if line == "" {
		s.code = ""
		s.fields = nil
		return
	}
	s.code = line[:1]
	s.fields = strings.Split(line[1:], "\t")
}

func (s *Scanner) Err() error {
	return s.scanner.Err()
}

func (s *Scanner) Code() string {
	return s.code
}

func (s *Scanner) Fields() []string {
	return s.fields
}

func (s *Scanner) Field(n int) string {
	if n >= len(s.fields) {
		return ""
	}
	return s.fields[n]
}
