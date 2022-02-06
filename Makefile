# Go parameters
GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -p 1

# PATHS
PROJECT_DIR       = $(shell pwd -L)
TOOLSDIR          = $(PROJECT_DIR)/tools
BIN_DIR           = $(PROJECT_DIR)/bin

# tools needed for the project
TOOLS := github.com/golang/mock/mockgen

# ------------------------------------------------------------------------------
# tools
# ------------------------------------------------------------------------------
tools: $(TOOLSDIR) $(BIN_DIR) ## install tools
	@ $(MAKE) $(BIN_DIR)

updatetools: $(TOOLSDIR) ## update all tools
	@ cd $(TOOLSDIR) && GOFLAGS="-mod=" go get -u -v $(TOOLS)
	@ $(MAKE) tools


all: test build
test:
	$(GOTEST)  ./...
clean:
	$(GOCLEAN)

# ------------------------------------------------------------------------------
# cleanup
# ------------------------------------------------------------------------------
clean-tools: ## remove tooling
	@ rm -rf $(BIN_DIR)

# ------------------------------------------------------------------------------
# directories
# ------------------------------------------------------------------------------
$(BIN_DIR):
	@ cd $(TOOLSDIR) && GOBIN=$(BIN_DIR) GOFLAGS="-mod=" go install $(TOOLS)

$(TOOLSDIR):
	@ mkdir -p $(TOOLSDIR)