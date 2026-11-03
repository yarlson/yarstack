CRAP4GO := github.com/unclebob/crap4go/cmd/crap4go@v0.0.0-20260521190544-bee16dbdadb4
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.4.0
CROSS_TARGETS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64
WORKFLOWS := ../../../.github/workflows/validate.yml ../../../.github/workflows/release-yarbrain.yml

.PHONY: check ci fmt-check lint vet test build mod-check release-check workflow-check race crap no-cgo cross-build vuln snapshot release-assets

check: lint vet test build mod-check release-check workflow-check

ci: check race crap no-cgo cross-build vuln snapshot release-assets

fmt-check:
	@golangci-lint version | grep -Eq 'version 2\.' || { echo "golangci-lint v2 is required; run mise install"; exit 1; }
	golangci-lint fmt --diff

lint: fmt-check
	golangci-lint run ./...

vet:
	go vet ./...

test:
	go test ./...

build:
	go build ./...

mod-check:
	go mod tidy -diff
	go mod verify

release-check:
	goreleaser check

workflow-check:
	actionlint $(WORKFLOWS)

snapshot:
	goreleaser release --snapshot --clean

# Stage the goreleaser binaries under their release names, with checksums.txt
# and the install scripts, and verify every binary against checksums.txt.
release-assets:
	rm -rf release
	mkdir release
	jq -r '.[] | select(.type == "Binary" and (.name | startswith("yarbrain_"))) | .path + " " + .name' dist/artifacts.json | \
		while read -r path name; do cp "$$path" "release/$$name"; done
	cp dist/checksums.txt install.sh install.ps1 release/
	cd release && sha256sum --check --strict checksums.txt

race:
	go test -race ./...

crap:
	@report="$$(go run $(CRAP4GO))" || { printf '%s\n' "$$report"; exit 1; }; \
	printf '%s\n' "$$report"; \
	printf '%s\n' "$$report" | awk '$$(NF - 1) ~ /%$$/ && ($$NF + 0) >= 15 { failed = 1 } END { exit failed }' || \
		{ echo "Every CRAP score must be below 15."; exit 1; }

no-cgo:
	CGO_ENABLED=0 go test ./...
	CGO_ENABLED=0 go build ./...

cross-build:
	@set -eu; for target in $(CROSS_TARGETS); do \
		CGO_ENABLED=0 GOOS="$${target%/*}" GOARCH="$${target#*/}" go build -o /dev/null ./cmd/yarbrain; \
	done

vuln:
	go run $(GOVULNCHECK) ./...
