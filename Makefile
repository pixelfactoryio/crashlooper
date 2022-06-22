include global.mk

GO_LDFLAGS := -s -w
GO_LDFLAGS := -X go.pixelfactory.io/pkg/version.REVISION=$(VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X go.pixelfactory.io/pkg/version.BUILDDATE=$(BUILD_DATE) $(GO_LDFLAGS)
bin/crashlooper: $(BUILD_FILES)
	@go build -trimpath -ldflags "$(GO_LDFLAGS)" -o "$@" 

test:
	@go test -v -race -coverprofile coverage.txt -covermode atomic ./...
.PHONY: test

lint:
	@golint -set_exit_status ./...
.PHONY: lint

vet:
	@go vet ./...
.PHONY: vet
