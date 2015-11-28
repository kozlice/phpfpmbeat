build:
	go build

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm phpfpmbeat || true
