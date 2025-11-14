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

# Version management targets
.PHONY: install-svu
install-svu:
	@if ! command -v svu &> /dev/null; then \
		echo "Installing svu..."; \
		SVU_VERSION="2.0.0"; \
		wget -q "https://github.com/caarlos0/svu/releases/download/v$${SVU_VERSION}/svu_$${SVU_VERSION}_linux_amd64.tar.gz"; \
		tar -xzf "svu_$${SVU_VERSION}_linux_amd64.tar.gz"; \
		sudo mv svu /usr/local/bin/; \
		rm "svu_$${SVU_VERSION}_linux_amd64.tar.gz"; \
		echo "✅ svu installed successfully"; \
	else \
		echo "✅ svu is already installed"; \
	fi

.PHONY: current-version
current-version: install-svu
	@echo "Current version:"
	@svu current || echo "v0.0.0 (no tags yet)"

.PHONY: next-version
next-version: install-svu
	@echo "Next version will be:"
	@NEXT=$$(svu next 2>/dev/null || echo ""); \
	if [ -z "$$NEXT" ]; then \
		echo "No version bump needed (no conventional commits since last tag)"; \
	else \
		echo "$$NEXT"; \
		echo ""; \
		echo "Changes since last tag:"; \
		CURRENT=$$(svu current 2>/dev/null || echo ""); \
		if [ -z "$$CURRENT" ]; then \
			git log --oneline | head -10; \
		else \
			git log $${CURRENT}..HEAD --oneline; \
		fi \
	fi

.PHONY: release
release: install-svu
	@echo "Creating release tag..."
	@NEXT=$$(svu next 2>/dev/null || echo ""); \
	if [ -z "$$NEXT" ]; then \
		echo "❌ No version bump needed. No conventional commits found."; \
		echo "Use 'make release-patch' to force a patch release."; \
		exit 1; \
	fi; \
	CURRENT=$$(svu current 2>/dev/null || echo "v0.0.0"); \
	echo "Current version: $$CURRENT"; \
	echo "New version: $$NEXT"; \
	echo ""; \
	read -p "Create and push tag $$NEXT? [y/N] " -n 1 -r; \
	echo ""; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a "$$NEXT" -m "Release $$NEXT"; \
		git push origin "$$NEXT"; \
		echo "✅ Tag $$NEXT created and pushed"; \
	else \
		echo "❌ Release cancelled"; \
	fi

.PHONY: release-major
release-major: install-svu
	@echo "Creating MAJOR release..."
	@NEXT=$$(svu major); \
	CURRENT=$$(svu current 2>/dev/null || echo "v0.0.0"); \
	echo "Current version: $$CURRENT"; \
	echo "New version: $$NEXT"; \
	echo ""; \
	read -p "Create and push tag $$NEXT? [y/N] " -n 1 -r; \
	echo ""; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a "$$NEXT" -m "Release $$NEXT"; \
		git push origin "$$NEXT"; \
		echo "✅ Tag $$NEXT created and pushed"; \
	else \
		echo "❌ Release cancelled"; \
	fi

.PHONY: release-minor
release-minor: install-svu
	@echo "Creating MINOR release..."
	@NEXT=$$(svu minor); \
	CURRENT=$$(svu current 2>/dev/null || echo "v0.0.0"); \
	echo "Current version: $$CURRENT"; \
	echo "New version: $$NEXT"; \
	echo ""; \
	read -p "Create and push tag $$NEXT? [y/N] " -n 1 -r; \
	echo ""; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a "$$NEXT" -m "Release $$NEXT"; \
		git push origin "$$NEXT"; \
		echo "✅ Tag $$NEXT created and pushed"; \
	else \
		echo "❌ Release cancelled"; \
	fi

.PHONY: release-patch
release-patch: install-svu
	@echo "Creating PATCH release..."
	@NEXT=$$(svu patch); \
	CURRENT=$$(svu current 2>/dev/null || echo "v0.0.0"); \
	echo "Current version: $$CURRENT"; \
	echo "New version: $$NEXT"; \
	echo ""; \
	read -p "Create and push tag $$NEXT? [y/N] " -n 1 -r; \
	echo ""; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		git tag -a "$$NEXT" -m "Release $$NEXT"; \
		git push origin "$$NEXT"; \
		echo "✅ Tag $$NEXT created and pushed"; \
	else \
		echo "❌ Release cancelled"; \
	fi
