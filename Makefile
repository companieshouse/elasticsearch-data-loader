CHS_ENV_HOME ?= $(HOME)/.chs_env
TESTS        ?= ./...

bin          := elasticsearch-data-loader
cb           := companybindex
chs_env      := $(CHS_ENV_HOME)/global_env
source_env   := for chs_env in $(chs_envs); do test -f $$chs_env && . $$chs_env; done
xunit_output := test.xml
lint_output  := lint.txt

commit       := $(shell git rev-parse --short HEAD)
tag          := $(shell git tag -l 'v*-rc*' --points-at HEAD)
version      := $(shell if [[ -n "$(tag)" ]]; then echo $(tag) | sed 's/^v//'; else echo $(commit); fi)

.EXPORT_ALL_VARIABLES:
GO111MODULE = on

.PHONY: all
all: clean build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build: fmt
	cd ./$(cb); go build

.PHONY: test
test: test-unit test-integration

.PHONY: test-unit
test-unit:
	go test $(TESTS) -run 'Unit' -coverprofile=coverage.out

.PHONY: test-integration
test-integration:
	$(source_env); go test $(TESTS) -run 'Integration'

.PHONY: clean
clean:
	go mod tidy
	rm -f $(bin)
	rm -f *.zip
	rm -rf build-*

.PHONY: package
package:
	$(eval tmpdir:=$(shell mktemp -d build-XXXXXXXXXX))
	cp -r ./$(cb) $(tmpdir)/$(cb)
	cp ./run-elastic-search.sh $(tmpdir)/run-elastic-search.sh
	cp -r ./config $(tmpdir)/config
	zip -r $(bin)-$(version).zip $(tmpdir)
	rm -rf $(tmpdir)

.PHONY: dist
dist: clean build package

.PHONY: xunit-tests
xunit-tests: GO111MODULE=off
xunit-tests:
	go get github.com/tebeka/go2xunit
	go test -v $(TESTS) -run 'Unit' | go2xunit -output $(xunit_output)

.PHONY: lint
lint: GO111MODULE=off
lint:
	go get github.com/golang/lint/golint
	golint ./... > $(lint_output)

.PHONY: convey
convey: GO111MODULE=off
convey:
	go get github.com/smartystreets/goconvey
	goconvey
