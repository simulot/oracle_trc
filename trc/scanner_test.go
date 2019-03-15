package trc

import (
	"bufio"
	"strings"
	"testing"
)

func scanAndTestString(t *testing.T, scanner *trc_scanner, expectedB bool, want string) {
	b := scanner.Scan()
	if b != expectedB {
		t.Errorf("Scan()==%t, expected %t", b, expectedB)
		return
	}
	got := scanner.Text()
	if got != want {
		t.Errorf("Text()==%s, expected %s", got, want)
	}
}

func Test_trc_scanner(t *testing.T) {
	scanner := newScanner(bufio.NewReader(strings.NewReader(`Line 1
Line 2
Line 3
Line 4
`)))

	scanAndTestString(t, scanner, true, "Line 1")
	scanAndTestString(t, scanner, true, "Line 2")
	scanner.Backup()
	scanAndTestString(t, scanner, true, "Line 2")
	scanAndTestString(t, scanner, true, "Line 3")
	scanner.Backup()
	scanner.Backup()
	scanAndTestString(t, scanner, true, "Line 3")
	scanAndTestString(t, scanner, true, "Line 4")
	scanAndTestString(t, scanner, false, "")

}
