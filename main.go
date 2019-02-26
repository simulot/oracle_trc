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
	for _, a := range os.Args {
		fn, err := filepath.Glob(a)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, f := range fn {
			f, err := os.Open(f)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer f.Close()
			p := New(f)
			p.dumpQueries()
		}
	}

}
