.PHONY: test lint build validate

test:
	go test ./...

build:
	go build -o bin/mwctl ./cmd/mwctl

validate:
	go test ./...
	bash -n scripts/mw_quick_survey.sh
