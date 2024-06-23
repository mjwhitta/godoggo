-include gomk/main.mk
-include local/Makefile

BS := 1
CS := 1
SC :=

clean: clean-default
ifeq ($(unameS),windows)
ifneq ($(wildcard resource.syso),)
	@remove-item -force ./resource_windows*.syso
endif
else
	@rm -f ./resource_windows*.syso
endif

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
