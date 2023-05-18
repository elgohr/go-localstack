all: # nothing - to speed up advanced security scan

test:
	go test -race ./...

generate:
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
