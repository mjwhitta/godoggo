package main

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func copyFile(from string, to string) {
	var b []byte
	var e error
	var f *os.File

	// Read from file
	b, e = ioutil.ReadFile(from)
	if e != nil {
		panic(e)
	}

	// Create to file
	if f, e = os.Create(to); e != nil {
		panic(e)
	}

	// Copy contents
	if _, e = f.Write(b); e != nil {
		panic(e)
	}
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
	var blocks int
	var blocksize int
	var chunksize int
	var e error
	var f *os.File
	var footer string
	var header string
	var hexsc string
	var name string
	var r *regexp.Regexp
	var sc []byte
	var scFile string

	// Store cli args
	blocksize, _ = strconv.Atoi(flag.Arg(0))
	chunksize, _ = strconv.Atoi(flag.Arg(1))
	name = flag.Arg(2)
	scFile = flag.Arg(3)

	// Validate file exists
	if _, e = os.Stat(scFile); (e != nil) && os.IsNotExist(e) {
		os.Exit(2)
	} else if e != nil {
		panic(e)
	}

	// Make cmd dir
	os.RemoveAll(filepath.Join("cmd", name))
	os.MkdirAll(filepath.Join("cmd", name), os.ModePerm)

	// Copy template files
	copyFile(
		filepath.Join("template", "main.go"),
		filepath.Join("cmd", name, "main.go"),
	)
	copyFile(
		filepath.Join("template", "first.go"),
		filepath.Join("cmd", name, "0.go"),
	)
	copyFile(
		filepath.Join("template", "last.go"),
		filepath.Join("cmd", name, "z.go"),
	)

	// Read scFile
	if b, e = ioutil.ReadFile(scFile); e != nil {
		panic(e)
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
		panic(e)
	}

	// Gzip the bytes
	if sc, e = zipUp(sc); e != nil {
		panic(e)
	}

	// Determine number of blocks
	blocks = len(sc) / blocksize
	if (len(sc) % blocksize) != 0 {
		blocks++
	}

	// Setup header and footer
	footer = "}\n"
	header = strings.Join(
		[]string{"package main", "", "func init() {", ""},
		"\n",
	)

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
			f = nextFile(
				f,
				i,
				blocks,
				blocksize,
				filepath.Join("cmd", name),
			)

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
	nextFile(f, 0, 0, 0, "")
}

func nextFile(
	f *os.File,
	block int,
	blocks int,
	blocksize int,
	path string,
) *os.File {
	var e error
	var fn string
	var fs string

	// Close file
	if f != nil {
		if e = f.Close(); e != nil {
			panic(e)
		}
	}

	if blocks == 0 {
		return nil
	}

	// Get new filename
	fs = "%0" + strconv.Itoa(len(strconv.Itoa(blocks))) + "d"
	fn = fmt.Sprintf(fs, (block/blocksize)+1) + ".go"

	// Open new file
	f, e = os.Create(filepath.Join(path, fn))
	if e != nil {
		panic(e)
	}
	if f == nil {
		panic(errors.New("Failed to open file"))
	}

	return f
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
		return []byte{}, e
	}

	if e = g.Close(); e != nil {
		return []byte{}, e
	}

	return buf.Bytes(), nil
}
