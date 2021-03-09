all: build format lint freeze

build: render-ci-pipeline

render-ci-pipeline:
    ./scripts/render-ci-pipeline.sh

format: format-dhall prettier format-shfmt format-golang

lint: lint-dhall shellcheck

prettier:
    yarn run prettier

format-golang:
    ./scripts/go-format.sh

freeze: freeze-dhall

freeze-dhall:
    ./scripts/dhall-freeze.sh

test: test-golang

test-golang:
    ./scripts/go-test.sh

format-dhall:
    ./scripts/dhall-format.sh

lint-dhall:
    ./scripts/dhall-lint.sh

shellcheck:
    ./scripts/shellcheck.sh

format-shfmt:
    shfmt -w .

install:
    just install-asdf
    just install-yarn

install-yarn:
    yarn

install-asdf:
    asdf install
