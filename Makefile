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

WAIT4X_BINARY_NAME ?= wait4x

# build output path
WAIT4X_BUILD_OUTPUT ?= ${CURDIR}/bin

# commit ref slug used for `wait4x version`
WAIT4X_COMMIT_REF_SLUG ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --always)
# commit short SHA (8 char) used for `wait4x version`
WAIT4X_COMMIT_SHORT_SHA ?= $(shell git rev-parse --verify --short=8 HEAD)
# commit datetime used for `wait4x version`
WAIT4X_COMMIT_DATETIME ?= $(shell git log -1 --format="%at" | TZ=UTC xargs -I{} date -d @{} '+%FT%TZ')

# build flags for the Wait4X binary
# - reproducible builds: -trimpath and -ldflags=-buildid=
# - smaller binaries: -w (trim debugger data, but not panics)
# - metadata: -X=... to bake in git commit
WAIT4X_BUILD_FLAGS ?= -trimpath -ldflags="-buildid= -w -X github.com/atkrad/wait4x/internal/app/wait4x/cmd.AppVersion=$(WAIT4X_COMMIT_REF_SLUG) -X github.com/atkrad/wait4x/internal/app/wait4x/cmd.GitCommit=$(WAIT4X_COMMIT_SHORT_SHA) -X github.com/atkrad/wait4x/internal/app/wait4x/cmd.BuildTime=$(WAIT4X_COMMIT_DATETIME)"

# run flags for run target
WAIT4X_RUN_FLAGS ?= -ldflags="-X github.com/atkrad/wait4x/internal/app/wait4x/cmd.AppVersion=$(WAIT4X_COMMIT_REF_SLUG) -X github.com/atkrad/wait4x/internal/app/wait4x/cmd.GitCommit=$(WAIT4X_COMMIT_SHORT_SHA) -X github.com/atkrad/wait4x/internal/app/wait4x/cmd.BuildTime=$(WAIT4X_COMMIT_DATETIME)"

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
	@echo '  You can pass subcommand and arguements with "Wait4X" e.g. "make run WAIT4X_FLAGS.'
	@echo ""

test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

check-gofmt:
	@ if [ -n "$(shell gofmt -s -l .)" ]; then \
		echo "Go code is not formatted, run 'gofmt -s -w .'"; \
		exit 1; \
	fi

check-revive:
	revive -config .revive.toml -formatter friendly ./...

build:
	go build -v $(WAIT4X_BUILD_FLAGS) -o $(WAIT4X_BUILD_OUTPUT)/$(WAIT4X_BINARY_NAME) cmd/wait4x/main.go

run:
	go run $(WAIT4X_RUN_FLAGS) cmd/wait4x/main.go $(WAIT4X_FLAGS)
