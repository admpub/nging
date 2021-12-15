# general
WORKDIR = $(PWD)

# coverage
COVERAGE_REPORT = coverage.txt
COVERAGE_PROFILE = profile.out
COVERAGE_MODE = atomic

coverage:
	@cd $(WORKDIR); \
	echo "" > $(COVERAGE_REPORT); \
	for dir in `find . -name "*.go" | grep -o '.*/' | sort | uniq`; do \
		go test -v  -race $$dir -coverprofile=$(COVERAGE_PROFILE) -covermode=$(COVERAGE_MODE); \
		if [ $$? != 0 ]; then \
			exit 2; \
		fi; \
		if [ -f $(COVERAGE_PROFILE) ]; then \
			cat $(COVERAGE_PROFILE) >> $(COVERAGE_REPORT); \
			rm $(COVERAGE_PROFILE); \
		fi; \
	done; \