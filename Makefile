.PHONY: clean
clean:
	rm -rf gen
	rm -rf dist

.PHONY: update
update:
	go mod tidy

.PHONY: gen
gen: clean update
	buf generate

.PHONY: test
test: gen
	go test -v ./pkg/