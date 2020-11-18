package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/simulot/oracle_trc/queries"
	"github.com/simulot/oracle_trc/ts"
)

type response struct {
	t   time.Time
	q   *queries.Query
	err error
}

var iAmDone = make(chan bool)

func main() {
	flag.Usage = func() {
		fmt.Println("Display all queries contained in trc files.")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	tsFormat := flag.String("tsFormat", "DD-MON-YYYY HH:MI:SS:FF3", "Timestamp format, oracle's way.")
	pAfter := flag.String("after", "", "Filter packets exchanged after this date. In same format as tsFormat parameter.")
	pSortByDate := flag.Bool("date-order", false, "Sort output by date")
	rowsToRead := flag.Int("row", 25, "Number of response row to extract, 0 to turn off the feature.")

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
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "Can't parse 'after' parameter"))
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
			fmt.Println(fn)
			err = parseFile(fn, timeParser, tAfter, rowsToRead, rChan)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	close(rChan)
	<-iAmDone
}

func parseFile(fn string, timeParser ts.TimeParserFn, tAfter time.Time, rowsToRead int, r chan response) error {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	p := queries.New(f, fn, rowsToRead)
	var q *queries.Query
	for err != io.EOF {
		q, err = p.Next()
		if err != nil || q == nil {
			r <- response{
				q:   nil,
				err: err,
			}
			break
		}
		var ts time.Time
		if len(q.Packet.TS) > 0 {
			ts, err = timeParser(q.Packet.TS)
			if err != nil {
				return err
			}
			if tAfter.After(ts) {
				continue
			}
		}
		r <- response{
			t:   ts,
			q:   q,
			err: err,
		}
	}
	return nil
}

func directOutput(ch chan response) {
	for r := range ch {
		q, err := r.q, r.err
		if q != nil {
			fmt.Fprintln(os.Stdout, q.String())
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
		if r.q != nil {
			l = append(l, r)
		}
	}
	sort.Sort(responseByDate(l))
	for _, r := range l {
		fmt.Println(r.q.String())
	}
	close(iAmDone)
}
