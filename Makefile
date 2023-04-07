.PHONY: tools
tools:
	go install github.com/RangelReale/trcache/cmd/troptgen

.PHONY: gen
gen: tools
	go generate ./...

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: update-dep-version
update-dep-version:
	test -n "$(TAG)"  # $$TAG
	sh -c 'go get github.com/RangelReale/trcache@$(TAG); go get github.com/RangelReale/trcache/mocks@$(TAG); go mod tidy'

git-status:
	@status=$$(git status --porcelain); \
	if [ ! -z "$${status}" ]; \
	then \
		echo "Error - working directory is dirty. Commit those changes!"; \
		exit 1; \
	fi

.PHONY: gittag
gittag: git-status
	test -n "$(TAG)"  # $$TAG
	git tag $(TAG)
