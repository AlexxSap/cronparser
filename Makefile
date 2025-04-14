.DEFAULT_GOAL := tests

.PHONY: tests

vet:
	go vet

tests: vet
	go test
