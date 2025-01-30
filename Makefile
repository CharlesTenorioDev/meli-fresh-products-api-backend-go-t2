# Build for Linux or Mac
.PHONY: build-linux-macos
build-linux-macos:
	go build -ldflags "-w" ./cmd/server/main.go

# Run in development mode
.PHONY: dev
dev:
	air

# Keep the tests running in the background
.PHONY: dev-test
dev-test:
	nodemon --watch ./**/*.go --exec 'make tests'

# This command will run the tests for the project
.PHONY tests:
tests:
	go test -v ./...

# This command will generate a coverage report for the project
.PHONY: coverage
coverage:
	go test -cover -coverprofile=coverage.out ./... && go tool cover -func="coverage.out"


a:
	make coverage | grep total | awk '{print substr($3, 1, length($3)-1)}'

# This command will display the coverage report in an HTML file
.PHONY: coverage-html
coverage-html: coverage
	go tool cover -html="coverage.out"


# This command will run the linter for the project
.PHONY: static-linter
static-linter:
	staticcheck ./...

.PHONY: linter
linter:
	golangci-lint run