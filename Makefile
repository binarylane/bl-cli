# NOTE: Adding a target so it shows up in the help listing
#    - The description is the text that is echoed in the first command in the target.
#    - Only 'public' targets (start with an alphanumeric character) display in the help listing.
#    - All public targets need a description

export CGO_ENABLED = 0

export GO111MODULE := on

# ORIGIN is used when testing release code
ORIGIN ?= origin
BUMP ?= patch

.PHONY: help
help:
	@echo "==> describe make commands"
	@echo ""
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null |\
	  awk -v RS= -F: \
	    '/^# File/,/^# Finished Make data base/ {if ($$1 ~ /^[a-zA-Z]/) {printf "%-20s%s\n", $$1, substr($$9, 9, length($$9)-9)}}' |\
	  sort

my_d = $(shell pwd)
OUT_D = $(shell echo $${OUT_D:-$(my_d)/builds})
DOCS_OUT = $(shell echo $${DOCS_OUT:-$(my_d)/builds/docs/yaml})

UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

GOOS = linux
ifeq ($(UNAME_S),Darwin)
  GOOS = darwin
endif

ifeq ($(GOARCH),)
  GOARCH = amd64
  ifneq ($(UNAME_M), x86_64)
    GOARCH = 386
  endif
endif

.PHONY: _build
_build:
	@echo "=> building bl via go build"
	@echo ""
	@OUT_D=${OUT_D} GOOS=${GOOS} GOARCH=${GOARCH} scripts/_build.sh
	@echo "built $(OUT_D)/bl_$(GOOS)_$(GOARCH)"

.PHONY: build
build: _build
	@echo "==> build local version"
	@echo ""
	@mv $(OUT_D)/bl_$(GOOS)_$(GOARCH) $(OUT_D)/bl
	@echo "installed as $(OUT_D)/bl"

.PHONY: native
native: build
	@echo ""
	@echo "==> The 'native' target is deprecated. Use 'make build'"

.PHONY: _build_linux_amd64
_build_linux_amd64: GOOS = linux
_build_linux_amd64: GOARCH = amd64
_build_linux_amd64: _build

.PHONY: docker_build
docker_build:
	@echo "==> build bl in local docker container"
	@echo ""
	@mkdir -p $(OUT_D)
	@docker build -f Dockerfile \
		--build-arg GOARCH=$(GOARCH) \
		. -t bl_local
	@docker run --rm \
		-v $(OUT_D):/copy \
		-it --entrypoint /bin/cp \
		bl_local /app/bl /copy/
	@docker run --rm \
		-v $(OUT_D):/copy \
		-it --entrypoint /bin/chown \
		alpine -R $(shell whoami | id -u): /copy
	@echo "Built binaries to $(OUT_D)"
	@echo "Created a local Docker container. To use, run: docker run --rm -it bl_local"

.PHONY: test_unit
test_unit:
	@echo "==> run unit tests"
	@echo ""
	go test -mod=vendor ./commands/... ./bl/... ./pkg/... .

.PHONY: test_integration
test_integration:
	@echo "==> run integration tests"
	@echo ""
	go test -v -mod=vendor ./integration

.PHONY: test
test: test_unit test_integration

.PHONY: shellcheck
shellcheck:
	@echo "==> analyze shell scripts"
	@echo ""
	@scripts/shell_check.sh

.PHONY: mocks
mocks:
	@echo "==> update mocks"
	@echo ""
	@scripts/regenmocks.sh

.PHONY: _upgrade_binarylane
_upgrade_binarylane:
	go get -u github.com/binarylane/go-binarylane

.PHONY: upgrade_binarylane
upgrade_binarylane: _upgrade_binarylane vendor mocks
	@echo "==> upgrade the binarylane version"
	@echo ""

.PHONY: vendor
vendor:
	@echo "==> vendor dependencies"
	@echo ""
	go mod vendor
	go mod tidy

.PHONY: clean
clean:
	@echo "==> remove build / release artifacts"
	@echo ""
	@rm -rf builds dist out

.PHONY: _install_github_release_notes
_install_github_release_notes:
	@GO111MODULE=off go get -u github.com/digitalocean/github-changelog-generator

.PHONY: _changelog
_changelog: _install_github_release_notes
	@scripts/changelog.sh

.PHONY: changes
changes: _install_github_release_notes
	@echo "==> list merged PRs since last release"
	@echo ""
	@changes=$(shell scripts/changelog.sh) && cat $$changes && rm -f $$changes

.PHONY: version
version:
	@echo "==> bl version"
	@echo ""
	@ORIGIN=${ORIGIN} scripts/version.sh

.PHONY: _install_sembump
_install_sembump:
	@echo "=> installing/updating sembump tool"
	@echo ""
	@GO111MODULE=off go get -u github.com/jessfraz/junk/sembump

.PHONY: tag
tag: _install_sembump
	@echo "==> BUMP=${BUMP} tag"
	@echo ""
	@ORIGIN=${ORIGIN} scripts/bumpversion.sh

.PHONY: _release
_release:
	@echo "=> releasing"
	@echo ""
	@scripts/release.sh

.PHONY: release
release:
	@echo "==> release (most recent tag, normally done by travis)"
	@echo ""
	@$(MAKE) _release

.PHONY: docs
docs:
	@echo "==> Generate YAML documentation in ${DOCS_OUT}"
	@echo ""
	@mkdir -p ${DOCS_OUT}
	@DOCS_OUT=${DOCS_OUT} go run scripts/gen-yaml-docs.go
