-include gomk/main.mk
-include local/Makefile

BS := 1
CS := 1
SC :=

go%:
	@go run ./tools/generator.go "$(BS)" "$(CS)" "$@" "$(SC)"
	@go generate "./cmd/$@"
	@go fmt ./...

superclean: clean
ifeq ($(unameS),windows)
	@$(foreach d,$(wildcard cmd/*),remove-item -force -recurse $d;)
else
	@rm -f -r ./cmd/*
endif
