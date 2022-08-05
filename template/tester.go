package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"gitlab.com/mjwhitta/errors"
)

func init() {
	defer func() {
		if r := recover(); r != nil {
			panic(r.(error).Error())
		}
	}()

	var e error
	var g *gzip.Reader

	// Return if no shellcode
	if len(sc) == 0 {
		return
	}

	// Uncompress
	if g, e = gzip.NewReader(bytes.NewReader(sc)); e != nil {
		panic(errors.Newf("failed to unzip: %w", e))
	}
	defer g.Close()

	if sc, e = io.ReadAll(g); e != nil {
		panic(errors.Newf("failed to unzip: %w", e))
	}

	// Print out the shellcode for verification
	for i, b := range sc {
		fmt.Printf("%02x", b)
		if ((i + 1) % 35) == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	os.Exit(0)
}
