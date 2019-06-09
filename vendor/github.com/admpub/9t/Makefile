OWNER      = gongo
REPOSITORY = 9t

COMMAND  = 9t
MAIN_DIR = cmd/9t
VERSION  = $(shell grep "const Version " $(MAIN_DIR)/version.go | sed -E 's/.*"(.+)"$$/\1/')

TOP       = $(shell pwd)
BUILD_DIR = $(TOP)/pkg
DIST_DIR  = $(TOP)/dist

XC_ARCH   = "386 amd64"
XC_OS     = "darwin linux windows"
XC_OUTPUT = "$(BUILD_DIR)/{{.OS}}_{{.Arch}}/{{.Dir}}"


setup:
	go get -u github.com/tcnksm/ghr

build:
	rm -rf $(BUILD_DIR)
	mkdir $(BUILD_DIR)
	gox -os $(XC_OS) -arch $(XC_ARCH) -output $(XC_OUTPUT) ./$(MAIN_DIR)

dist: build
	rm -rf $(DIST_DIR)
	mkdir $(DIST_DIR)

	@for dir in $$(find $(BUILD_DIR) -mindepth 1 -maxdepth 1 -type d); do \
		platform=$$(basename $$dir) ; \
		archive=$(COMMAND)_$(VERSION)_$$platform ;\
		zip -j $(DIST_DIR)/$$archive.zip $$dir/* ;\
	done

	@pushd $(DIST_DIR) ; shasum -a 256 *.zip > ./SHA256SUMS ; popd

release:
	ghr -u $(OWNER) -r $(REPOSITORY) $(VERSION) $(DIST_DIR)

clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
