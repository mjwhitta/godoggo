-include gomk/main.mk
-include local/Makefile

BS := 1
CS := 1
SC :=

go%:
	@go run ./tools/generator.go "$(BS)" "$(CS)" "$@" "$(SC)"
	@go generate "./cmd/$@"
	@go fmt ./...

mr: fmt
	@make GOOS=darwin reportcard spellcheck vslint
	@make GOOS=linux reportcard spellcheck vslint
	@make CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows \
	    reportcard spellcheck vslint
	@make test

superclean: clean
ifeq ($(unameS),windows)
	@$(foreach d,$(wildcard cmd/*),remove-item -force -recurse $d;)
else
	@rm -f -r ./cmd/*
endif
