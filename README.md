# GoDogGO

<a href="https://www.buymeacoffee.com/mjwhitta">🍪 Buy me a cookie</a>

[![Go Report Card](https://goreportcard.com/badge/gitlab.com/mjwhitta/godoggo)](https://goreportcard.com/report/gitlab.com/mjwhitta/godoggo)

This dog's byte is much worse than its bark.

## What is this?

This is not a real Go module. You should `git clone` it rather than
`go get` or `go install`.

This tool aims to hide a malicious payload in multiple `init()`
functions. Specifically it appends `BS` (blocksize, defaults to 1)
bytes of the payload in each `init()`. Each block is divided into `CS`
(chunk size, defaults to 1) bytes to prevent excessive calls to
`append()` for larger payloads. Then the final `init()` runs the
payload while `main()` simply keeps the process from exiting.

# Usage

Target must start with `go`. This is a limitation of the current
Makefile but can be changed after your code has been generated.

```
$ git clone https://gitlab.com/mjwhitta/godogo.git
$ cd godoggo
$ git submodule update --init
$ make gocalc SC=hex/windows/calc
$ make GOOS=windows
```

## Other examples

Use `BS` and `CS` to adjust splitting of shellcode.

```
$ make gobeacon BS=1024 SC=/path/to/beacon
$ make gostageless BS=4096 CS=512 SC=/path/to/beacon_stageless
$ make gokatz BS=1024 SC=/path/to/mimikatz
$ make GOOS=windows
```

**Note:** Compiling larger shellcode can take a considerable amount of
time. `make superclean` is your friend.

## Links

- [Source](https://gitlab.com/mjwhitta/godoggo)
