package main

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mjwhitta/errors"
)

func copyFile(from string, to string) error {
	var e error
	var src *os.File
	var dst *os.File

	// Open src file
	if src, e = os.Open(from); e != nil {
		return errors.Newf("failed to open %s: %w", from, e)
	}

	// Create dst file
	if dst, e = os.Create(to); e != nil {
		return errors.Newf("failed to create %s: %w", to, e)
	}

	// Copy contents
	if _, e = io.Copy(dst, src); e != nil {
		return errors.Newf("failed to copy contents: %w", e)
	}

	return nil
}

func copyTemplateFiles(name string) error {
	var e error

	e = copyFile(
		filepath.Join("template", "main.go"),
		filepath.Join("cmd", name, "main.go"),
	)
	if e != nil {
		return e
	}

	e = copyFile(
		filepath.Join("template", "first.go"),
		filepath.Join("cmd", name, "0.go"),
	)
	if e != nil {
		return e
	}

	e = copyFile(
		filepath.Join("template", "last.go"),
		filepath.Join("cmd", name, "z.go"),
	)
	if e != nil {
		return e
	}

	return nil
}

func init() {
	// Parse cli args
	flag.Parse()

	// Exit if wrong number of cli args
	if flag.NArg() != 4 {
		os.Exit(1)
	}
}

func main() {
	var b []byte
	var e error
	var hexsc string
	var name string
	var r *regexp.Regexp
	var sc []byte
	var scFile string

	// Store cli args
	name = flag.Arg(2)
	scFile = flag.Arg(3)

	// Validate file exists
	if _, e = os.Stat(scFile); (e != nil) && os.IsNotExist(e) {
		panic(errors.Newf("file %s not found", scFile))
	} else if e != nil {
		e = errors.Newf(
			"failed to get file info for %s: %w",
			scFile,
			e,
		)
		panic(e)
	}

	// Make cmd dir
	os.RemoveAll(filepath.Join("cmd", name))
	os.MkdirAll(filepath.Join("cmd", name), os.ModePerm)

	// Copy template files
	if e = copyTemplateFiles(name); e != nil {
		panic(e)
	}

	// Read scFile
	if b, e = os.ReadFile(scFile); e != nil {
		panic(errors.Newf("failed to read %s: %w", scFile, e))
	}

	// Get hex string ignoring comments
	r = regexp.MustCompile(`\s+`)
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		} else if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}

		hexsc += r.ReplaceAllString(line, "")
	}

	// Decode hex to []byte
	if sc, e = hex.DecodeString(hexsc); e != nil {
		panic(errors.Newf("failed to decode hex: %w", e))
	}

	// Gzip the bytes
	if sc, e = zipUp(sc); e != nil {
		panic(e)
	}

	writeFiles(name, sc)
}

func nextFile(
	f *os.File,
	block int,
	blocks int,
	blocksize int,
	path string,
) (*os.File, error) {
	var e error
	var fn string
	var fs string

	// Close file
	if f != nil {
		if e = f.Close(); e != nil {
			return nil, errors.Newf("failed to close file: %w", e)
		}
	}

	if blocks == 0 {
		return nil, nil
	}

	// Get new filename
	fs = "%0" + strconv.Itoa(len(strconv.Itoa(blocks))) + "d"
	fn = fmt.Sprintf(fs, (block/blocksize)+1) + ".go"
	fn = filepath.Join(path, fn)

	// Open new file
	if f, e = os.Create(fn); e != nil {
		return nil, errors.Newf("failed to create %s: %w", fn, e)
	}

	return f, nil
}

func writeFiles(name string, sc []byte) error {
	var b []byte
	var blocks int
	var blocksize int
	var chunksize int
	var e error
	var f *os.File
	var footer string = "}\n"
	var header string = strings.Join(
		[]string{"package main", "", "func init() {", ""},
		"\n",
	)

	// Get block and chunk size
	if blocksize, e = strconv.Atoi(flag.Arg(0)); e != nil {
		return errors.Newf(
			"failed to parse %s as blocksize: %w",
			flag.Arg(0),
			e,
		)
	}
	if chunksize, e = strconv.Atoi(flag.Arg(1)); e != nil {
		return errors.Newf(
			"failed to parse %s as chunksize: %w",
			flag.Arg(1),
			e,
		)
	}

	// Determine number of blocks
	blocks = len(sc) / blocksize
	if (len(sc) % blocksize) != 0 {
		blocks++
	}

	// Create numerous go files
	b = []byte{}
	for i, c := range sc {
		if (i % blocksize) == 0 {
			// Write partial chunk
			b = writeSC(b, f)

			if f != nil {
				// Write footer
				f.WriteString(footer)
			}

			// Get next file
			f, e = nextFile(
				f,
				i,
				blocks,
				blocksize,
				filepath.Join("cmd", name),
			)
			if e != nil {
				return e
			}

			// Write header
			f.WriteString(header)
		} else if (i % blocksize % chunksize) == 0 {
			// Write chunk
			b = writeSC(b, f)
		}

		b = append(b, c)
	}

	// Write partial chunk
	writeSC(b, f)
	f.WriteString(footer)
	if _, e = nextFile(f, 0, 0, 0, ""); e != nil {
		return e
	}

	return nil
}

func writeSC(b []byte, f *os.File) []byte {
	if len(b) > 0 {
		f.WriteString("\tsc = append(sc, ")

		for _, c := range b {
			f.WriteString(fmt.Sprintf("%#x,", c))
		}

		f.WriteString(")\n")
	}

	return []byte{}
}

func zipUp(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	var e error
	var g *gzip.Writer = gzip.NewWriter(&buf)

	if _, e = g.Write(b); e != nil {
		return []byte{}, errors.Newf("failed to write gzip: %w", e)
	}

	if e = g.Close(); e != nil {
		return []byte{}, errors.Newf("failed to close gzip: %w", e)
	}

	return buf.Bytes(), nil
}
