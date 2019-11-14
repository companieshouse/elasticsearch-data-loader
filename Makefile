CHS_ENV_HOME ?= $(HOME)/.chs_env
TESTS        ?= ./...

bin          := elasticsearch-data-loader
cb           := companybindex
test_path    := ./test
chs_env      := $(CHS_ENV_HOME)/global_env
source_env   := for chs_env in $(chs_envs); do test -f $$chs_env && . $$chs_env; done
xunit_output := test.xml
lint_output  := lint.txt

commit       := $(shell git rev-parse --short HEAD)
tag          := $(shell git tag -l 'v*-rc*' --points-at HEAD)
version      := $(shell if [[ -n "$(tag)" ]]; then echo $(tag) | sed 's/^v//'; else echo $(commit); fi)

.PHONY: all
all: build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: deps
deps:
	go get github.com/mitchellh/gox
	go get github.com/companieshouse/$(bin)/$(cb)

.PHONY: build
build: deps fmt
	cd ./$(cb); go build

.PHONY: test-deps
test-deps: deps
	go get github.com/smartystreets/goconvey
	go get github.com/smartystreets/assertions

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit: test-deps
	@set -a; go test $(TESTS) -run 'Unit'

.PHONY: test-integration
test-integration: test-deps
	$(source_env); go test $(TESTS) -run 'Integration'

.PHONY: convey
convey: clean build
	$(source_env); goconvey

.PHONY: clean
clean:
	rm -rf ./$(cb)/$(cb) ./$(bin)-*.zip $(test_path)

.PHONY: package
package: deps
	$(eval tmpdir:=$(shell mktemp -d build-XXXXXXXXXX))
	cp -r ./$(cb) $(tmpdir)/$(cb)
	cp ./run-elastic-search.sh $(tmpdir)/run-elastic-search.sh
	cp -r ./config $(tmpdir)/config
	zip -r $(bin)-$(version).zip $(tmpdir)
	rm -rf $(tmpdir)

.PHONY: dist
dist: clean build package

.PHONY: xunit-tests
xunit-tests: test-deps
	go get github.com/tebeka/go2xunit
	@set -a; go test -v $(TESTS) -run 'Unit' | go2xunit -output $(xunit_output)

.PHONY: lint
lint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	gometalinter ./... --deadline=30s > $(lint_output); true