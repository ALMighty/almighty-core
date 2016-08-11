#===============================================================================
# Testing has become a rather big and interconnected topic and that's why it
# has arrived in it's own file.
#
# We have to types of tests available:
#
#  1. unit tests and
#  2. integration tests.
#
# While the unit tests can be executed fairly simply be running `go test`, the
# integration tests have a little bit more setup going on. That's why they are
# split up in to tests.
#
# Usage
# -----
# If you want to run the unit tests, type
#
#     $ make test-unit
#
# To run the integration tests, type
#
#     $ make test-integration
#
# To run both tests, type
#
#     $ make test-all
#
# To output unit-test coverage profile information for each function, type
#
#     $ make coverage-unit
#
# To generate unit-test HTML representation of coverage profile (opens a browser), type
#
#     $ make coverage-unit-html
#
# If you replace the "unit" with "integration" you get the same for integration
# tests.
#
# To output all coverage profile information for each function, type
#
#     $ make coverage
#
# Artifacts and coverage modes
# ----------------------------
# Each package generates coverage outputs under tmp/coverage/$(PACKAGE) where
# $(PACKAGE) resolves to the Go package. Here's an example of a coverage file
# for the package "github.com/almighty/almighty-core/models" with coverage mode
# "set" generated by the unit tests:
#
#   tmp/coverage/github.com/almighty/almighty-core/models/coverage.unit.mode-set
#
# For unit-tests all results are combined into this file:
#
#   tmp/coverage.unit.mode-$(COVERAGE_MODE)
#
# For integration-tests all results are combined into this file:
#
#   tmp/coverage.integration.mode-$(COVERAGE_MODE)
#
# The overall coverage gets combined into this file:
#
#   tmp/coverage.mode-$(COVERAGE_MODE)
#
# The $(COVERAGE_MODE) in each filename indicates what coverage mode was used.
#
# These are possible coverage modes (see https://blog.golang.org/cover):
#
# 	set: did each statement run? (default)
# 	count: how many times did each statement run?
# 	atomic: like count, but counts precisely in parallel programs
#
# To choose another coverage mode, simply prefix the invovation of `make`:
#
#     $ COVERAGE_MODE=count make test-unit
#===============================================================================

# mode can be: set, count, or atomic
COVERAGE_MODE ?= set

# By default use localhost or specify manually during make invocation:
#
# 	ALMIGHTY_DB_HOST=somehost make test-integration
#
ALMIGHTY_DB_HOST ?= localhost

# Output directory for coverage information
COV_DIR = $(TMP_PATH)/coverage

# Files that combine package coverages for unit- and integration-tests separately
COV_PATH_UNIT = $(TMP_PATH)/coverage.unit.mode-$(COVERAGE_MODE)
COV_PATH_INTEGRATION = $(TMP_PATH)/coverage.integration.mode-$(COVERAGE_MODE)

# File that stores overall coverge for all packages and unit- and integration-tests
COV_PATH_OVERALL = $(TMP_PATH)/coverage.mode-$(COVERAGE_MODE)

#-------------------------------------------------------------------------------
# Normal test targets
#
# These test targets are the ones that will be invoked from the outside. If
# they are called and the artifacts already exist, then the artifacts will
# first be cleaned and recreated. This ensures that the tests are always
# executed.
#-------------------------------------------------------------------------------

.PHONY: test-all
## Runs test-unit and test-integration targets.
test-all: prebuild-check test-unit test-integration

.PHONY: test-unit
## Runs the unit tests and produces coverage files for each package.
test-unit: prebuild-check clean-coverage-unit $(COV_PATH_UNIT)

.PHONY: test-integration
## Runs the integration tests and produces coverage files for each package.
test-integration: prebuild-check clean-coverage-integration $(COV_PATH_INTEGRATION)

.PHONY: integration-test-env-prepare
## Prepares all services needed to run the integration tests
integration-test-env-prepare:
ifndef DOCKER_COMPOSE_BIN
	$(error The "$(DOCKER_COMPOSE_BIN_NAME)" executable could not be found in your PATH)
endif
	@$(DOCKER_COMPOSE_BIN) -f $(CUR_DIR)/.make/docker-compose.integration-test.yaml up -d

.PHONY: integration-test-env-tear-down
## Tears down all services needed to run the integration tests
integration-test-env-tear-down:
ifndef DOCKER_COMPOSE_BIN
	$(error The "$(DOCKER_COMPOSE_BIN_NAME)" executable could not be found in your PATH)
endif
	@$(DOCKER_COMPOSE_BIN) -f $(CUR_DIR)/.make/docker-compose.integration-test.yaml down

#-------------------------------------------------------------------------------
# Inspect coverage of unit tests or integration tests in either pure
# console mode or in a browser (*-html).
#
# If the test coverage files to be evaluated already exist, then no new
# tests are executed. If they don't exist, we first run the tests.
#-------------------------------------------------------------------------------


$(COV_PATH_OVERALL): $(COV_PATH_UNIT) $(COV_PATH_INTEGRATION) $(GOCOVMERGE_BIN)
	@$(GOCOVMERGE_BIN) $(COV_PATH_UNIT) $(COV_PATH_INTEGRATION) > $(COV_PATH_OVERALL)

# Console coverage output:

.PHONY: coverage-unit
## Output coverage profile information for each function (only based on unit-tests).
## Re-runs unit-tests if coverage information is outdated.
coverage-unit: prebuild-check $(COV_PATH_UNIT)
	@go tool cover -func=$(COV_PATH_UNIT)

.PHONY: coverage-integration
## Output coverage profile information for each function (only based on integration tests).
## Re-runs integration-tests if coverage information is outdated.
coverage-integration: prebuild-check $(COV_PATH_INTEGRATION)
	@go tool cover -func=$(COV_PATH_INTEGRATION)

.PHONY: coverage-all
## Output coverage profile information for each function.
## Re-runs unit- and integration-tests if coverage information is outdated.
coverage-all: prebuild-check clean-coverage-overall $(COV_PATH_OVERALL)
	@go tool cover -func=$(COV_PATH_OVERALL)

# HTML coverage output:

.PHONY: coverage-unit-html
## Generate HTML representation (and show in browser) of coverage profile (based on unit tests).
## Re-runs unit tests if coverage information is outdated.
coverage-unit-html: prebuild-check $(COV_PATH_UNIT)
	@go tool cover -html=$(COV_PATH_UNIT)

.PHONY: coverage-integration-html
## Generate HTML representation (and show in browser) of coverage profile (based on integration tests).
## Re-runs integration tests if coverage information is outdated.
coverage-integration-html: prebuild-check $(COV_PATH_INTEGRATION)
	@go tool cover -html=$(COV_PATH_INTEGRATION)

.PHONY: coverage-all-html
## Output coverage profile information for each function.
## Re-runs unit- and integration-tests if coverage information is outdated.
coverage-all-html: prebuild-check clean-coverage-overall $(COV_PATH_OVERALL)
	@go tool cover -html=$(COV_PATH_OVERALL)

# Experimental:

.PHONY: gocov-unit-annotate
## (EXPERIMENTAL) Show actual code and how it is covered with unit tests.
##                This target only runs the tests if the coverage file does exist.
gocov-unit-annotate: prebuild-check $(GOCOV_BIN) $(COV_PATH_UNIT)
	@$(GOCOV_BIN) convert $(COV_PATH_UNIT) | $(GOCOV_BIN) annotate -

.PHONY: .gocov-unit-report
.gocov-unit-report: prebuild-check $(GOCOV_BIN) $(COV_PATH_UNIT)
	@$(GOCOV_BIN) convert $(COV_PATH_UNIT) | $(GOCOV_BIN) report

.PHONY: gocov-integration-annotate
## (EXPERIMENTAL) Show actual code and how it is covered with integration tests.
##                This target only runs the tests if the coverage file does exist.
gocov-integration-annotate: prebuild-check $(GOCOV_BIN) $(COV_PATH_INTEGRATION)
	@$(GOCOV_BIN) convert $(COV_PATH_INTEGRATION) | $(GOCOV_BIN) annotate -

.PHONY: .gocov-integration-report
.gocov-integration-report: prebuild-check $(GOCOV_BIN) $(COV_PATH_INTEGRATION)
	@$(GOCOV_BIN) convert $(COV_PATH_INTEGRATION) | $(GOCOV_BIN) report

#-------------------------------------------------------------------------------
# Test artifacts are coverage files for unit and integration tests.
#-------------------------------------------------------------------------------

# The test-package function executes tests for a package and saves the collected
# coverage output to a directory. After storing the coverage information it is
# also appended to a file of choice (without the "mode"-line)
#
# Parameters:
#  1. Test name (e.g. "unit" or "integration")
#  2. package name "github.com/almighty/almighty-core/model"
#  3. File in which to combine the output
#  4. (optional) parameters for "go test" command
define test-package
$(eval TEST_NAME := $(1))
$(eval PACKAGE_NAME := $(2))
$(eval COMBINED_OUT_FILE := $(3))
$(eval EXTRA_TEST_PARAMS := $(4))

@mkdir -p $(COV_DIR)/$(PACKAGE_NAME);
$(eval COV_OUT_FILE := $(COV_DIR)/$(PACKAGE_NAME)/coverage.$(TEST_NAME).mode-$(COVERAGE_MODE))
@ALMIGHTY_DB_HOST=$(ALMIGHTY_DB_HOST) go test $(PACKAGE_NAME) -v -coverprofile $(COV_OUT_FILE) -covermode=$(COVERAGE_MODE) -timeout 10m $(EXTRA_TEST_PARAMS);

@if [ -e "$(COV_OUT_FILE)" ]; then \
	tail -n +2 $(COV_OUT_FILE) >> $(COMBINED_OUT_FILE); \
fi
endef

# NOTE: We don't have prebuild-check as a dependency here because it would cause
#       the recipe to be always executed.
$(COV_PATH_UNIT): $(SOURCES)
	$(eval TEST_NAME := unit)
	$(call log-info,"Running test: $(TEST_NAME)")
	@mkdir -p $(COV_DIR)
	@echo "mode: $(COVERAGE_MODE)" > $(COV_PATH_UNIT)
	$(eval TEST_PACKAGES:=$(shell go list ./... | grep -v vendor))
	$(foreach package, $(TEST_PACKAGES), $(call test-package,$(TEST_NAME),$(package),$(COV_PATH_UNIT),-tags=unit))

# NOTE: We don't have prebuild-check as a dependency here because it would cause
#       the recipe to be always executed.
$(COV_PATH_INTEGRATION): $(SOURCES)
	$(eval TEST_NAME := integration)
	$(call log-info,"Running test: $(TEST_NAME)")
	@mkdir -p $(COV_DIR)
	@echo "mode: $(COVERAGE_MODE)" > $(COV_PATH_INTEGRATION)
	$(eval TEST_PACKAGES:=$(shell go list ./... | grep -v vendor))
	$(foreach package, $(TEST_PACKAGES), $(call test-package,$(TEST_NAME),$(package),$(COV_PATH_INTEGRATION),-tags=integration))

#-------------------------------------------------------------------------------
# Additional tools to build
#-------------------------------------------------------------------------------

$(GOCOV_BIN): prebuild-check
	@cd $(VENDOR_DIR)/github.com/axw/gocov/gocov/ && go build

$(GOCOVMERGE_BIN): prebuild-check
	@cd $(VENDOR_DIR)/github.com/wadey/gocovmerge && go build

#-------------------------------------------------------------------------------
# Clean targets
#-------------------------------------------------------------------------------

CLEAN_TARGETS += clean-coverage
.PHONY: clean-coverage
## Removes all coverage files
clean-coverage: clean-coverage-unit clean-coverage-integration clean-coverage-overall
	-@rm -rf $(COV_DIR)

CLEAN_TARGETS += clean-coverage-overall
.PHONY: clean-coverage-overall
## Removes overall coverage file
clean-coverage-overall:
	-@rm -f $(COV_PATH_OVERALL)

CLEAN_TARGETS += clean-coverage-unit
.PHONY: clean-coverage-unit
## Removes unit test coverage file
clean-coverage-unit:
	-@rm -f $(COV_PATH_UNIT)

CLEAN_TARGETS += clean-coverage-integration
.PHONY: clean-coverage-integration
## Removes integreation test coverage file
clean-coverage-integration:
	-@rm -f $(COV_PATH_INTEGRATION)
