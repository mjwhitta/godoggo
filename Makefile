-include gomk/main.mk
-include local/Makefile

BS := 1
CS := 1
SC :=

hide-%:
	@go run ./tools/generator.go "$(BS)" "$(CS)" "go$(subst hide-,,$@)" "$(SC)"
	@go generate "./cmd/go$(subst hide-,,$@)"
	@go fmt ./...

superclean: clean
ifeq ($(unameS),windows)
	@$(foreach d,$(wildcard cmd/*),remove-item -force -recurse $d;)
else
	@rm -f -r ./cmd/*
endif
