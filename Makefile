-include gomk/main.mk
-include local/Makefile

BS := 1
CS := 1
SC :=

go%:
	@go run ./tools/generator.go "$(BS)" "$(CS)" "$@" "$(SC)"
	@go fmt ./...

superclean: clean
ifeq ($(unameS),windows)
	$(foreach d,$(wildcard cmd/*),$(shell powershell -c Remove-Item -Force -Recurse $d))
else
	@rm -f -r ./cmd/*
endif
