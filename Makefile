COP=cover.out

.PHONY: vet lint test test-cover clean

vet:
	@- go vet ./...

lint: vet
	@- golangci-lint run

test: vet
	@- go test -race -count 1 -v -coverprofile="$(COP)" ./...

test-cover: test
	@- go tool cover -func="$(COP)"


bench:
	@- go test -count 1 -bench=. -benchmem -timeout 15m

clean:
	@- rm -f "$(COP)"
