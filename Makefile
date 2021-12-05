BS := 1
CS := 1
SC :=

-include gomk/main.mk

go%:
	@go run ./tools/generator.go "$(BS)" "$(CS)" "$@" "$(SC)"

superclean: clean
ifeq ($(unameS),Windows)
	$(foreach d,$(wildcard cmd/*),$(shell powershell -c Remove-Item -Force -Recurse $d))
else
	@rm -fr ./cmd/*
endif
