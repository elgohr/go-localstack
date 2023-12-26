all: # nothing - to speed up advanced security scan

test:
	go test -race -timeout 0 ./...

update-linting-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest

lint: update-linting-tools
	golangci-lint run -v ./...
	govulncheck ./...

update-dependencies:
	go get -u ./...
	go mod tidy
	go mod vendor

generate:
	find . -type d -name '*fakes' | xargs rm -r
	go generate ./...

install-release-tool:
	go install github.com/elgohr/semv@latest

porcelain:
	./scripts/porcelain.sh

new-patch-release: porcelain
	./scripts/release.sh --patch

new-minor-release: porcelain
	./scripts/release.sh --minor

new-major-release: porcelain
	./scripts/release.sh --major
