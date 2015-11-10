# setup any dependencies for Build
GO_PACKAGES := github.com/docker/docker/pkg/term
GO_PACKAGES := $(GO_PACKAGES) github.com/docker/machine
GO_PACKAGES := $(GO_PACKAGES) github.com/stretchr/testify/assert
GO_PACKAGES := $(GO_PACKAGES) golang.org/x/crypto/ssh
GO_PACKAGES := $(GO_PACKAGES) github.com/Sirupsen/logrus
GO_PACKAGES := $(GO_PACKAGES) github.com/Azure/go-ansiterm
GO_PACKAGES := $(GO_PACKAGES) github.com/davecgh/go-spew/spew
GO_PACKAGES := $(GO_PACKAGES) github.com/pmezard/go-difflib/difflib

GO15VENDOREXPERIMENT := 1

# Cross builder helper
define godeps-get
	GOPATH=$(GOVENDORPATH) godep get $(1);
endef

define godeps-save
	GOPATH=$(GOVENDORPATH) godep save $(1);
endef

define GOVENDORPATH
$(shell godep path)
endef

define godeps-clean
	echo 'Clean Package $(1)';
	[ -d $(GOVENDORPATH)/src/$(1) ] && \
		( cd $(GOVENDORPATH)/src/$(1); \
			_PKG_CLEAN=`git rev-parse --show-toplevel`; \
			[ -d $$_PKG_CLEAN ] && rm -rf $$_PKG_CLEAN; ) || \
		echo "Skipting clean for $(1)";
endef

define godeps-vendor-gitclean
	echo 'Clean up git repos in $(1)'; \
	cd $(GOVENDORPATH)/src/$(1); \
	_GIT_ROOT=`git rev-parse --show-toplevel`; \
	[ -d $$_GIT_ROOT/.git ] && rm -rf $$_GIT_ROOT/.git || \
		echo "Skipting .git clean for $(1)";
endef

# TODO: remove this call when we OSRB
godeps-vendor-clean-oneview:
	$(call godeps-vendor-gitclean,github.com/HewlettPackard/oneview-golang)

godeps-init-oneview:
		$(call godeps-get,github.com/HewlettPackard/oneview-golang)

# END TODO:

vendor-clean:
		@rm -rf $(PREFIX)/vendor/*
		@echo cleaning up in $(PREFIX)/vendor/*

# for fresh setup so we can do godep save -r
godeps-clean: vendor-clean
		@echo "Removing all dependent packages from $(GOVENDORPATH)"
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-clean,$(GOPCKG)))
		rm -rf $(GOPATH)/src/github.com/$(GH_USER)/$(GH_REPO)

# setup a fresh GOPATH directory with what would be needed to build
# TODO: clean dependencies on oneview
godeps-init: godeps-clean godeps-init-oneview
		@echo "Get dependent packages"
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-get,$(GOPCKG)))

godeps-save:
		$(call godeps-save, $(GO_PACKAGES))

# setup the vendor folder with required packages that have been committed
# TODO: clean dependencies on oneview
godeps-vendor: godeps-vendor-clean-oneview
		echo "Placing packages into $(GOVENDORPATH)"
		[ ! -h $(PREFIX)/vendor ] && ln -s Godeps/_workspace/src vendor; \
		[ ! -d $(PREFIX)/Godeps/_workspace/src ] && mkdir -p $(PREFIX)/Godeps/_workspace/src; \
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-vendor-gitclean,$(GOPCKG)))

godeps: godeps-init godeps-save godeps-vendor
		echo "All done! run git status and commit to save any changes."

godep: godeps
