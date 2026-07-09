COG ?= cog
FIRST_CONVENTIONAL_COMMIT := $(shell git log --reverse --format='%H %s' | awk '$$2 ~ /^[a-z]+(\(.+\))?!?:$$/ { print $$1; exit }')

.PHONY: build bump-auto bump-major bump-minor bump-patch changelog check-commits generate install-hooks lint pre-commit pre-push test version

build:
	$(MAKE) -C backend build

generate:
	$(MAKE) -C backend generate

lint:
	$(MAKE) -C backend lint

test:
	cd backend && go test ./...

pre-commit: generate lint

pre-push: build test

install-hooks:
	$(COG) install-hook --all --overwrite

check-commits:
	@if git describe --tags --abbrev=0 >/dev/null 2>&1; then \
		$(COG) check --from-latest-tag; \
	elif [ -n "$(FIRST_CONVENTIONAL_COMMIT)" ]; then \
		$(COG) check "$(FIRST_CONVENTIONAL_COMMIT)^..HEAD"; \
	else \
		echo "No tag or conventional commit baseline found." >&2; \
		exit 1; \
	fi

changelog:
	@if git describe --tags --abbrev=0 >/dev/null 2>&1; then \
		$(COG) changelog --unified; \
	elif [ -n "$(FIRST_CONVENTIONAL_COMMIT)" ]; then \
		$(COG) changelog --unified "$(FIRST_CONVENTIONAL_COMMIT)^..HEAD"; \
	else \
		echo "No tag or conventional commit baseline found." >&2; \
		exit 1; \
	fi

version:
	$(COG) get-version --fallback 0.0.0 --tag

bump-auto:
	$(COG) bump --auto

bump-patch:
	$(COG) bump --patch

bump-minor:
	$(COG) bump --minor

bump-major:
	$(COG) bump --major
