# Go
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gotv_plus
DIST_DIR=build
SERVER_DIR=server
OS_Linux=linux
OS_Windows=windows
OS_Mac=darwin
ARCH_386=386
ARCH_AMD64=amd64

# Replacing "RM" command for Windows PowerShell.
RM = rm -rf
ifeq ($(OS),Windows_NT)
    RM = Remove-Item -Recurse -Force
endif

# Replacing "MKDIR" command for Windows PowerShell.
MKDIR = mkdir -p
ifeq ($(OS),Windows_NT)
    MKDIR = New-Item -ItemType Directory
endif

# Replacing "CP" command for Windows PowerShell.
CP = cp -R
ifeq ($(OS),Windows_NT)
	CP = powershell -Command Copy-Item -Recurse -Force
endif

# Replacing "GOPATH" command for Windows PowerShell.
GOPATHDIR = $GOPATH
ifeq ($(OS),Windows_NT)
    GOPATHDIR = $$env:GOPATH
endif

.DEFAULT_GOAL := build-all

test:
	$(GOTEST) -v ./...
clean:
	@$(GOCLEAN)
	-@$(RM) $(DIST_DIR)/*
deps-go:
	@$(GOGET)
# Cross compile for go
build-prepare:
	@cd ./server && $(GOGET) github.com/mitchellh/gox \
	github.com/konsorten/go-windows-terminal-sequences
build-all: clean build-linux build-windows build-mac
build-linux: build-prepare
	@cd ./server && gox \
	-os="$(OS_Linux)" \
	-arch="$(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/$(BINARY_NAME)"
	$(CP) ./server/templates $(DIST_DIR)/$(BINARY_NAME)_$(OS_Linux)_$(ARCH_AMD64)/
build-windows: build-prepare
	@cd ./server && gox \
	-os="$(OS_Windows)" \
	-arch="$(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/$(BINARY_NAME)"
	$(CP) ./server/templates $(DIST_DIR)/$(BINARY_NAME)_$(OS_Windows)_$(ARCH_AMD64)/
build-mac: build-prepare
	@cd ./server && gox \
	-os="$(OS_Mac)" \
	-arch="$(ARCH_AMD64)" \
	--output "../$(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/$(BINARY_NAME)"
	$(CP) ./server/templates $(DIST_DIR)/$(BINARY_NAME)_$(OS_Mac)_$(ARCH_AMD64)/