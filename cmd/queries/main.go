package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/simulot/oracle_trc/queries"
	"github.com/simulot/oracle_trc/ts"
)

func main() {
	flag.Usage = func() {
		fmt.Println("Display all queries contained in trc files.")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	tsFormat := flag.String("tsFormat", "DD-MON-YYYY HH:MI:SS:FF3", "Timestamp format, oracle's way.")
	pAfter := flag.String("after", "", "Filter packets exchanged after this date. In same format as tsFormat parameter.")

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

	for _, a := range flag.Args() {
		fns, err := filepath.Glob(a)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, fn := range fns {
			fmt.Println(fn)
			err = parseFile(fn, timeParser, tAfter)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func parseFile(fn string, timeParser ts.TimeParserFn, tAfter time.Time) error {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	p := queries.New(f, fn)
	var q *queries.Query
	for err != io.EOF {
		q, err = p.Next()
		if err != nil || q == nil {
			break
		}
		if !tAfter.IsZero() && len(q.Packet.TS) > 0 {
			ts, err := timeParser(q.Packet.TS)
			if err != nil {
				return err
			}
			if tAfter.After(ts) {
				continue
			}
		}
		fmt.Println(q.String())
	}
	return nil
}
