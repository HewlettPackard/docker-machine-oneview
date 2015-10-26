# setup any dependencies for Build
GO_PACKAGES := github.com/docker/machine github.com/stretchr/testify/assert
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

vendor-clean:
		@rm -rf $(PREFIX)/vendor/*
		@echo cleaning up in $(PREFIX)/vendor/*

# for fresh setup so we can do godep save -r
godeps-clean: vendor-clean
		@echo "Removing all dependent packages from $(GOVENDORPATH)"
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-clean,$(GOPCKG)))
		rm -rf $(GOPATH)/src/github.com/$(GH_USER)/$(GH_REPO)

# setup a fresh GOPATH directory with what would be needed to build
godeps-init: godeps-clean
		@echo "Pulling required packages into $(GOPATH)"
		mkdir -p $(GOPATH)/src/github.com/$(GH_USER)
		ln -s $(PREFIX) $(GOPATH)/src/github.com/$(GH_USER)/$(GH_REPO)
		@echo "Get dependent packages"
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-get,$(GOPCKG)))

godeps-save:
		$(call godeps-save, $(GO_PACKAGES))

# setup the vendor folder with required packages that have been committed
godeps-vendor:
		echo "Placing packages into $(GOVENDORPATH)"
		[ ! -h $(PREFIX)/vendor ] && ln -s Godeps/_workspace/src vendor; \
		[ ! -d $(PREFIX)/Godeps/_workspace/src ] && mkdir -p $(PREFIX)/Godeps/_workspace/src; \
		$(foreach GOPCKG,$(GO_PACKAGES),$(call godeps-vendor-gitclean,$(GOPCKG)))
		# GOPATH=$(GOVENDORPATH) godep restore; TODO: this makes Godep path submodules

godeps: godeps-init godeps-save godeps-vendor
		echo "All done! run git status and commit to save any changes."

godep: godeps godeps-vendor
