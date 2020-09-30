DEPS = $(wildcard */*.go)
VERSION = $(shell git describe --always --dirty)
COMMIT_SHA1 = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date +%Y-%m-%d)

all: lint vet prometheus-kerberos-exporter

prometheus-kerberos-exporter: main.go $(DEPS)
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
	  go build -a \
		  -ldflags="-X main.version=$(VERSION) -X main.commit=$(COMMIT_SHA1) -X main.date=$(BUILD_DATE)" \
	    -installsuffix cgo -o $@ $<
	strip $@

clean:
	rm -f prometheus-kerberos-exporter

lint:
	@ go get -v golang.org/x/lint/golint
	@for file in $$(git ls-files '*.go' | grep -v '_workspace/' | grep -v 'vendor/'); do \
		export output="$$(golint $${file} | grep -v 'type name will be used as docker.DockerInfo')"; \
		[ -n "$${output}" ] && echo "$${output}" && export status=1; \
	done; \
	exit $${status:-0}

vet: main.go
	go vet $<

.PHONY: all lint vet clean vendor
