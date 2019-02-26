package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: oracle_trc trace [, trace...]")
		os.Exit(1)
	}
	for _, a := range os.Args[1:] {
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
			p := New(f, fn)
			err = p.dumpQueries()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}

}
