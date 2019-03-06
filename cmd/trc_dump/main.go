package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/simulot/oracle_trc/trc"

	"github.com/pkg/errors"
	"github.com/simulot/oracle_trc/ts"
)

func main() {
	tsFormat := flag.String("tsFormat", "DD-MON-YYYY HH:MI:SS:FF3", "Timestamp format, oracle's way.")
	pAfter := flag.String("after", "", "Filter queries executed after this date. In same format as tsFormat parameter.")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: trc_dump trace [, trace...]")
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
				if err != nil || pk == nil {
					break
				}
				if !tAfter.IsZero() && len(pk.TS) > 0 {
					ts, err := timeParser(pk.TS)
					if err != nil {
						break
					}
					if tAfter.After(ts) {
						break
					}
				}
				fmt.Println(pk.String())
				var input string
				fmt.Scanln(&input)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
