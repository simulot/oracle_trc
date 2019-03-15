package trc

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type trc_scanner struct {
	*bufio.Scanner
	Line      int
	buff      *bytes.Buffer
	didBackup bool
	hasRead   bool
}

func newScanner(in io.Reader) *trc_scanner {
	s := &trc_scanner{
		Scanner: bufio.NewScanner(in),
		buff:    bytes.NewBuffer([]byte{}),
	}
	return s
}

func (s *trc_scanner) Scan() bool {
	if s.didBackup && !s.hasRead {
		return true
	}
	r := s.Scanner.Scan()
	s.buff.Reset()
	s.buff.Write(s.Scanner.Bytes())
	s.didBackup = false
	s.hasRead = false
	s.Line++
	return r
}

func (s *trc_scanner) Bytes() []byte {
	if s.didBackup {
		s.hasRead = true
	}
	return s.buff.Bytes()
}

func (s *trc_scanner) Text() string {
	if s.didBackup {
		s.hasRead = true
	}
	return s.buff.String()
}

func (s *trc_scanner) Backup() {
	if !s.didBackup {
		s.didBackup = true
		s.hasRead = false
	}
}

func (s *trc_scanner) Err() error {
	return errors.Wrapf(s.Scanner.Err(), "Error at line %d", s.Line)
}
