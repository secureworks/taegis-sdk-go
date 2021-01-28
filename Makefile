PACKAGE_DIRS = $(shell go list -f '{{ .Dir }}' ./...)
PACKAGES = $(shell go list ./...)

flake:
	go test -short -v ./... -test.failfast -test.count 10

lint:
	golangci-lint run

vet:
	go vet ./... || go clean ./...; go vet ./...

coverage: $(patsubst %,%.coverage,$(PACKAGES))
	@rm -f .gocoverage/cover.txt
	gocovmerge .gocoverage/*.out > coverage.txt
	go tool cover -html=coverage.txt -o .gocoverage/index.html
	go tool cover -func=coverage.txt

coverage-html: coverage
	go tool cover -html=coverage.txt

%.coverage:
	@[ -d .gocoverage ] || mkdir .gocoverage
	go test -covermode=count -coverprofile=.gocoverage/$(subst /,-,$*).out $* -v



