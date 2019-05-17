package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/simulot/oracle_trc/trc"

	"github.com/pkg/errors"
	"github.com/simulot/oracle_trc/ts"
)

type response struct {
	t   time.Time
	pk  *trc.Packet
	err error
}

var iAmDone = make(chan bool)

func main() {
	flag.Usage = func() {
		fmt.Println("Display all packets contained into given files in hexadecimal format like hex -C would do.")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	tsFormat := flag.String("tsFormat", "DD-MON-YYYY HH:MI:SS:FF3", "Timestamp format, oracle's way.")
	pAfter := flag.String("after", "", "Filter packets exchanged after this date. In same format as tsFormat parameter.")
	pSortByDate := flag.Bool("date-order", false, "Sort output by date")

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	timeParser, err := ts.GetParser(*tsFormat)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tAfter := time.Time{}
	if *pAfter != "" {
		var err error

		tAfter, err = timeParser([]byte(*pAfter))
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "Can't parse after parameter"))
			os.Exit(1)
		}
	}

	rChan := make(chan response)
	if *pSortByDate {
		go dateSortedOutput(rChan)
	} else {
		go directOutput(rChan)
	}

	for _, a := range flag.Args() {
		fns, err := filepath.Glob(a)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, fn := range fns {
			err = parseFile(fn, timeParser, tAfter, rChan)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	close(rChan)
	<-iAmDone
}

func parseFile(fn string, timeParser ts.TimeParserFn, tAfter time.Time, r chan response) error {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	p := trc.New(f, fn)
	var pk *trc.Packet
	for {
		pk, err = p.NextPacket()
		if pk == nil {
			break
		}

		if err != nil && pk == nil {
			r <- response{
				pk:  pk,
				err: err,
			}
			break
		}
		var ts time.Time
		if len(pk.TS) > 0 {
			ts, err = timeParser(pk.TS)
			if err != nil {
				return err
			}
			if tAfter.After(ts) {
				continue
			}
		}
		r <- response{
			t:   ts,
			pk:  pk,
			err: err,
		}
	}
	return nil
}

func directOutput(ch chan response) {
	for r := range ch {
		pk, err := r.pk, r.err
		if pk != nil {
			fmt.Fprintln(os.Stdout, pk.String())
		}
		if err != nil {
			fmt.Println(os.Stderr, err)
		}
	}
	close(iAmDone)
}

type responseByDate []response

func (r responseByDate) Len() int           { return len(r) }
func (r responseByDate) Less(i, j int) bool { return r[i].t.Before(r[j].t) }
func (r responseByDate) Swap(i, j int)      { r[j], r[i] = r[i], r[j] }

func dateSortedOutput(ch chan response) {
	l := []response{}

	for r := range ch {
		if r.pk != nil {
			l = append(l, r)
		}
	}
	sort.Sort(responseByDate(l))
	for _, r := range l {
		fmt.Println(r.pk.String())
	}
	close(iAmDone)
}
