# Copyright 2020 Mohammad Abdolirad
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO_BINARY ?= $(shell which go)
GO_ENVIRONMENTS ?=

# Wait4X output name
WAIT4X_BINARY_NAME ?= wait4x
ifeq ($(GOOS),windows)
WAIT4X_BINARY_NAME := ${WAIT4X_BINARY_NAME}.exe
endif

# build output path
WAIT4X_BUILD_OUTPUT ?= ${CURDIR}/dist

# commit ref slug used for `wait4x version`
WAIT4X_COMMIT_REF_SLUG ?= $(shell [ -d ./.git ] && (git symbolic-ref -q --short HEAD || git describe --tags --always))

# commit hash used for `wait4x version`
WAIT4X_COMMIT_HASH ?= $(shell [ -d ./.git ] && git rev-parse HEAD)
# build time used for `wait4x version`
WAIT4X_BUILD_TIME ?= $(shell date -u '+%FT%TZ')

# build flags for the Wait4X binary
# - reproducible builds: -ldflags=-buildid=
# - smaller binaries: -w (trim debugger data, but not panics)
# - metadata: -X=... to bake in git commit
WAIT4X_BUILD_LDFLAGS ?= -buildid= -w -X github.com/atkrad/wait4x/v2/internal/app/wait4x/cmd.BuildTime=$(WAIT4X_BUILD_TIME)

# pass the AppVersion if the WAIT4X_COMMIT_REF_SLUG isn't empty.
ifneq ($(WAIT4X_COMMIT_REF_SLUG),)
WAIT4X_BUILD_LDFLAGS += -X github.com/atkrad/wait4x/v2/internal/app/wait4x/cmd.AppVersion=$(WAIT4X_COMMIT_REF_SLUG)
endif

# pass the GitCommit if the WAIT4X_COMMIT_HASH isn't empty.
ifneq ($(WAIT4X_COMMIT_HASH),)
WAIT4X_BUILD_LDFLAGS += -X github.com/atkrad/wait4x/v2/internal/app/wait4x/cmd.GitCommit=$(WAIT4X_COMMIT_HASH)
endif

# build flags for the Wait4X binary
# - reproducible builds: -trimpath
WAIT4X_BUILD_FLAGS ?= -trimpath -ldflags="$(WAIT4X_BUILD_LDFLAGS)"

# run flags for run target
WAIT4X_RUN_FLAGS ?= -ldflags="$(WAIT4X_BUILD_LDFLAGS)"

# flags for wait4x
WAIT4X_FLAGS ?=

help:
	@echo " __      __        .__  __     _________  ___"
	@echo "/  \    /  \_____  |__|/  |_  /  |  \   \/  /"
	@echo "\   \/\/   /\__  \ |  \   __\/   |  |\     / "
	@echo " \        /  / __ \|  ||  | /    ^   /     \ "
	@echo "  \__/\  /  (____  /__||__| \____   /___/\  \\"
	@echo "       \/        \/              |__|     \_/"
	@echo ""
	@echo ""
	@echo "build"
	@echo "  Build Wait4X."
	@echo ""
	@echo "run"
	@echo '  Run Wait4X.'
	@echo '  You can pass subcommand and arguments with "Wait4X" e.g. "make run WAIT4X_FLAGS.'
	@echo ""

test:
	$(GO_ENVIRONMENTS) $(GO_BINARY) test -v -covermode=count -coverprofile=coverage.out ./...

check-gofmt:
	@ if [ -n "$(shell gofmt -s -l .)" ]; then \
		echo "Go code is not formatted, run 'gofmt -s -w .'"; \
		exit 1; \
	fi

check-revive:
	revive -config .revive.toml -formatter friendly ./...

build:
	$(GO_ENVIRONMENTS) $(GO_BINARY) build -v $(WAIT4X_BUILD_FLAGS) -o $(WAIT4X_BUILD_OUTPUT)/$(WAIT4X_BINARY_NAME) cmd/wait4x/main.go

run:
	$(GO_ENVIRONMENTS) $(GO_BINARY) run $(WAIT4X_RUN_FLAGS) cmd/wait4x/main.go $(WAIT4X_FLAGS)
