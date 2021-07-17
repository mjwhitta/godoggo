BS := 1
CS := 1
SC :=

-include gomk/main.mk

go%: havego
	@go run ./tools/generator.go "$(BS)" "$(CS)" "$@" "$(SC)"

superclean: clean
	@rm -fr cmd/*
